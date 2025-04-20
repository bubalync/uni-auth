package app

import (
	"fmt"
	"github.com/bubalync/uni-auth/internal/api/grpc"
	"github.com/bubalync/uni-auth/internal/api/http"
	"github.com/bubalync/uni-auth/internal/config"
	"github.com/bubalync/uni-auth/internal/lib/email"
	"github.com/bubalync/uni-auth/internal/lib/jwtgen"
	"github.com/bubalync/uni-auth/internal/repo"
	"github.com/bubalync/uni-auth/internal/service"
	"github.com/bubalync/uni-auth/pkg/hasher"
	"github.com/bubalync/uni-auth/pkg/httpserver"
	"github.com/bubalync/uni-auth/pkg/logger"
	"github.com/bubalync/uni-auth/pkg/logger/sl"
	"github.com/bubalync/uni-auth/pkg/postgres"
	"github.com/bubalync/uni-auth/pkg/redis"
	"github.com/gin-gonic/gin"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	log := logger.New(cfg.Env, cfg.Log.Level)

	// Postgres
	log.Info("Initializing postgres...")
	pg, err := postgres.New(cfg.PG.Url, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		log.Error("app - Run - postgres.New", sl.Err(err))
	}
	defer pg.Close()

	// Redis
	redisClient := redis.NewRedisClient(redis.Addr(cfg.Redis.Host), redis.Db(cfg.Redis.Db))

	// Repositories
	log.Info("Initializing repositories...")
	repositories := repo.NewRepositories(pg)

	// services
	log.Info("Initializing services...")
	deps := service.ServicesDependencies{
		Repos:  repositories,
		Hasher: hasher.NewBcryptHasher(),
		Cache:  redisClient,
		TokenGenerator: jwtgen.NewJwtTokenGenerator(
			cfg.JWT.AccessSignKey,
			cfg.JWT.RefreshSignKey,
			cfg.JWT.AccessTokenTTL,
			cfg.JWT.RefreshTokenTTL,
		),
		RefreshTokenTTL: cfg.JWT.RefreshTokenTTL,
		EmailSender: email.NewSmtpSender(
			cfg.EmailSender.SMTPHost,
			cfg.EmailSender.SMTPPort,
			cfg.EmailSender.Username,
			cfg.EmailSender.Password,
			cfg.EmailSender.From,
		),
	}
	services := service.NewServices(log, deps)

	// gRPC server
	gRPCServer := grpc.NewServer(log, services, cfg.GRPC.Port)

	go func() {
		gRPCServer.MustRun()
	}()

	// Gin handler
	log.Info("Initializing handlers and routes...")
	handler := gin.New()
	http.NewRouter(handler, cfg, log, services)

	// HTTP server
	httpServer := httpserver.New(
		handler,
		httpserver.Port(cfg.HTTP.Port),
		httpserver.ReadTimeout(cfg.HTTP.Timeout),
		httpserver.WriteTimeout(cfg.HTTP.Timeout),
		httpserver.IdleTimeout(cfg.HTTP.IdleTimeout),
	)

	log.Info("Starting http server...", slog.String("Port", cfg.HTTP.Port))
	httpServer.Start()

	log.Info(fmt.Sprintf("%s service ready to work", cfg.App.Name))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-interrupt:
		log.Info("app - Run - signal: " + sig.String())
	case err = <-httpServer.Notify():
		log.Error("app - Run - httpServer.Notify:", sl.Err(err))
	}

	// Graceful shutdown
	log.Info("Shutting down server...")
	err = httpServer.Shutdown()
	if err != nil {
		log.Error("app - Run - httpServer.Shutdown:", sl.Err(err))
	}

	err = redisClient.Close()
	if err != nil {
		log.Error("app - Run - redis.Close:", sl.Err(err))
	}

	gRPCServer.Stop()
}
