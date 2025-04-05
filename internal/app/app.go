package app

import (
	"fmt"
	"github.com/bubalync/uni-auth/internal/api/http"
	"github.com/bubalync/uni-auth/internal/config"
	repo "github.com/bubalync/uni-auth/internal/repo/persistent"
	"github.com/bubalync/uni-auth/internal/service"
	"github.com/bubalync/uni-auth/pkg/httpserver"
	"github.com/bubalync/uni-auth/pkg/logger"
	"github.com/bubalync/uni-auth/pkg/logger/sl"
	"github.com/bubalync/uni-auth/pkg/postgres"
	"github.com/gin-gonic/gin"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	log := logger.New(cfg.Env, cfg.Log.Level)

	// Postgres
	pg, err := postgres.New(cfg.PG.Url, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		log.Error("app - Run - postgres.New", sl.Err(err))
	}
	defer pg.Close()

	// services
	services := service.NewServices(log, cfg, repo.NewUserRepo(pg))

	// Gin handler
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
}
