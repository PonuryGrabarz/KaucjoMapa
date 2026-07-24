package usecase

import (
	"context"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
)

type ICommunityReportVoteService interface {
	Create(ctx context.Context, vote entity.CommunityReportVote) (entity.CommunityReportVote, error)
}

type communityReportVoteService struct {
	repo      ICommunityReportVoteRepository
	userSvc   IUserService
	reportSvc ICommunityReportsService
}

func NewCommunityReportVoteService(repo ICommunityReportVoteRepository, userSvc IUserService, reportSvc ICommunityReportsService) ICommunityReportVoteService {
	return &communityReportVoteService{
		repo:      repo,
		userSvc:   userSvc,
		reportSvc: reportSvc,
	}
}

func (s *communityReportVoteService) Create(ctx context.Context, vote entity.CommunityReportVote) (entity.CommunityReportVote, error) {
	votedReport, err := s.repo.Create(ctx, vote)
	if err != nil {
		return entity.CommunityReportVote{}, err
	}

	return votedReport, nil
}
