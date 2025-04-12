package auth

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/bubalync/uni-auth/internal/entity"
	"github.com/bubalync/uni-auth/internal/lib/jwtgen"
	"github.com/bubalync/uni-auth/internal/mocks/redismocks"
	"github.com/bubalync/uni-auth/internal/mocks/repomocks"
	"github.com/bubalync/uni-auth/internal/mocks/utilmocks"
	"github.com/bubalync/uni-auth/internal/repo/repoErrs"
	"github.com/bubalync/uni-auth/internal/service/svcErrs"
	"github.com/bubalync/uni-auth/pkg/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

const (
	refreshTokenTTL = time.Minute
)

type partialUserMatcher struct {
	Email        string
	PasswordHash []byte
}

func (m partialUserMatcher) Matches(x interface{}) bool {
	user, ok := x.(entity.User)
	return ok && user.Email == m.Email && bytes.Equal(user.PasswordHash, m.PasswordHash)
}

func (m partialUserMatcher) String() string {
	return fmt.Sprintf("matches User with Email=%s", m.Email)
}

func TestAuthService_CreateUser(t *testing.T) {
	type args struct {
		ctx   context.Context
		input CreateUserInput
	}

	type MockBehavior func(o *repomocks.MockUser, h *utilmocks.MockPasswordHasher, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		wantErr      bool
		err          error
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				input: CreateUserInput{
					Email:    "test@example.com",
					Password: "Qwerty!1",
				},
			},
			mockBehavior: func(r *repomocks.MockUser, h *utilmocks.MockPasswordHasher, args args) {
				hash := []byte{1, 2, 3}
				h.EXPECT().Hash(args.input.Password).Return(hash, nil)

				user := partialUserMatcher{
					Email:        args.input.Email,
					PasswordHash: hash,
				}
				r.EXPECT().Create(args.ctx, user).
					Return(nil)
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "hasher error",
			args: args{
				ctx: context.Background(),
				input: CreateUserInput{
					Email:    "test@example.com",
					Password: "Qwerty!1",
				},
			},
			mockBehavior: func(r *repomocks.MockUser, h *utilmocks.MockPasswordHasher, args args) {
				h.EXPECT().Hash(args.input.Password).Return(nil, errors.New("some error"))
			},
			wantErr: true,
			err:     svcErrs.ErrCannotCreateUser,
		},
		{
			name: "error already exists",
			args: args{
				ctx: context.Background(),
				input: CreateUserInput{
					Email:    "test@example.com",
					Password: "Qwerty!1",
				},
			},
			mockBehavior: func(r *repomocks.MockUser, h *utilmocks.MockPasswordHasher, args args) {
				hash := []byte{1, 2, 3}
				h.EXPECT().Hash(args.input.Password).Return(hash, nil)

				user := partialUserMatcher{
					Email:        args.input.Email,
					PasswordHash: hash,
				}
				r.EXPECT().Create(args.ctx, user).
					Return(repoErrs.ErrAlreadyExists)
			},
			wantErr: true,
			err:     svcErrs.ErrUserAlreadyExists,
		},
		{
			name: "some error error",
			args: args{
				ctx: context.Background(),
				input: CreateUserInput{
					Email:    "test@example.com",
					Password: "Qwerty!1",
				},
			},
			mockBehavior: func(r *repomocks.MockUser, h *utilmocks.MockPasswordHasher, args args) {
				hash := []byte{1, 2, 3}
				h.EXPECT().Hash(args.input.Password).Return(hash, nil)

				user := partialUserMatcher{
					Email:        args.input.Email,
					PasswordHash: hash,
				}
				r.EXPECT().Create(args.ctx, user).
					Return(errors.New("some error"))
			},
			wantErr: true,
			err:     svcErrs.ErrCannotCreateUser,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// init deps
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// init repo mock
			repo := repomocks.NewMockUser(ctrl)
			hasher := utilmocks.NewMockPasswordHasher(ctrl)
			tc.mockBehavior(repo, hasher, tc.args)

			// Log
			log := logger.New("local", "info")

			// init service
			s := New(log, nil, repo, hasher, nil, refreshTokenTTL)

			// run test
			got, err := s.CreateUser(tc.args.ctx, tc.args.input)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.IsType(t, uuid.UUID{}, got)
		})
	}
}

