package app

import (
	"context"
	"errors"
	"github.com/bubalync/uni-auth/internal/api/http"
	"github.com/bubalync/uni-auth/internal/config"
	"github.com/bubalync/uni-auth/pkg/httpserver"
	"github.com/bubalync/uni-auth/pkg/logger"
	"github.com/bubalync/uni-auth/pkg/logger/sl"
	"github.com/bubalync/uni-auth/pkg/postgres"
	"log/slog"
	nethttp "net/http"

	"os"
	"os/signal"
	"syscall"
	"time"
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

	// HTTP Server
	server := httpserver.NewServer(
		log,
		httpserver.Port(cfg.HTTP.Port),
		httpserver.ReadTimeout(cfg.HTTP.Timeout),
		httpserver.WriteTimeout(cfg.HTTP.Timeout),
		httpserver.IdleTimeout(cfg.HTTP.IdleTimeout),
	)

	http.FillRouter(server.Router, cfg, log)

	log.Info("Starting server.", slog.String("Port", cfg.HTTP.Port))
	go func() {
		if err := server.Start(); err != nil && !errors.Is(err, nethttp.ErrServerClosed) {
			log.Error("listen", sl.Err(err))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Stop(ctx); err != nil {
		log.Error("Server Shutdown:", sl.Err(err))
	}
	<-ctx.Done()
	log.Info("Timeout of 5 seconds.")
	log.Info("Server exiting")
}
