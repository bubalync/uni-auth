package persistent

import (
	"context"
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
	//TODO implement me
	panic("implement me")
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

func (r *Repo) UserByEmailIsExists(ctx context.Context, email string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repo) UserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	//TODO implement me
	panic("implement me")
}
