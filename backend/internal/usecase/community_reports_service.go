package usecase

import (
	"context"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/google/uuid"
)

type ICommunityReportsService interface {
	Create(ctx context.Context, report entity.CommunityReport) (entity.CommunityReport, error)
	List(ctx context.Context, id uuid.UUID) ([]entity.CommunityReport, error)
}

type communityReportService struct {
	repo    ICommunityReportsRepository
	userSvc IUserService
}

func NewCommunityReportsService(repo ICommunityReportsRepository, userSvc IUserService) ICommunityReportsService {
	return &communityReportService{
		repo:    repo,
		userSvc: userSvc,
	}
}

func (s *communityReportService) Create(ctx context.Context, report entity.CommunityReport) (entity.CommunityReport, error) {
	newReport, err := s.repo.Create(ctx, report)
	if err != nil {
		return entity.CommunityReport{}, err
	}

	return newReport, nil
}

func (s *communityReportService) List(ctx context.Context, id uuid.UUID) ([]entity.CommunityReport, error) {
	reports, err := s.repo.List(ctx, id)
	if err != nil {
		return []entity.CommunityReport{}, err
	}

	return reports, nil
}
