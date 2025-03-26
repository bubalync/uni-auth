package services

import (
	"context"
	"github.com/bubalync/uni-auth/internal/entity"
	"github.com/google/uuid"
)

type User interface {
	Register(ctx context.Context, u *entity.User) error
	Login(ctx context.Context, u *entity.User) (string, error)
	Logout(ctx context.Context, u *entity.User) error
	ResetPassword(ctx context.Context, u *entity.User) error
	Update(ctx context.Context, u *entity.User) error
	UserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	UserByEmail(ctx context.Context, email string) (*entity.User, error)
	Delete(ctx context.Context, u *entity.User) error
}
