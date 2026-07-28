package postgres

import (
	"context"
	"time"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/infrastructure/database"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase"
	"github.com/google/uuid"
)

type PasswordRecoveryRepo struct {
	q *database.Queries
}

func NewPasswordRecoveryRepo(q *database.Queries) *PasswordRecoveryRepo {
	return &PasswordRecoveryRepo{
		q: q,
	}
}

var _ usecase.IPasswordRecoveryRepository = (*PasswordRecoveryRepo)(nil)

func (r *PasswordRecoveryRepo) Create(ctx context.Context, challenge entity.PasswordResetChallenge) (entity.PasswordResetChallenge, error) {
	row, err := r.q.CreatePasswordResetChallenge(ctx, toCreatePasswordResetChallengeParams(challenge))
	if err != nil {
		return entity.PasswordResetChallenge{}, mapError(err)
	}

	return toPasswordResetChallengeEntity(row), nil
}

func (r *PasswordRecoveryRepo) GetByID(ctx context.Context, id uuid.UUID) (entity.PasswordResetChallenge, error) {
	row, err := r.q.GetPasswordResetChallengeByID(ctx, id)
	if err != nil {
		return entity.PasswordResetChallenge{}, mapError(err)
	}

	return toPasswordResetChallengeEntity(row), nil
}

func (r *PasswordRecoveryRepo) ConsumeActiveByUserID(ctx context.Context, userID uuid.UUID) error {
	err := r.q.ConsumeActivePasswordResetChallengesByUserID(ctx, toNullUUID(&userID))
	if err != nil {
		return mapError(err)
	}

	return nil
}

func (r *PasswordRecoveryRepo) IncrementAttempts(ctx context.Context, id uuid.UUID) (entity.PasswordResetChallenge, error) {
	row, err := r.q.IncrementPasswordResetChallengeAttempts(ctx, id)
	if err != nil {
		return entity.PasswordResetChallenge{}, mapError(err)
	}

	return toPasswordResetChallengeEntity(row), nil
}

func (r *PasswordRecoveryRepo) UpdateCode(ctx context.Context, challenge entity.PasswordResetChallenge) (entity.PasswordResetChallenge, error) {
	row, err := r.q.UpdatePasswordResetChallengeCode(ctx, toUpdatePasswordResetChallengeCodeParams(challenge))
	if err != nil {
		return entity.PasswordResetChallenge{}, mapError(err)
	}

	return toPasswordResetChallengeEntity(row), nil
}

func (r *PasswordRecoveryRepo) SetResetToken(ctx context.Context, id uuid.UUID, tokenHash string, expiresAt time.Time) (entity.PasswordResetChallenge, error) {
	row, err := r.q.SetPasswordResetToken(ctx, toSetPasswordResetTokenParams(id, tokenHash, expiresAt))
	if err != nil {
		return entity.PasswordResetChallenge{}, mapError(err)
	}

	return toPasswordResetChallengeEntity(row), nil
}

func (r *PasswordRecoveryRepo) Consume(ctx context.Context, id uuid.UUID) error {
	err := r.q.ConsumePasswordResetChallenge(ctx, id)
	if err != nil {
		return mapError(err)
	}

	return nil
}
