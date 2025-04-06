package persistent

import (
	"context"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/bubalync/uni-auth/internal/entity"
	"github.com/bubalync/uni-auth/pkg/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

func boolPointer(b bool) *bool {
	return &b
}

func TestUserRepo_Create(t *testing.T) {
	type args struct {
		ctx  context.Context
		user entity.User
	}

	type MockBehavior func(m pgxmock.PgxPoolIface, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				user: entity.User{
					ID:           uuid.New(),
					Email:        "test@example.com",
					PasswordHash: make([]byte, 1),
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectExec("INSERT INTO users").
					WithArgs(args.user.ID, args.user.Email, args.user.PasswordHash).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			wantErr: false,
		},
		{
			name: "user already exists",
			args: args{
				ctx: context.Background(),
				user: entity.User{
					Email:        "test@example.com",
					PasswordHash: make([]byte, 1),
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectExec("INSERT INTO users").
					WithArgs(args.user.ID, args.user.Email, args.user.PasswordHash).
					WillReturnError(&pgconn.PgError{
						ConstraintName: "users_email_lower_unique",
					})
			},
			wantErr: true,
		},
		{
			name: "unexpected error",
			args: args{
				ctx: context.Background(),
				user: entity.User{
					Email:        "test@example.com",
					PasswordHash: make([]byte, 1),
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectQuery("INSERT INTO users").
					WithArgs(args.user.ID, args.user.Email, args.user.PasswordHash).
					WillReturnError(errors.New("some error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			poolMock, _ := pgxmock.NewPool()
			defer poolMock.Close()
			tc.mockBehavior(poolMock, tc.args)

			postgresMock := &postgres.Postgres{
				Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
				Pool:    poolMock,
			}
			userRepoMock := NewUserRepo(postgresMock)

			err := userRepoMock.Create(tc.args.ctx, tc.args.user)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			err = poolMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestUserRepo_UserByEmailIsExists(t *testing.T) {
	type args struct {
		ctx   context.Context
		email string
	}

	type MockBehavior func(m pgxmock.PgxPoolIface, args args)

	var testCases = []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         *bool
		wantErr      bool
	}{
		{
			name: "OK true",
			args: args{
				ctx:   context.Background(),
				email: "test@example.com",
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				rows := pgxmock.NewRows([]string{"exists"}).
					AddRow(false)

				m.ExpectQuery("SELECT 1 FROM users").
					WithArgs(args.email).
					WillReturnRows(rows)
			},
			want:    boolPointer(true),
			wantErr: false,
		},
		{
			name: "OK false",
			args: args{
				ctx:   context.Background(),
				email: "test@example.com",
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				rows := pgxmock.NewRows([]string{"exists"}).
					AddRow(false)

				m.ExpectQuery("SELECT 1 FROM users").
					WithArgs(args.email).
					WillReturnRows(rows)
			},
			want:    boolPointer(false),
			wantErr: false,
		},
		{
			name: "unexpected error",
			args: args{
				ctx:   context.Background(),
				email: "test@example.com",
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectQuery("SELECT 1 FROM users").
					WithArgs(args.email).
					WillReturnError(errors.New("some error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			poolMock, _ := pgxmock.NewPool()
			defer poolMock.Close()
			tc.mockBehavior(poolMock, tc.args)

			postgresMock := &postgres.Postgres{
				Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
				Pool:    poolMock,
			}
			userRepoMock := NewUserRepo(postgresMock)

			got, err := userRepoMock.UserByEmailIsExists(tc.args.ctx, tc.args.email)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)

			err = poolMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
