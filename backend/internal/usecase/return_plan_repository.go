package usecase

import (
	"context"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/google/uuid"
)

type IReturnPlanRepository interface {
	CreateReturnPlan(ctx context.Context, plan entity.ReturnPlan) (entity.ReturnPlan, error)
	GetReturnPlanByID(ctx context.Context, id uuid.UUID) (entity.ReturnPlan, error)
	GetReturnPlanByIDAndUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (entity.ReturnPlan, error)
	GetDraftReturnPlanByUserID(ctx context.Context, userID uuid.UUID) (entity.ReturnPlan, error)
	ListReturnPlansByUserID(ctx context.Context, userID uuid.UUID) ([]entity.ReturnPlan, error)
	UpdateReturnPlanTargets(ctx context.Context, planID uuid.UUID, targetPet int, targetCan int, targetGlass int) (entity.ReturnPlan, error)
	UpdateReturnPlanStatus(ctx context.Context, planID uuid.UUID, status entity.PlanStatus) (entity.ReturnPlan, error)
	DeleteReturnPlanByIDAndUserID(ctx context.Context, planID uuid.UUID, userID uuid.UUID) error
	ConfirmReturnPlan(ctx context.Context, planID uuid.UUID) (entity.ReturnPlan, error)
	CreateReturnPlanPoint(ctx context.Context, point entity.ReturnPlanPoint) (entity.ReturnPlanPoint, error)
	GetReturnPlanPointByID(ctx context.Context, id uuid.UUID) (entity.ReturnPlanPoint, error)
	GetReturnPlanPointByPlanAndPoint(ctx context.Context, planID uuid.UUID, pointID uuid.UUID) (entity.ReturnPlanPoint, error)
	ListReturnPlanPointsByPlanID(ctx context.Context, planID uuid.UUID) ([]entity.ReturnPlanPoint, error)
	UpdateReturnPlanPointQuantities(ctx context.Context, pointID uuid.UUID, petCount int, canCount int, glassCount int) (entity.ReturnPlanPoint, error)
	ConfirmReturnPlanPoint(ctx context.Context, pointID uuid.UUID) (entity.ReturnPlanPoint, error)
	ConfirmAllPendingReturnPlanPoints(ctx context.Context, planID uuid.UUID) ([]entity.ReturnPlanPoint, error)
	DeleteReturnPlanPoint(ctx context.Context, pointID uuid.UUID) error
}
