package postgres

import (
	"context"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/infrastructure/database"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase"
	"github.com/google/uuid"
)

type UserRepo struct {
	q *database.Queries
}

func NewUserRepo(q *database.Queries) *UserRepo {
	return &UserRepo{
		q: q,
	}
}

var _ usecase.IUserRepository = (*UserRepo)(nil)

func (r *UserRepo) Create(ctx context.Context, user entity.User) (entity.User, error) {
	row, err := r.q.CreateUser(ctx, toCreateUserParams(user))
	if err != nil {
		return entity.User{}, entity.ErrServerError
	}

	return toUserEntity(row), nil
}

func (r *UserRepo) GetUserByUsername(ctx context.Context, username string) (entity.User, error) {
	row, err := r.q.GetUserByUsername(ctx, username)
	if err != nil {
		return entity.User{}, entity.ErrNotFound
	}

	return toUserEntity(row), nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, id uuid.UUID) (entity.User, error) {
	row, err := r.q.GetUserById(ctx, id)
	if err != nil {
		return entity.User{}, entity.ErrNotFound
	}

	return toUserEntity(row), nil
}
