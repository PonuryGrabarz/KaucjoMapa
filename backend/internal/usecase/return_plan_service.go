package usecase

import (
	"context"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
)

type ReturnPlanService interface {
	CreateReturnPlan(ctx context.Context, plan entity.ReturnPlan) (entity.ReturnPlan, error)
	CreateReturnPlanPoint(ctx context.Context, planPoint entity.ReturnPlanPoint) (entity.ReturnPlanPoint, error)
}
