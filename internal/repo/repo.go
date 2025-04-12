package repo

import (
	"context"
	"github.com/bubalync/uni-auth/internal/entity"
	"github.com/bubalync/uni-auth/internal/repo/persistent"
	"github.com/bubalync/uni-auth/pkg/postgres"
	"github.com/google/uuid"
)

type (
	User interface {
		Create(ctx context.Context, u entity.User) error
		Delete(ctx context.Context, id uuid.UUID) error
		ResetPassword(ctx context.Context, passwordHash []byte) error
		Update(ctx context.Context, u entity.User) error
		UpdateLastLoginAttempt(ctx context.Context, id uuid.UUID) error
		UserByEmail(ctx context.Context, email string) (entity.User, error)
		UserByEmailIsExists(ctx context.Context, email string) (*bool, error)
		UserByID(ctx context.Context, id uuid.UUID) (entity.User, error)
	}
)

type Repositories struct {
	User
}

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		User: persistent.NewUserRepo(pg),
	}
}
