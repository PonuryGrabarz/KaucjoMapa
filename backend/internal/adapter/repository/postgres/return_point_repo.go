package postgres

import (
	"context"
	"errors"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/infrastructure/database"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/returnpoint"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
)

type ReturnPointRepo struct {
	q *database.Queries
}

func NewReturnPointRepo(q *database.Queries) *ReturnPointRepo {
	return &ReturnPointRepo{q: q}
}

// Kompilacyjna gwarancja, że spełniamy port.
var _ returnpoint.IReturnPointRepository = (*ReturnPointRepo)(nil)

func (r *ReturnPointRepo) ExistsNearby(ctx context.Context, lat, lng float64) (bool, error) {
	row, err := r.q.GetPointByDistance(ctx, database.GetPointByDistanceParams{
		LlToEarth: lat, LlToEarth_2: lng,
	})

	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}

	if row.ID != uuid.Nil {
		return true, entity.ErrPointExistsNearby
	}

	return false, err
}

func (r *ReturnPointRepo) Create(ctx context.Context, p entity.ReturnPoint) (entity.ReturnPoint, error) {
	created, err := r.q.CreatePoint(ctx, toCreateReturnPointParams(p)) // mapper.go: entity -> database
	if err != nil {
		return entity.ReturnPoint{}, err
	}

	return toReturnPointEntity(created), nil // mapper.go: database -> entity
}

func (r *ReturnPointRepo) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeletePointById(ctx, id)
	if err != nil {
		return entity.ErrNotFound
	}

	return nil
}

func (r *ReturnPointRepo) GetByID(ctx context.Context, id uuid.UUID) (entity.ReturnPoint, error) {
	row, err := r.q.GetPointById(ctx, id)
	if err != nil {
		return entity.ReturnPoint{}, entity.ErrNotFound
	}

	return toReturnPointEntity(row), nil
}

func (r *ReturnPointRepo) List(ctx context.Context, f returnpoint.Filter) ([]entity.ReturnPoint, error) {
	row, err := r.q.ListReturnPoints(ctx, toListReturnPointsParams(f))
	if err != nil {
		return []entity.ReturnPoint{}, entity.ErrNotFound
	}

	var returnPoints []entity.ReturnPoint
	for _, point := range row {
		returnPoints = append(returnPoints, toReturnPointEntity(point))
	}

	return returnPoints, nil
}

func (r *ReturnPointRepo) Update(ctx context.Context, p entity.ReturnPoint) (entity.ReturnPoint, error) {
	row, err := r.q.UpdatePointById(ctx, toUpdatePointParams(p))
	if err != nil {
		return entity.ReturnPoint{}, entity.ErrNotFound
	}

	return toReturnPointEntity(row), nil
}
