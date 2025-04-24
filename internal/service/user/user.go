package user

import (
	"context"
	"errors"
	"github.com/bubalync/uni-auth/internal/entity"
	"github.com/bubalync/uni-auth/internal/repo"
	"github.com/bubalync/uni-auth/internal/repo/repoErrs"
	"github.com/bubalync/uni-auth/internal/service/svcErrs"
	"github.com/bubalync/uni-auth/pkg/logger/sl"
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

func (s *Service) Delete(ctx context.Context, u entity.User) error {
	//TODO implement me
	panic("implement me")
}

func (s *Service) Logout(ctx context.Context, u entity.User) error {
	//TODO implement me
	panic("implement me")
}

func (s *Service) Update(ctx context.Context, u entity.User) error {
	//TODO implement me
	panic("implement me")
}

func (s *Service) UserByEmail(ctx context.Context, email string) (entity.User, error) {
	log := s.log.With(slog.String("op", "service.user.UserByEmail"))

	user, err := s.repo.UserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repoErrs.ErrNotFound) {
			log.Error("Cannot get user", sl.Err(err))
			return entity.User{}, svcErrs.ErrUserNotFound
		}

		log.Error("Cannot get user", sl.Err(err))
		return entity.User{}, svcErrs.ErrCannotGetUser
	}

	return user, nil
}

func (s *Service) UserById(ctx context.Context, id uuid.UUID) (entity.User, error) {
	log := s.log.With(slog.String("op", "service.user.UserById"))

	user, err := s.repo.UserById(ctx, id)
	if err != nil {
		if errors.Is(err, repoErrs.ErrNotFound) {
			log.Error("Cannot get user", sl.Err(err))
			return entity.User{}, svcErrs.ErrUserNotFound
		}

		log.Error("Cannot get user", sl.Err(err))
		return entity.User{}, svcErrs.ErrCannotGetUser
	}

	return user, nil
}
