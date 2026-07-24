package usecase

import (
	"context"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
)

type ICommunityReportVoteRepository interface {
	Create(ctx context.Context, vote entity.CommunityReportVote) (entity.CommunityReportVote, error)
}
