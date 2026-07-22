package postgres

import (
	"context"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/infrastructure/database"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase"
	"github.com/google/uuid"
)

type RefreshTokenRepo struct {
	q *database.Queries
}

func NewRefreshTokenRepo(q *database.Queries) *RefreshTokenRepo {
	return &RefreshTokenRepo{
		q: q,
	}
}

var _ usecase.IRefreshTokenRepository = (*RefreshTokenRepo)(nil)

func (r *RefreshTokenRepo) Create(ctx context.Context, token entity.RefreshToken) (entity.RefreshToken, error) {
	row, err := r.q.CreateRefreshToken(ctx, toCreateRefreshTokenParams(token))
	if err != nil {
		return entity.RefreshToken{}, entity.ErrServerError
	}

	return toRefreshTokenEntity(row), nil
}

func (r *RefreshTokenRepo) GetTokensByUserID(ctx context.Context, id uuid.UUID) ([]entity.RefreshToken, error) {
	row, err := r.q.GetUserTokens(ctx, id)
	if err != nil {
		return []entity.RefreshToken{}, entity.ErrServerError
	}

	var userTokens []entity.RefreshToken

	for _, token := range row {
		userTokens = append(userTokens, toRefreshTokenEntity(token))
	}

	return userTokens, nil
}

func (r *RefreshTokenRepo) GetToken(ctx context.Context, id uuid.UUID, tokenString string) (entity.RefreshToken, error) {
	row, err := r.q.GetUserRefreshToken(ctx, toGetUserRefresTokenParams(id, tokenString))
	if err != nil {
		return entity.RefreshToken{}, entity.ErrNotFound
	}

	return toRefreshTokenEntity(row), nil
}

func (r *RefreshTokenRepo) RevokeToken(ctx context.Context, id uuid.UUID, tokenString string) error {
	_, err := r.GetToken(ctx, id, tokenString)
	if err != nil {
		return entity.ErrUnauthorized
	}
	err = r.q.RevokeToken(ctx, tokenString)
	if err != nil {
		return entity.ErrNotFound
	}

	return nil
}
