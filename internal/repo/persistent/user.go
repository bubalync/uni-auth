package persistent

import (
	"context"
	"errors"
	"fmt"
	"github.com/bubalync/uni-auth/internal/entity"
	"github.com/bubalync/uni-auth/internal/repo/repoErrs"
	"github.com/bubalync/uni-auth/pkg/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (r *UserRepo) Create(ctx context.Context, u entity.User) error {
	const op = "repo.persistent.user.Create"

	sql, args, _ := r.Builder.
		Insert("users").
		Columns("id, email, password_hash").
		Values(u.ID, u.Email, u.PasswordHash).
		ToSql()

	_, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.ConstraintName == "users_email_lower_unique" {
				return repoErrs.ErrAlreadyExists
			}
		}

		return fmt.Errorf("%s: r.Pool.Exec: %w", op, err)
	}

	return nil
}

func (r *UserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (r *UserRepo) ResetPassword(ctx context.Context, passwordHash []byte) error {
	//TODO implement me
	panic("implement me")
}

func (r *UserRepo) Update(ctx context.Context, u entity.User) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (r *UserRepo) UserByEmail(ctx context.Context, email string) (entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *UserRepo) UserByEmailIsExists(ctx context.Context, email string) (*bool, error) {
	const op = "repo.persistent.user.UserByEmailIsExists"

	sql, _, _ := r.Builder.
		Select("1").
		Prefix("SELECT EXISTS (").
		From("users").
		Where("LOWER(email) = LOWER(?)").
		Suffix(")").
		ToSql()

	var isExists bool

	err := r.Pool.QueryRow(ctx, sql, email).Scan(&isExists)
	if err != nil {
		return nil, fmt.Errorf("%s: r.Pool.QueryRow: %w", op, err)
	}

	return &isExists, nil
}

func (r *UserRepo) UserByID(ctx context.Context, id uuid.UUID) (entity.User, error) {
	//TODO implement me
	panic("implement me")
}
