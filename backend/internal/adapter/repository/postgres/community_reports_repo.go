package postgres

import (
	"context"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/infrastructure/database"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase"
	"github.com/google/uuid"
)

type CommunityReportsRepo struct {
	q *database.Queries
}

func NewCommunityReportsRepo(q *database.Queries) *CommunityReportsRepo {
	return &CommunityReportsRepo{
		q: q,
	}
}

var _ usecase.ICommunityReportsRepository = (*CommunityReportsRepo)(nil)

func (r *CommunityReportsRepo) Create(ctx context.Context, report entity.CommunityReport) (entity.CommunityReport, error) {
	row, err := r.q.CreateCommunityReport(ctx, toCreateCommunityReportParams(report))
	if err != nil {
		return entity.CommunityReport{}, entity.ErrServerError
	}

	return toCommunityReportEntity(row), nil
}

func (r *CommunityReportsRepo) GetByID(ctx context.Context, id uuid.UUID) (entity.CommunityReport, error) {
	row, err := r.q.GetCommunityReportById(ctx, id)
	if err != nil {
		return entity.CommunityReport{}, entity.ErrNotFound
	}

	return toCommunityReportEntity(row), nil
}

func (r *CommunityReportsRepo) List(ctx context.Context, id uuid.UUID) ([]entity.CommunityReport, error) {
	rows, err := r.q.ListCommunityReportsByReturnPoint(ctx, id)
	if err != nil {
		return []entity.CommunityReport{}, entity.ErrNotFound
	}

	reports := []entity.CommunityReport{}
	for _, row := range rows {
		reports = append(reports, toCommunityReportEntity(row))
	}

	return reports, nil
}
