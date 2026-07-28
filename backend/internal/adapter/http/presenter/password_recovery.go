package presenter

import (
	"time"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/google/uuid"
)

type StartRecoveryRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type VerifyCodeRequest struct {
	Code string `json:"code" validate:"required,len=8,numeric"`
}

type ResetPasswordRequest struct {
	ResetToken  string `json:"reset_token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type RecoveryChallengeResponse struct {
	ChallengeID       uuid.UUID `json:"challenge_id"`
	ExpiresAt         time.Time `json:"expires_at"`
	ResendAvailableAt time.Time `json:"resend_available_at"`
}

type RecoveryChallengeEnvelope struct {
	Data RecoveryChallengeResponse `json:"data"`
}

type ResetTokenResponse struct {
	ResetToken string    `json:"reset_token"`
	ExpiresAt  time.Time `json:"expires_at"`
}

type ResetTokenEnvelope struct {
	Data ResetTokenResponse `json:"data"`
}

func RecoveryChallengeSuccessResponse(challenge entity.PasswordResetChallenge) RecoveryChallengeEnvelope {
	return RecoveryChallengeEnvelope{
		Data: RecoveryChallengeResponse{
			ChallengeID:       challenge.ID,
			ExpiresAt:         challenge.ExpiresAt,
			ResendAvailableAt: challenge.ResendAvailableAt,
		},
	}
}

func ResetTokenSuccessResponse(token string, expiresAt time.Time) ResetTokenEnvelope {
	return ResetTokenEnvelope{
		Data: ResetTokenResponse{
			ResetToken: token,
			ExpiresAt:  expiresAt,
		},
	}
}
