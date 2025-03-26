package repo

import (
	"context"
	"github.com/bubalync/uni-auth/internal/entity"
	"github.com/google/uuid"
)

type UserRepo interface {
	Create(ctx context.Context, u *entity.User) error
	Update(ctx context.Context, u *entity.User) (string, error)
	ResetPassword(ctx context.Context, passwordHash []byte) error
	UserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	UserByEmail(ctx context.Context, email string) (*entity.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
