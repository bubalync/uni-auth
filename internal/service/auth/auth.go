package auth

import (
	"context"
	"errors"
	"github.com/bubalync/uni-auth/internal/entity"
	"github.com/bubalync/uni-auth/internal/lib/jwtgen"
	"github.com/bubalync/uni-auth/internal/repo"
	"github.com/bubalync/uni-auth/internal/repo/repoErrs"
	"github.com/bubalync/uni-auth/internal/service/svcErrs"
	"github.com/bubalync/uni-auth/pkg/hasher"
	"github.com/bubalync/uni-auth/pkg/logger/sl"
	"github.com/bubalync/uni-auth/pkg/redis"
	"github.com/google/uuid"
	"log/slog"
	"strings"
	"time"
)

const (
	cacheRefreshKeyPrefix = "refresh:"
)

type Service struct {
	log             *slog.Logger
	cache           redis.Cache
	userRepo        repo.User
	hasher          hasher.PasswordHasher
	tokenGenerator  jwtgen.TokenGenerator
	refreshTokenTTL time.Duration
}

// New -.
func New(log *slog.Logger, cache redis.Cache, userRepo repo.User, hasher hasher.PasswordHasher, tokenGenerator jwtgen.TokenGenerator, refreshTokenTTL time.Duration) *Service {
	return &Service{
		log:             log,
		cache:           cache,
		userRepo:        userRepo,
		hasher:          hasher,
		tokenGenerator:  tokenGenerator,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (s *Service) CreateUser(ctx context.Context, input CreateUserInput) (uuid.UUID, error) {
	const op = "service.auth.CreateUser"
	log := s.log.With(slog.String("op", op))

	hashedPassword, err := s.hasher.Hash(input.Password)
	if err != nil {
		log.Error("failed to generate hashed password", sl.Err(err))
		return uuid.Nil, svcErrs.ErrCannotCreateUser
	}

	user := entity.User{
		Id:           uuid.New(),
		Email:        strings.ToLower(input.Email),
		PasswordHash: hashedPassword,
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		if errors.Is(err, repoErrs.ErrAlreadyExists) {
			log.Error("failed to create new user", sl.Err(err))
			return uuid.Nil, svcErrs.ErrUserAlreadyExists
		}

		log.Error("failed to create new user", sl.Err(err))
		return uuid.Nil, svcErrs.ErrCannotCreateUser
	}
	return user.Id, nil
}

func (s *Service) GenerateToken(ctx context.Context, input GenerateTokenInput) (GenerateTokenOutput, error) {
	const op = "service.auth.GenerateToken"
	log := s.log.With(slog.String("op", op))

	user, err := s.userRepo.UserByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, repoErrs.ErrNotFound) {
			log.Error("Cannot get user", sl.Err(err))
			return GenerateTokenOutput{}, svcErrs.ErrInvalidCredentials
		}

		log.Error("Cannot get user", sl.Err(err))
		return GenerateTokenOutput{}, svcErrs.ErrCannotGetUser
	}

	if err = s.hasher.Compare(user.PasswordHash, []byte(input.Password)); err != nil {
		log.Error("failed to compare password", sl.Err(err))
		return GenerateTokenOutput{}, svcErrs.ErrInvalidCredentials
	}

	accessToken, err := s.tokenGenerator.GenerateAccessToken(user)
	if err != nil {
		log.Error("failed to generate access token", sl.Err(err))
		return GenerateTokenOutput{}, svcErrs.ErrCannotSignToken
	}

	refreshToken, err := s.tokenGenerator.GenerateRefreshToken(user)
	if err != nil {
		log.Error("failed to generate refresh token", sl.Err(err))
		return GenerateTokenOutput{}, svcErrs.ErrCannotSignToken
	}

	if err = s.userRepo.UpdateLastLoginAttempt(ctx, user.Id); err != nil {
		log.Error("failed to update last_login_attempt", sl.Err(err))
	}

	err = s.cache.Set(ctx, cacheRefreshKeyPrefix+user.Id.String(), refreshToken, s.refreshTokenTTL)
	if err != nil {
		log.Error("failed to save refresh token to cache", sl.Err(err))
		return GenerateTokenOutput{}, svcErrs.ErrAccessToCache
	}

	return GenerateTokenOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) ResetPassword(ctx context.Context, input ResetPasswordInput) error {
	//TODO implement me
	panic("implement me")
}

func (s *Service) ParseToken(token string) (uuid.UUID, error) {
	const op = "service.auth.ParseToken"
	log := s.log.With(slog.String("op", op))

	claims, err := s.tokenGenerator.ParseAccessToken(token)
	if err != nil {
		log.Error("failed to parse access token", sl.Err(err))
		return uuid.Nil, svcErrs.ErrCannotParseToken
	}

	return claims.UserId, nil
}
