package user

import (
	"context"
	"github.com/bubalync/uni-auth/internal/entity"
	"github.com/bubalync/uni-auth/internal/repo"
	"github.com/google/uuid"
	"log/slog"
)

// Service -.
type Service struct {
	repo repo.User
	log  *slog.Logger
}

// New -.
func New(log *slog.Logger, r repo.User) *Service {
	return &Service{
		r,
		log,
	}
}

func (s *Service) Delete(ctx context.Context, u *entity.User) error {
	//TODO implement me
	panic("implement me")
}

func (s *Service) Logout(ctx context.Context, u *entity.User) error {
	//TODO implement me
	panic("implement me")
}

func (s *Service) Update(ctx context.Context, u *entity.User) error {
	//TODO implement me
	panic("implement me")
}

func (s *Service) UserByEmail(ctx context.Context, email string) (*entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) UserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	//TODO implement me
	panic("implement me")
}
