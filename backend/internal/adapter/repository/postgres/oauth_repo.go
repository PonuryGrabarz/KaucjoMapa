package postgres

import (
	"context"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/infrastructure/database"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase"
)

type OAuthRepo struct {
	q *database.Queries
}

func NewOAuthRepo(q *database.Queries) *OAuthRepo {
	return &OAuthRepo{
		q: q,
	}
}

var _ usecase.IGoogleOAuthRepository = (*OAuthRepo)(nil)

func (r *OAuthRepo) GetUserByGoogleID(ctx context.Context, googleID string) (entity.User, error) {
	row, err := r.q.GetUserByGoogleSubject(ctx, toNullString(googleID))
	if err != nil {
		return entity.User{}, mapError(err)
	}

	return toUserEntity(row), nil
}

func (r *OAuthRepo) CreateUser(ctx context.Context, user entity.User) (entity.User, error) {
	row, err := r.q.UpsertGoogleUserByEmail(ctx, toGoogleUserParams(user))
	if err != nil {
		return entity.User{}, mapError(err)
	}

	return toUserEntity(row), nil
}
