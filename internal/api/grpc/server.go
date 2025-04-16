package grpc

import (
	"context"
	"fmt"
	v1 "github.com/bubalync/uni-auth/internal/api/grpc/v1"
	"github.com/bubalync/uni-auth/internal/service"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

type Server struct {
	log     *slog.Logger
	grpcSrv *grpc.Server
	port    int
}

func NewServer(log *slog.Logger, services *service.Services, port int) *Server {
	// TODO continue server setup: otel, etc...

	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.StartCall, logging.FinishCall,
			logging.PayloadReceived, logging.PayloadSent,
		),
		// Add any other option (check functions starting with logging.With).
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			log.Error("Recovered from panic", slog.Any("panic", p))

			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	grpcSrv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(recoveryOpts...),
			logging.UnaryServerInterceptor(interceptorLogger(log), loggingOpts...),
		),
	)

	// handlers
	v1.NewAuthServer(grpcSrv, services.Auth)

	return &Server{
		log:     log,
		grpcSrv: grpcSrv,
		port:    port,
	}
}

// interceptorLogger adapts slog logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func interceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

// MustRun runs gRPC server and panics if any error occurs.
func (s *Server) MustRun() {
	if err := s.Run(); err != nil {
		panic(err)
	}

}

// Run runs gRPC server.
func (s *Server) Run() error {
	const op = "api.grpc.server.Run"

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	s.log.Info("grpc server started", slog.String("addr", l.Addr().String()))

	if err := s.grpcSrv.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Server) Stop() {
	const op = "api.grpc.server.Stop"

	s.log.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.Int("port", s.port))

	s.grpcSrv.GracefulStop()
}
