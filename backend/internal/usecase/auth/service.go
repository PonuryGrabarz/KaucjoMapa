package auth

import (
	"context"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/google/uuid"
)

type IRefreshTokenService interface {
	Create(ctx context.Context, token entity.RefreshToken) (entity.RefreshToken, error)
	GetTokensByUserID(ctx context.Context, id uuid.UUID) ([]entity.RefreshToken, error)
	GetToken(ctx context.Context, id uuid.UUID, tokenString string) (entity.RefreshToken, error)
	RevokeToken(ctx context.Context, id uuid.UUID, tokenString string) error
}

type refreshTokenService struct {
	refreshTokenRepo IRefreshTokenRepository
}

func NewRefreshTokenService(refreshTokenRepo IRefreshTokenRepository) IRefreshTokenService {
	return &refreshTokenService{
		refreshTokenRepo: refreshTokenRepo,
	}
}

func (s *refreshTokenService) Create(ctx context.Context, token entity.RefreshToken) (entity.RefreshToken, error) {
	token.ID = uuid.New()

	newToken, err := s.refreshTokenRepo.Create(ctx, token)
	if err != nil {
		return entity.RefreshToken{}, err
	}

	return newToken, err
}

func (s *refreshTokenService) GetTokensByUserID(ctx context.Context, id uuid.UUID) ([]entity.RefreshToken, error) {
	tokens, err := s.refreshTokenRepo.GetTokensByUserID(ctx, id)
	if err != nil {
		return []entity.RefreshToken{}, err
	}

	return tokens, nil
}

func (s *refreshTokenService) GetToken(ctx context.Context, id uuid.UUID, tokenString string) (entity.RefreshToken, error) {
	token, err := s.refreshTokenRepo.GetToken(ctx, id, tokenString)
	if err != nil {
		return entity.RefreshToken{}, err
	}

	return token, nil
}

func (s *refreshTokenService) RevokeToken(ctx context.Context, id uuid.UUID, tokenString string) error {
	err := s.refreshTokenRepo.RevokeToken(ctx, id, tokenString)
	if err != nil {
		return err
	}

	return nil
}
