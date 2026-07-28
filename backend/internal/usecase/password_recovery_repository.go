package usecase

import (
	"context"
	"time"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/google/uuid"
)

type IPasswordRecoveryRepository interface {
	Create(ctx context.Context, challenge entity.PasswordResetChallenge) (entity.PasswordResetChallenge, error)
	GetByID(ctx context.Context, id uuid.UUID) (entity.PasswordResetChallenge, error)
	ConsumeActiveByUserID(ctx context.Context, userID uuid.UUID) error
	IncrementAttempts(ctx context.Context, id uuid.UUID) (entity.PasswordResetChallenge, error)
	UpdateCode(ctx context.Context, challenge entity.PasswordResetChallenge) (entity.PasswordResetChallenge, error)
	SetResetToken(ctx context.Context, id uuid.UUID, tokenHash string, expiresAt time.Time) (entity.PasswordResetChallenge, error)
	Consume(ctx context.Context, id uuid.UUID) error
}
