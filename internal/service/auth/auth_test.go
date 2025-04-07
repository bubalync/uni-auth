package auth

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/bubalync/uni-auth/internal/entity"
	"github.com/bubalync/uni-auth/internal/mocks/repomocks"
	"github.com/bubalync/uni-auth/internal/mocks/utilmocks"
	"github.com/bubalync/uni-auth/internal/repo/repoErrs"
	"github.com/bubalync/uni-auth/internal/service/svcErrs"
	"github.com/bubalync/uni-auth/pkg/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
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
			err:     svcErrs.ErrInternal,
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
			err:     svcErrs.ErrInternal,
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
			s := New(log, repo, hasher, "sign_key", 0)

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