func TestAuthService_GenerateToken(t *testing.T) {
	type args struct {
		ctx   context.Context
		input GenerateTokenInput
	}

	type MockBehavior func(o *repomocks.MockUser, h *utilmocks.MockPasswordHasher, c *redismocks.MockCache, g *utilmocks.MockTokenGenerator, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		wantErr      bool
		err          error
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				input: GenerateTokenInput{
					Email:    "test@example.com",
					Password: "Qwerty!1",
				},
			},
			mockBehavior: func(r *repomocks.MockUser, h *utilmocks.MockPasswordHasher, c *redismocks.MockCache, g *utilmocks.MockTokenGenerator, args args) {
				hash := []byte(args.input.Password)
				user := entity.User{Id: uuid.New(), PasswordHash: hash, Email: args.input.Email}

				r.EXPECT().UserByEmail(args.ctx, args.input.Email).Return(user, nil)
				h.EXPECT().Compare(hash, hash).Return(nil)
				g.EXPECT().GenerateAccessToken(user).Return("access_token", nil)
				g.EXPECT().GenerateRefreshToken(user).Return("refresh_token", nil)
				r.EXPECT().UpdateLastLoginAttempt(args.ctx, user.Id).Return(nil)
				c.EXPECT().Set(args.ctx, "refresh:"+user.Id.String(), "refresh_token", refreshTokenTTL).Return(nil)
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "get user: user not found",
			args: args{
				ctx: context.Background(),
				input: GenerateTokenInput{
					Email:    "test@example.com",
					Password: "Qwerty!1",
				},
			},
			mockBehavior: func(r *repomocks.MockUser, h *utilmocks.MockPasswordHasher, c *redismocks.MockCache, g *utilmocks.MockTokenGenerator, args args) {
				r.EXPECT().UserByEmail(args.ctx, args.input.Email).Return(entity.User{}, repoErrs.ErrNotFound)
			},
			wantErr: true,
			err:     svcErrs.ErrInvalidCredentials,
		},
		{
			name: "get user: some error",
			args: args{
				ctx: context.Background(),
				input: GenerateTokenInput{
					Email:    "test@example.com",
					Password: "Qwerty!1",
				},
			},
			mockBehavior: func(r *repomocks.MockUser, h *utilmocks.MockPasswordHasher, c *redismocks.MockCache, g *utilmocks.MockTokenGenerator, args args) {
				r.EXPECT().UserByEmail(args.ctx, args.input.Email).Return(entity.User{}, errors.New("some error"))
			},
			wantErr: true,
			err:     svcErrs.ErrCannotGetUser,
		},
		{
			name: "compare passwords: some error",
			args: args{
				ctx: context.Background(),
				input: GenerateTokenInput{
					Email:    "test@example.com",
					Password: "Qwerty!1",
				},
			},
			mockBehavior: func(r *repomocks.MockUser, h *utilmocks.MockPasswordHasher, c *redismocks.MockCache, g *utilmocks.MockTokenGenerator, args args) {
				hash := []byte(args.input.Password)
				user := entity.User{Id: uuid.New(), PasswordHash: hash, Email: args.input.Email}

				r.EXPECT().UserByEmail(args.ctx, args.input.Email).Return(user, nil)
				h.EXPECT().Compare(hash, hash).Return(errors.New("some error"))
			},
			wantErr: true,
			err:     svcErrs.ErrInvalidCredentials,
		},
		{
			name: "generate access token error",
			args: args{
				ctx: context.Background(),
				input: GenerateTokenInput{
					Email:    "test@example.com",
					Password: "Qwerty!1",
				},
			},
			mockBehavior: func(r *repomocks.MockUser, h *utilmocks.MockPasswordHasher, c *redismocks.MockCache, g *utilmocks.MockTokenGenerator, args args) {
				hash := []byte(args.input.Password)
				user := entity.User{Id: uuid.New(), PasswordHash: hash, Email: args.input.Email}

				r.EXPECT().UserByEmail(args.ctx, args.input.Email).Return(user, nil)

				h.EXPECT().Compare(hash, hash).Return(nil)
				g.EXPECT().GenerateAccessToken(user).Return("", errors.New("some error"))
			},
			wantErr: true,
			err:     svcErrs.ErrCannotSignToken,
		},
		{
			name: "generate refresh token error",
			args: args{
				ctx: context.Background(),
				input: GenerateTokenInput{
					Email:    "test@example.com",
					Password: "Qwerty!1",
				},
			},
			mockBehavior: func(r *repomocks.MockUser, h *utilmocks.MockPasswordHasher, c *redismocks.MockCache, g *utilmocks.MockTokenGenerator, args args) {
				hash := []byte(args.input.Password)
				user := entity.User{Id: uuid.New(), PasswordHash: hash, Email: args.input.Email}

				r.EXPECT().UserByEmail(args.ctx, args.input.Email).Return(user, nil)

				h.EXPECT().Compare(hash, hash).Return(nil)
				g.EXPECT().GenerateAccessToken(user).Return("access_token", nil)
				g.EXPECT().GenerateRefreshToken(user).Return("", errors.New("some error"))
			},
			wantErr: true,
			err:     svcErrs.ErrCannotSignToken,
		},
		{
			name: "save to cache: refresh token error",
			args: args{
				ctx: context.Background(),
				input: GenerateTokenInput{
					Email:    "test@example.com",
					Password: "Qwerty!1",
				},
			},
			mockBehavior: func(r *repomocks.MockUser, h *utilmocks.MockPasswordHasher, c *redismocks.MockCache, g *utilmocks.MockTokenGenerator, args args) {
				hash := []byte(args.input.Password)
				user := entity.User{Id: uuid.New(), PasswordHash: hash, Email: args.input.Email}

				r.EXPECT().UserByEmail(args.ctx, args.input.Email).Return(user, nil)

				h.EXPECT().Compare(hash, hash).Return(nil)
				g.EXPECT().GenerateAccessToken(user).Return("access_token", nil)
				g.EXPECT().GenerateRefreshToken(user).Return("refresh_token", nil)
				r.EXPECT().UpdateLastLoginAttempt(args.ctx, user.Id).Return(errors.New("some update error"))
				c.EXPECT().Set(args.ctx, "refresh:"+user.Id.String(), "refresh_token", refreshTokenTTL).Return(errors.New("some error"))
			},
			wantErr: true,
			err:     svcErrs.ErrAccessToCache,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// init deps
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// init repo mock
			repo := repomocks.NewMockUser(ctrl)
			hasher := utilmocks.NewMockPasswordHasher(ctrl)
			cache := redismocks.NewMockCache(ctrl)
			tokenGenerator := utilmocks.NewMockTokenGenerator(ctrl)

			tc.mockBehavior(repo, hasher, cache, tokenGenerator, tc.args)

			// Log
			log := logger.New("local", "info")

			// init service
			s := New(log, cache, repo, hasher, tokenGenerator, refreshTokenTTL)

			// run test
			got, err := s.GenerateToken(tc.args.ctx, tc.args.input)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.NotEqual(t, "", got.AccessToken)
			assert.NotEqual(t, "", got.RefreshToken)
		})
	}
}

