package postgres

import (
	"context"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/infrastructure/database"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase"
	"github.com/google/uuid"
)

type ReturnPlanRepo struct {
	q *database.Queries
}

func NewReturnPlanRepo(q *database.Queries) *ReturnPlanRepo {
	return &ReturnPlanRepo{
		q: q,
	}
}

var _ usecase.IReturnPlanRepository = (*ReturnPlanRepo)(nil)

func (r *ReturnPlanRepo) CreateReturnPlan(ctx context.Context, plan entity.ReturnPlan) (entity.ReturnPlan, error) {
	row, err := r.q.CreateReturnPlan(ctx, toCreateReturnPlanParams(plan))
	if err != nil {
		return entity.ReturnPlan{}, mapError(err)
	}

	return toReturnPlanEntity(row), nil
}

func (r *ReturnPlanRepo) GetReturnPlanByID(ctx context.Context, id uuid.UUID) (entity.ReturnPlan, error) {
	row, err := r.q.GetReturnPlanByID(ctx, id)
	if err != nil {
		return entity.ReturnPlan{}, mapError(err)
	}

	return toReturnPlanEntity(row), nil
}

func (r *ReturnPlanRepo) GetReturnPlanByIDAndUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (entity.ReturnPlan, error) {
	row, err := r.q.GetReturnPlanByIDAndUserID(ctx, database.GetReturnPlanByIDAndUserIDParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		return entity.ReturnPlan{}, mapError(err)
	}

	return toReturnPlanEntity(row), nil
}

func (r *ReturnPlanRepo) GetDraftReturnPlanByUserID(ctx context.Context, userID uuid.UUID) (entity.ReturnPlan, error) {
	row, err := r.q.GetDraftReturnPlanByUserID(ctx, userID)
	if err != nil {
		return entity.ReturnPlan{}, mapError(err)
	}

	return toReturnPlanEntity(row), nil
}

func (r *ReturnPlanRepo) ListReturnPlansByUserID(ctx context.Context, userID uuid.UUID) ([]entity.ReturnPlan, error) {
	rows, err := r.q.ListReturnPlansByUserID(ctx, userID)
	if err != nil {
		return []entity.ReturnPlan{}, mapError(err)
	}

	plans := make([]entity.ReturnPlan, 0, len(rows))
	for _, row := range rows {
		plans = append(plans, toReturnPlanEntity(row))
	}

	return plans, nil
}

func (r *ReturnPlanRepo) UpdateReturnPlanTargets(ctx context.Context, planID uuid.UUID, targetPet int, targetCan int, targetGlass int) (entity.ReturnPlan, error) {
	row, err := r.q.UpdateReturnPlanTargets(ctx, database.UpdateReturnPlanTargetsParams{
		ID:          planID,
		TargetPet:   int32(targetPet),
		TargetCan:   int32(targetCan),
		TargetGlass: int32(targetGlass),
	})
	if err != nil {
		return entity.ReturnPlan{}, mapError(err)
	}

	return toReturnPlanEntity(row), nil
}

func (r *ReturnPlanRepo) UpdateReturnPlanStatus(ctx context.Context, planID uuid.UUID, status entity.PlanStatus) (entity.ReturnPlan, error) {
	row, err := r.q.UpdateReturnPlanStatus(ctx, database.UpdateReturnPlanStatusParams{
		ID:     planID,
		Status: database.ReturnPlanStatus(status),
	})
	if err != nil {
		return entity.ReturnPlan{}, mapError(err)
	}

	return toReturnPlanEntity(row), nil
}

func (r *ReturnPlanRepo) DeleteReturnPlanByIDAndUserID(ctx context.Context, planID uuid.UUID, userID uuid.UUID) error {
	err := r.q.DeleteReturnPlanByIDAndUserID(ctx, database.DeleteReturnPlanByIDAndUserIDParams{
		ID:     planID,
		UserID: userID,
	})
	if err != nil {
		return mapError(err)
	}

	return nil
}

