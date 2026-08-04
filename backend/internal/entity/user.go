package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                 uuid.UUID
	Email              string
	Username           string
	PasswordHash       string
	Token              string
	Reputation         int32
	AvatarPath         string
	TotalConfirmations int32
	TotalRejections    int32
	IsAdmin            bool
	IsBanned           bool
	GoogleID           string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
