package usecase

import (
	"context"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/google/uuid"
)

type IDataErrorReportsService interface {
	Create(ctx context.Context, report entity.ErrorReport) (entity.ErrorReport, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]entity.UserErrorReport, error)
	ChangeStatus(ctx context.Context, id uuid.UUID, status entity.ErrorReportStatus) error
}

type DataErrorReportsService struct {
	repo IDataErrorReportsRepository
}

func NewDataErrorReportsService(repo IDataErrorReportsRepository) *DataErrorReportsService {
	return &DataErrorReportsService{
		repo: repo,
	}
}

func (s *DataErrorReportsService) Create(ctx context.Context, report entity.ErrorReport) (entity.ErrorReport, error) {
	newReport, err := s.repo.Create(ctx, report)
	if err != nil {
		return entity.ErrorReport{}, err
	}

	return newReport, nil
}

func (s *DataErrorReportsService) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *DataErrorReportsService) List(ctx context.Context) ([]entity.UserErrorReport, error) {
	list, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (s *DataErrorReportsService) ChangeStatus(ctx context.Context, id uuid.UUID, status entity.ErrorReportStatus) error {
	report, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	report.Report.Status = status

	_, err = s.repo.Update(ctx, report.Report)
	if err != nil {
		return err
	}

	return nil
}
