package usecase

import (
	"context"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/google/uuid"
)

type ICommunityReportsService interface {
	Create(ctx context.Context, report entity.CommunityReport) (entity.CommunityReport, error)
	List(ctx context.Context, id uuid.UUID) ([]entity.UserCommunityReport, error)
	GetByID(ctx context.Context, id uuid.UUID) (entity.CommunityReport, error)
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
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

func (s *communityReportService) List(ctx context.Context, id uuid.UUID) ([]entity.UserCommunityReport, error) {
	reports, err := s.repo.List(ctx, id)
	if err != nil {
		return []entity.UserCommunityReport{}, err
	}

	return reports, nil
}

func (s *communityReportService) GetByID(ctx context.Context, id uuid.UUID) (entity.CommunityReport, error) {
	report, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return entity.CommunityReport{}, err
	}

	return report, nil
}

func (s *communityReportService) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	user, err := s.userSvc.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if !user.IsAdmin {
		return entity.NewForbidden("only an administrator can delete reports")
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