func (r *ReturnPlanRepo) ConfirmReturnPlan(ctx context.Context, planID uuid.UUID) (entity.ReturnPlan, error) {
	row, err := r.q.ConfirmReturnPlan(ctx, planID)
	if err != nil {
		return entity.ReturnPlan{}, mapError(err)
	}

	return toReturnPlanEntity(row), nil
}

func (r *ReturnPlanRepo) CreateReturnPlanPoint(ctx context.Context, point entity.ReturnPlanPoint) (entity.ReturnPlanPoint, error) {
	row, err := r.q.CreateReturnPlanPoint(ctx, toCreateReturnPlanPointParams(point))
	if err != nil {
		return entity.ReturnPlanPoint{}, mapError(err)
	}

	return toReturnPlanPointEntity(row), nil
}

func (r *ReturnPlanRepo) GetReturnPlanPointByID(ctx context.Context, id uuid.UUID) (entity.ReturnPlanPoint, error) {
	row, err := r.q.GetReturnPlanPointByID(ctx, id)
	if err != nil {
		return entity.ReturnPlanPoint{}, mapError(err)
	}

	return toReturnPlanPointEntity(row), nil
}

func (r *ReturnPlanRepo) GetReturnPlanPointByPlanAndPoint(ctx context.Context, planID uuid.UUID, pointID uuid.UUID) (entity.ReturnPlanPoint, error) {
	row, err := r.q.GetReturnPlanPointByPlanAndPoint(ctx, database.GetReturnPlanPointByPlanAndPointParams{
		ReturnPlanID:  planID,
		ReturnPointID: pointID,
	})
	if err != nil {
		return entity.ReturnPlanPoint{}, mapError(err)
	}

	return toReturnPlanPointEntity(row), nil
}

func (r *ReturnPlanRepo) ListReturnPlanPointsByPlanID(ctx context.Context, planID uuid.UUID) ([]entity.ReturnPlanPoint, error) {
	rows, err := r.q.ListReturnPlanPointsByPlanID(ctx, planID)
	if err != nil {
		return []entity.ReturnPlanPoint{}, mapError(err)
	}

	return toReturnPlanPointEntities(rows), nil
}

func (r *ReturnPlanRepo) UpdateReturnPlanPointQuantities(ctx context.Context, pointID uuid.UUID, petCount int, canCount int, glassCount int) (entity.ReturnPlanPoint, error) {
	row, err := r.q.UpdateReturnPlanPointQuantities(ctx, database.UpdateReturnPlanPointQuantitiesParams{
		ID:       pointID,
		PetQty:   int32(petCount),
		CanQty:   int32(canCount),
		GlassQty: int32(glassCount),
	})
	if err != nil {
		return entity.ReturnPlanPoint{}, mapError(err)
	}

	return toReturnPlanPointEntity(row), nil
}

func (r *ReturnPlanRepo) ConfirmReturnPlanPoint(ctx context.Context, pointID uuid.UUID) (entity.ReturnPlanPoint, error) {
	row, err := r.q.ConfirmReturnPlanPoint(ctx, pointID)
	if err != nil {
		return entity.ReturnPlanPoint{}, mapError(err)
	}

	return toReturnPlanPointEntity(row), nil
}

func (r *ReturnPlanRepo) ConfirmAllPendingReturnPlanPoints(ctx context.Context, planID uuid.UUID) ([]entity.ReturnPlanPoint, error) {
	rows, err := r.q.ConfirmAllPendingReturnPlanPoints(ctx, planID)
	if err != nil {
		return []entity.ReturnPlanPoint{}, mapError(err)
	}

	return toReturnPlanPointEntities(rows), nil
}

func (r *ReturnPlanRepo) DeleteReturnPlanPoint(ctx context.Context, pointID uuid.UUID) error {
	err := r.q.DeleteReturnPlanPoint(ctx, pointID)
	if err != nil {
		return mapError(err)
	}

	return nil
}
