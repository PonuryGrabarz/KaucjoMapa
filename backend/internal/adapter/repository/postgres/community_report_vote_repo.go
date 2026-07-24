package postgres

import (
	"context"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/infrastructure/database"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase"
)

type CommunityReportVoteRepo struct {
	q *database.Queries
}

func NewCommunityReportVoteRepo(q *database.Queries) *CommunityReportVoteRepo {
	return &CommunityReportVoteRepo{
		q: q,
	}
}

var _ usecase.ICommunityReportVoteRepository = (*CommunityReportVoteRepo)(nil)

func (r *CommunityReportVoteRepo) Create(ctx context.Context, vote entity.CommunityReportVote) (entity.CommunityReportVote, error) {
	reportVote, err := r.q.UpsertCommunityReportVote(ctx, toCreateCommunityReportVoteParams(vote))
	if err != nil {
		return entity.CommunityReportVote{}, entity.ErrServerError
	}

	return toCommunityReportVoteEntity(reportVote), nil
}
