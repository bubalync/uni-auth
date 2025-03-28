package user

import (
	"context"
	"github.com/bubalync/uni-auth/internal/entity"
	"github.com/bubalync/uni-auth/internal/repo"
	"github.com/bubalync/uni-auth/internal/services"
	"github.com/bubalync/uni-auth/pkg/logger/sl"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"strings"
)

// Service -.
type Service struct {
	repo repo.UserRepo
	log  *slog.Logger
}

// New -.
func New(log *slog.Logger, r repo.UserRepo) *Service {
	return &Service{
		r,
		log,
	}
}

func (s *Service) Delete(ctx context.Context, u *entity.User) error {
	//TODO implement me
	panic("implement me")
}

func (s *Service) Login(ctx context.Context, u *entity.User) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) Logout(ctx context.Context, u *entity.User) error {
	//TODO implement me
	panic("implement me")
}

func (s *Service) Register(ctx context.Context, email, password string) (uuid.UUID, error) {
	const op = "services.user.Register"
	log := s.log.With(slog.String("op", op))

	alreadyExists, err := s.repo.UserByEmailIsExists(ctx, email)
	if err != nil {
		log.Error("failed check user existence by email", sl.Err(err))
		return uuid.Nil, services.ErrInternal
	}

	if *alreadyExists {
		return uuid.Nil, services.ErrUserExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate hashed password", sl.Err(err))
		return uuid.Nil, services.ErrInternal
	}

	user := &entity.User{
		ID:           uuid.New(),
		Email:        strings.ToLower(email),
		PasswordHash: hashedPassword,
	}
	if err = s.repo.Create(ctx, user); err != nil {
		log.Error("failed to create new user", sl.Err(err))
		return uuid.Nil, services.ErrInternal
	}

	return user.ID, nil
}

func (s *Service) ResetPassword(ctx context.Context, u *entity.User) error {
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
