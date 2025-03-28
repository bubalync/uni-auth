package persistent

import (
	"context"
	"fmt"
	"github.com/bubalync/uni-auth/internal/entity"
	"github.com/bubalync/uni-auth/pkg/postgres"
	"github.com/google/uuid"
)

type Repo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *Repo {
	return &Repo{pg}
}

func (r *Repo) Create(ctx context.Context, u *entity.User) error {
	const op = "repo.persistent.user.Create"

	sql, args, err := r.Builder.
		Insert("users").
		Columns("id, email, password_hash").
		Values(u.ID, u.Email, u.PasswordHash).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: r.Builder: %w", op, err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("%s: r.Pool.Exec: %w", op, err)
	}

	return nil
}

func (r *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (r *Repo) ResetPassword(ctx context.Context, passwordHash []byte) error {
	//TODO implement me
	panic("implement me")
}

func (r *Repo) Update(ctx context.Context, u *entity.User) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repo) UserByEmail(ctx context.Context, email string) (*entity.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repo) UserByEmailIsExists(ctx context.Context, email string) (*bool, error) {
	const op = "repo.persistent.user.UserByEmailIsExists"

	sql, _, err := r.Builder.
		Select("1").
		Prefix("SELECT EXISTS (").
		From("users").
		Where("LOWER(email) = LOWER(?)").
		Suffix(")").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: r.Builder: %w", op, err)
	}

	var isExists bool

	err = r.Pool.QueryRow(ctx, sql, email).Scan(&isExists)
	if err != nil {
		return nil, fmt.Errorf("%s: r.Pool.QueryRow: %w", op, err)
	}

	return &isExists, nil
}

func (r *Repo) UserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	//TODO implement me
	panic("implement me")
}
