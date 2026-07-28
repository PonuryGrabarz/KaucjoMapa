package usecase

import (
	"context"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/google/uuid"
)

type IUserRepository interface {
	Create(ctx context.Context, user entity.User) (entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (entity.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error
}
