package usecase

import (
	"context"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
)

type IGoogleOAuthRepository interface {
	GetUserByGoogleID(ctx context.Context, googleID string) (entity.User, error)
	CreateUser(ctx context.Context, user entity.User) (entity.User, error)
}
