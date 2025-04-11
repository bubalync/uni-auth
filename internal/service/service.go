package service

import (
	"context"
	"github.com/bubalync/uni-auth/internal/entity"
	"github.com/bubalync/uni-auth/internal/repo"
	"github.com/bubalync/uni-auth/internal/service/auth"
	"github.com/bubalync/uni-auth/internal/service/user"
	"github.com/bubalync/uni-auth/pkg/hasher"
	"github.com/bubalync/uni-auth/pkg/redis"
	"github.com/google/uuid"
	"log/slog"
	"time"
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

type (
	ServicesDependencies struct {
		Repos  *repo.Repositories
		Hasher hasher.PasswordHasher
		Cache  redis.Cache

		SignKey  string
		TokenTTL time.Duration
	}

	Services struct {
		Auth Auth
		User User
	}
)

func NewServices(log *slog.Logger, deps ServicesDependencies) *Services {
	return &Services{
		Auth: auth.New(log, deps.Repos.User, deps.Hasher, deps.SignKey, deps.TokenTTL),
		User: user.New(log, deps.Repos.User),
	}
}
