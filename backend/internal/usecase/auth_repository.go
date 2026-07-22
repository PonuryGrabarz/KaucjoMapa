package usecase

import (
	"context"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/google/uuid"
)

type IRefreshTokenRepository interface {
	Create(ctx context.Context, token entity.RefreshToken) (entity.RefreshToken, error)
	GetTokensByUserID(ctx context.Context, id uuid.UUID) ([]entity.RefreshToken, error)
	GetToken(ctx context.Context, id uuid.UUID, tokenString string) (entity.RefreshToken, error)
	RevokeToken(ctx context.Context, id uuid.UUID, tokenString string) error
}
