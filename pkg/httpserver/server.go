package httpserver

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

type Server struct {
	Router *gin.Engine
	server *http.Server

	address      string
	readTimeout  time.Duration
	writeTimeout time.Duration
	idleTimeout  time.Duration
}

func NewServer(log *slog.Logger, opts ...Option) *Server {
	s := &Server{}

	router := gin.New()

	// Middleware
	router.Use(gin.Recovery())
	router.Use(sloggin.New(log))

	// Custom options
	for _, opt := range opts {
		opt(s)
	}

	s.Router = router
	s.server = &http.Server{
		Addr:         s.address,
		Handler:      s.Router,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
		IdleTimeout:  s.idleTimeout,
	}

	return s
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	if s.server != nil {
		return s.server.Shutdown(ctx)
	}
	return nil
}
