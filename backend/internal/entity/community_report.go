package entity

import (
	"time"

	"github.com/google/uuid"
)

type CommunityReport struct {
	ID                 uuid.UUID
	UserID             *uuid.UUID
	ReturnPointID      *uuid.UUID
	Status             PointStatus
	Description        string
	ConfirmationsCount int
	RejectionsCount    int
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type UserCommunityReport struct {
	Report     CommunityReport
	Username   string
	Reputation int32
}
