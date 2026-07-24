package entity

import (
	"time"

	"github.com/google/uuid"
)

type CommunityReportVoteValue string

const (
	VoteConfirm CommunityReportVoteValue = "confirm"
	VoteReject  CommunityReportVoteValue = "reject"
)

type CommunityReportVote struct {
	CommunityReportID *uuid.UUID
	UserID            *uuid.UUID
	Vote              CommunityReportVoteValue
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
