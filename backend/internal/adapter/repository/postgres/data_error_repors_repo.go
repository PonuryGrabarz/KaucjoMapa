package postgres

import (
	"context"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/infrastructure/database"
	"github.com/google/uuid"
)

type DataErrorReportsRepo struct {
	q *database.Queries
}

func NewDataErrorReportsRepo(q *database.Queries) *DataErrorReportsRepo {
	return &DataErrorReportsRepo{
		q: q,
	}
}

func (r *DataErrorReportsRepo) Create(ctx context.Context, report entity.ErrorReport) (entity.ErrorReport, error) {
	row, err := r.q.CreateDataErrorReport(ctx, toCreateDataErrorReportParams(report))
	if err != nil {
		return entity.ErrorReport{}, mapError(err)
	}

	return toErrorReportEntity(row), nil
}

func (r *DataErrorReportsRepo) GetByID(ctx context.Context, id uuid.UUID) (entity.UserErrorReport, error) {
	row, err := r.q.GetDataErrorReportById(ctx, id)
	if err != nil {
		return entity.UserErrorReport{}, mapError(err)
	}

	return toUserErrorReportEntity(row), nil
}

func (r *DataErrorReportsRepo) List(ctx context.Context) ([]entity.UserErrorReport, error) {
	rows, err := r.q.ListDataErrorReports(ctx, database.NullDataErrorReportStatus{
		Valid: false,
	})

	if err != nil {
		return []entity.UserErrorReport{}, err
	}

	list := []entity.UserErrorReport{}

	for _, row := range rows {
		list = append(list, toUserErrorReportEntity(database.GetDataErrorReportByIdRow{
			ID:                  row.ID,
			UserID:              row.UserID,
			ReturnPointID:       row.ReturnPointID,
			Categories:          row.Categories,
			Description:         row.Description,
			SuggestedCorrection: row.SuggestedCorrection,
			Status:              row.Status,
			CreatedAt:           row.CreatedAt,
			UpdatedAt:           row.UpdatedAt,
			Username:            row.Username,
			ReturnPointName:     row.ReturnPointName,
		}))
	}

	return list, nil
}

func (r *DataErrorReportsRepo) Update(ctx context.Context, report entity.ErrorReport) (entity.ErrorReport, error) {
	row, err := r.q.UpdateDataErrorReportStatus(ctx, database.UpdateDataErrorReportStatusParams{
		ID:     report.ID,
		Status: database.DataErrorReportStatus(report.Status),
	})
	if err != nil {
		return entity.ErrorReport{}, err
	}

	return toErrorReportEntity(row), nil
}

func (r *DataErrorReportsRepo) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeleteDataErrorReportById(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