func TestAuthService_ParseToken(t *testing.T) {
	type args struct {
		ctx   context.Context
		token string
	}

	type MockBehavior func(g *utilmocks.MockTokenGenerator, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		wantErr      bool
		err          error
	}{
		{
			name: "OK",
			args: args{
				ctx:   context.Background(),
				token: "valid_access_token",
			},
			mockBehavior: func(g *utilmocks.MockTokenGenerator, args args) {
				g.EXPECT().ParseAccessToken(args.token).Return(&jwtgen.Claims{}, nil)
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "some error",
			args: args{
				ctx:   context.Background(),
				token: "invalid_access_token",
			},
			mockBehavior: func(g *utilmocks.MockTokenGenerator, args args) {
				g.EXPECT().ParseAccessToken(args.token).Return(nil, errors.New("some error"))
			},
			wantErr: true,
			err:     svcErrs.ErrCannotParseToken,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// init deps
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// init repo mock
			repo := repomocks.NewMockUser(ctrl)
			hasher := utilmocks.NewMockPasswordHasher(ctrl)
			cache := redismocks.NewMockCache(ctrl)
			tokenGenerator := utilmocks.NewMockTokenGenerator(ctrl)

			tc.mockBehavior(tokenGenerator, tc.args)

			// Log
			log := logger.New("local", "info")

			// init service
			s := New(log, cache, repo, hasher, tokenGenerator, refreshTokenTTL)

			// run test
			got, err := s.ParseToken(tc.args.token)
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)
		})
	}
}
