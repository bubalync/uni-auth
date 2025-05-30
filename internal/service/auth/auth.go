package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/bubalync/uni-auth/internal/entity"
	"github.com/bubalync/uni-auth/internal/lib/email"
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
	refreshKeyTemplate = "refresh:%s"
	resetKeyTemplate   = "reset:%s"
)

type Service struct {
	log             *slog.Logger
	cache           redis.Cache
	userRepo        repo.User
	hasher          hasher.PasswordHasher
	tokenGenerator  jwtgen.TokenGenerator
	refreshTokenTTL time.Duration
	emailSender     email.Sender
}

// New -.
func New(
	log *slog.Logger,
	cache redis.Cache,
	userRepo repo.User,
	hasher hasher.PasswordHasher,
	tokenGenerator jwtgen.TokenGenerator,
	emailSender email.Sender,
	refreshTokenTTL time.Duration,
) *Service {
	return &Service{
		log:             log,
		cache:           cache,
		userRepo:        userRepo,
		hasher:          hasher,
		tokenGenerator:  tokenGenerator,
		emailSender:     emailSender,
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

	return s.generateTokens(ctx, log, user)
}

func (s *Service) Refresh(ctx context.Context, token string) (GenerateTokenOutput, error) {
	const op = "service.auth.Refresh"
	log := s.log.With(slog.String("op", op))

	claims, err := s.tokenGenerator.ParseRefreshToken(token)
	if err != nil {
		log.Error("failed to parse refresh token", sl.Err(err))
		return GenerateTokenOutput{}, svcErrs.ErrCannotParseToken
	}

	stored, err := s.cache.Get(ctx, fmt.Sprintf(refreshKeyTemplate, claims.UserId))
	if err != nil {
		log.Error("failed to refresh token", sl.Err(err))
		return GenerateTokenOutput{}, svcErrs.ErrTokenIsExpired
	}

	if token != stored {
		log.Error("stored and input token is not equal")
		return GenerateTokenOutput{}, svcErrs.ErrTokenIsExpired
	}

	return s.generateTokens(ctx, log, entity.User{Id: claims.UserId, Email: claims.Email})
}

func (s *Service) generateTokens(ctx context.Context, log *slog.Logger, user entity.User) (GenerateTokenOutput, error) {
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

	//TODO delete field and func`s
	if err = s.userRepo.UpdateLastLoginAttempt(ctx, user.Id); err != nil {
		log.Error("failed to update last_login_attempt", sl.Err(err))
	}

	// TODO solve the problem with authentication from two different realms.
	//  storing token as a key?
	err = s.cache.Set(ctx, fmt.Sprintf(refreshKeyTemplate, user.Id), refreshToken, s.refreshTokenTTL)
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
	const op = "service.auth.ResetPassword"
	log := s.log.With(slog.String("op", op))

	isExists, err := s.userRepo.UserByEmailIsExists(ctx, input.Email)
	if err != nil {
		log.Error("failed to get user from db", sl.Err(err))
		return svcErrs.ErrCannotGetUser
	}

	if !*isExists {
		return svcErrs.ErrUserNotFound
	}

	// todo think about generating a token
	token := uuid.New().String() + uuid.New().String() + uuid.New().String()

	if err := s.cache.Set(ctx, fmt.Sprintf(resetKeyTemplate, token), input.Email, 15*time.Minute); err != nil {
		log.Error("failed to save the reset password token to cache", sl.Err(err))
		return svcErrs.ErrAccessToCache
	}

	if err := s.emailSender.SendResetPasswordEmail(input.Email, token); err != nil {
		log.Error("failed to send the reset password email", sl.Err(err))
		return svcErrs.ErrSendResetPasswordEmail
	}

	return nil
}

func (s *Service) RecoveryPassword(ctx context.Context, input RecoveryPasswordInput) error {
	const op = "service.auth.RecoveryPassword"
	log := s.log.With(slog.String("op", op))

	userEmail, err := s.cache.Get(ctx, fmt.Sprintf(resetKeyTemplate, input.Token))
	if err != nil {
		log.Error("failed to get the reset token", sl.Err(err))
		return svcErrs.ErrTokenIsExpired
	}

	pwd, err := s.hasher.Hash(input.Password)
	if err != nil {
		log.Error("failed to generate hashed password", sl.Err(err))
		return svcErrs.ErrCannotUpdateUser
	}

	err = s.userRepo.UpdatePassword(ctx, userEmail, pwd)
	if err != nil {
		log.Error("failed to update password", sl.Err(err))
		return svcErrs.ErrCannotUpdateUser
	}

	err = s.cache.Delete(ctx, fmt.Sprintf(resetKeyTemplate, input.Token))
	if err != nil {
		log.Error("failed to delete token from cache", sl.Err(err))
	}

	return nil
}

func (s *Service) ParseToken(token string) (*jwtgen.Claims, error) {
	const op = "service.auth.ParseToken"
	log := s.log.With(slog.String("op", op))

	claims, err := s.tokenGenerator.ParseAccessToken(token)
	if err != nil {
		log.Error("failed to parse access token", sl.Err(err))
		return nil, svcErrs.ErrCannotParseToken
	}

	return claims, nil
}
