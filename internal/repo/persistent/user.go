package persistent

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/bubalync/uni-auth/internal/entity"
	"github.com/bubalync/uni-auth/internal/repo/repoErrs"
	"github.com/bubalync/uni-auth/pkg/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
		Values(u.Id, u.Email, u.PasswordHash).
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

func (r *UserRepo) Update(ctx context.Context, u entity.User) error {
	//TODO implement me
	panic("implement me")
}

func (r *UserRepo) UpdateLastLoginAttempt(ctx context.Context, id uuid.UUID) error {
	const op = "repo.persistent.user.UpdateLastLoginAttempt"

	sql, args, _ := r.Builder.
		Update("users").
		Set("last_login_attempt", squirrel.Expr("NOW()")).
		Where("id = ?", id).
		ToSql()

	_, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("%s: r.Pool.Exec: %w", op, err)
	}

	return nil
}

func (r *UserRepo) UserByEmail(ctx context.Context, email string) (entity.User, error) {
	const op = "repo.persistent.user.UserByEmail"

	sql, args, _ := r.Builder.
		Select("id, email, password_hash, name, is_active, last_login_attempt, created_at, updated_at").
		From("users").
		Where("LOWER(email) = LOWER(?)", email).
		ToSql()

	var user entity.User
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&user.Id,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.IsActive,
		&user.LastLoginAttempt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, repoErrs.ErrNotFound
		}
		return entity.User{}, fmt.Errorf("%s: r.Pool.QueryRow: %w", op, err)
	}

	return user, nil
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
