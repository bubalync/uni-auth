package auth

import (
	"context"
	"errors"
	"github.com/bubalync/uni-auth/internal/entity"
	"github.com/bubalync/uni-auth/internal/repo"
	"github.com/bubalync/uni-auth/internal/repo/repoErrs"
	"github.com/bubalync/uni-auth/internal/service/svcErrs"
	"github.com/bubalync/uni-auth/pkg/hasher"
	"github.com/bubalync/uni-auth/pkg/logger/sl"
	"github.com/google/uuid"
	"log/slog"
	"strings"
	"time"
)

type Service struct {
	userRepo repo.User
	hasher   hasher.PasswordHasher
	log      *slog.Logger

	signKey  string
	tokenTTL time.Duration
}

// New -.
func New(log *slog.Logger, userRepo repo.User, hasher hasher.PasswordHasher, signKey string, tokenTTL time.Duration) *Service {
	return &Service{
		userRepo: userRepo,
		hasher:   hasher,
		log:      log,
		signKey:  signKey,
		tokenTTL: tokenTTL,
	}
}

func (s *Service) CreateUser(ctx context.Context, input CreateUserInput) (uuid.UUID, error) {
	const op = "service.auth.CreateUser"
	log := s.log.With(slog.String("op", op))

	hashedPassword, err := s.hasher.Hash(input.Password)
	if err != nil {
		log.Error("failed to generate hashed password", sl.Err(err))
		return uuid.Nil, svcErrs.ErrInternal
	}

	user := entity.User{
		ID:           uuid.New(),
		Email:        strings.ToLower(input.Email),
		PasswordHash: hashedPassword,
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		if errors.Is(err, repoErrs.ErrAlreadyExists) {
			return uuid.Nil, svcErrs.ErrUserAlreadyExists
		}

		log.Error("failed to create new user", sl.Err(err))
		return uuid.Nil, svcErrs.ErrInternal
	}
	return user.ID, nil
}

func (s *Service) GenerateToken(ctx context.Context, input GenerateTokenInput) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) ResetPassword(ctx context.Context, input ResetPasswordInput) error {
	//TODO implement me
	panic("implement me")
}

func (s *Service) ParseToken(token string) (uuid.UUID, error) {
	//TODO implement me
	panic("implement me")
}
