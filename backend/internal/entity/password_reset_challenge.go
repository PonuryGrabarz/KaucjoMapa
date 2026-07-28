package entity

import (
	"time"

	"github.com/google/uuid"
)

type PasswordResetChallenge struct {
	ID                  uuid.UUID
	UserID              *uuid.UUID
	CodeHash            string
	Attempts            int
	ExpiresAt           time.Time
	ResendAvailableAt   time.Time
	ResetTokenHash      string
	ResetTokenExpiresAt *time.Time
	ConsumedAt          *time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
