package dto

import "github.com/google/uuid"

type StartRecoveryInput struct {
	Email string
}

type VerifyCodeInput struct {
	ChallengeID uuid.UUID
	Code        string
}

type ResendCodeInput struct {
	ChallengeID uuid.UUID
}

type ResetPasswordInput struct {
	ChallengeID uuid.UUID
	ResetToken  string
	NewPassword string
}
