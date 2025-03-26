package user

import (
	"context"
	"github.com/bubalync/uni-auth/internal/entity"
	"github.com/bubalync/uni-auth/internal/repo"
	"github.com/google/uuid"
)

// Service -.
type Service struct {
	repo repo.UserRepo
}

// New -.
func New(r repo.UserRepo) *Service {
	return &Service{repo: r}
}

func (us *Service) Register(ctx context.Context, u *entity.User) error {
	//TODO implement me
	panic("implement me")
}

func (us *Service) Login(ctx context.Context, u *entity.User) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (us *Service) ResetPassword(ctx context.Context, u *entity.User) error {
	//TODO implement me
	panic("implement me")
}

func (us *Service) Update(ctx context.Context, u *entity.User) error {
	//TODO implement me
	panic("implement me")
}

func (us *Service) UserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (us *Service) UserByEmail(ctx context.Context, email string) (*entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (us *Service) Logout(ctx context.Context, u *entity.User) error {
	//TODO implement me
	panic("implement me")
}

func (us *Service) Delete(ctx context.Context, u *entity.User) error {
	//TODO implement me
	panic("implement me")
}
