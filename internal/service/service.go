package service

import (
	"context"
	"github.com/bubalync/uni-auth/internal/config"
	"github.com/bubalync/uni-auth/internal/entity"
	"github.com/bubalync/uni-auth/internal/repo"
	"github.com/bubalync/uni-auth/internal/service/auth"
	"github.com/bubalync/uni-auth/internal/service/user"
	"github.com/google/uuid"
	"log/slog"
)

type (
	Auth interface {
		CreateUser(ctx context.Context, input auth.CreateUserInput) (uuid.UUID, error)
		GenerateToken(ctx context.Context, input auth.GenerateTokenInput) (string, error)
		ResetPassword(ctx context.Context, input auth.ResetPasswordInput) error
		ParseToken(token string) (uuid.UUID, error)
	}

	User interface {
		Delete(ctx context.Context, u entity.User) error
		Logout(ctx context.Context, u entity.User) error
		Update(ctx context.Context, u entity.User) error
		UserByEmail(ctx context.Context, email string) (entity.User, error)
		UserByID(ctx context.Context, id uuid.UUID) (entity.User, error)
	}
)

type Services struct {
	Auth Auth
	User User
}

func NewServices(log *slog.Logger, cfg *config.Config, userRepo repo.User) *Services {
	return &Services{
		Auth: auth.New(log, userRepo, cfg.JWT.SignKey, cfg.JWT.TokenTTL),
		User: user.New(log, userRepo),
	}
}
