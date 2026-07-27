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
		return entity.CommunityReport{}, mapError(err)
	}

	return toCommunityReportEntity(row), nil
}

func (r *CommunityReportsRepo) GetByID(ctx context.Context, id uuid.UUID) (entity.CommunityReport, error) {
	row, err := r.q.GetCommunityReportById(ctx, id)
	if err != nil {
		return entity.CommunityReport{}, mapError(err)
	}

	return toCommunityReportEntity(row), nil
}

func (r *CommunityReportsRepo) List(ctx context.Context, id uuid.UUID) ([]entity.UserCommunityReport, error) {
	rows, err := r.q.ListCommunityReportsByReturnPoint(ctx, id)
	if mapped := mapError(err); mapped != nil {
		// Brak zgłoszeń dla punktu to pusta lista, nie błąd.
		if entity.IsKind(mapped, entity.KindNotFound) {
			return []entity.UserCommunityReport{}, nil
		}
		return []entity.UserCommunityReport{}, mapped
	}

	userReports := toListCommunityReportsByReturnPoint(rows)

	return userReports, nil
}

func (r *CommunityReportsRepo) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeleteCommunityReportById(ctx, id)
	if err != nil {
		return mapError(err)
	}

	return nil
}
