package presenter

import (
	"time"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/google/uuid"
)

type CreateReturnPlanRequest struct {
	TargetPet   int `json:"target_pet" validate:"omitempty,min=0" example:"11"`
	TargetCan   int `json:"target_can" validate:"omitempty,min=0" example:"15"`
	TargetGlass int `json:"target_glass" validate:"omitempty,min=0" example:"10"`
}

type UpdateReturnPlanTargetsRequest struct {
	TargetPet   int `json:"target_pet" validate:"omitempty,min=0" example:"11"`
	TargetCan   int `json:"target_can" validate:"omitempty,min=0" example:"15"`
	TargetGlass int `json:"target_glass" validate:"omitempty,min=0" example:"10"`
}

type AddReturnPlanPointRequest struct {
	ReturnPointID string `json:"return_point_id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	PetCount      int    `json:"pet_count" validate:"omitempty,min=0" example:"5"`
	CanCount      int    `json:"can_count" validate:"omitempty,min=0" example:"3"`
	GlassCount    int    `json:"glass_count" validate:"omitempty,min=0" example:"0"`
}

type UpdateReturnPlanPointQuantitiesRequest struct {
	PetCount   int `json:"pet_count" validate:"omitempty,min=0" example:"5"`
	CanCount   int `json:"can_count" validate:"omitempty,min=0" example:"3"`
	GlassCount int `json:"glass_count" validate:"omitempty,min=0" example:"2"`
}

type ReturnPlanPointResponse struct {
	ID            uuid.UUID  `json:"id" format:"uuid"`
	ReturnPlanID  uuid.UUID  `json:"return_plan_id" format:"uuid"`
	ReturnPointID uuid.UUID  `json:"return_point_id" format:"uuid"`
	PointName     string     `json:"point_name"`
	PointAddress  string     `json:"point_address"`
	PetCount      int        `json:"pet_count"`
	CanCount      int        `json:"can_count"`
	GlassCount    int        `json:"glass_count"`
	Status        string     `json:"status" example:"pending"`
	ValueGrosze   int        `json:"value_grosze"`
	ConfirmedAt   *time.Time `json:"confirmed_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type ReturnPlanResponse struct {
	ID               uuid.UUID                 `json:"id" format:"uuid"`
	UserID           uuid.UUID                 `json:"user_id" format:"uuid"`
	Status           string                    `json:"status" example:"draft"`
	TargetPet        int                       `json:"target_pet"`
	TargetCan        int                       `json:"target_can"`
	TargetGlass      int                       `json:"target_glass"`
	RatePetGrosze    int                       `json:"rate_pet_grosze"`
	RateCanGrosze    int                       `json:"rate_can_grosze"`
	RateGlassGrosze  int                       `json:"rate_glass_grosze"`
	AssignedPet      int                       `json:"assigned_pet"`
	AssignedCan      int                       `json:"assigned_can"`
	AssignedGlass    int                       `json:"assigned_glass"`
	TotalValueGrosze int                       `json:"total_value_grosze"`
	IsComplete       bool                      `json:"is_complete"`
	ConfirmedAt      *time.Time                `json:"confirmed_at,omitempty"`
	CreatedAt        time.Time                 `json:"created_at"`
	UpdatedAt        time.Time                 `json:"updated_at"`
	Points           []ReturnPlanPointResponse `json:"points"`
}

type ReturnPlanSummaryResponse struct {
	ID          uuid.UUID  `json:"id" format:"uuid"`
	UserID      uuid.UUID  `json:"user_id" format:"uuid"`
	Status      string     `json:"status" example:"saved"`
	TargetPet   int        `json:"target_pet"`
	TargetCan   int        `json:"target_can"`
	TargetGlass int        `json:"target_glass"`
	ConfirmedAt *time.Time `json:"confirmed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func ReturnPlanPointSuccessResponse(point entity.ReturnPlanPoint, ratePet, rateCan, rateGlass int) ReturnPlanPointResponse {
	return ReturnPlanPointResponse{
		ID:            point.ID,
		ReturnPlanID:  point.ReturnPlanID,
		ReturnPointID: point.ReturnPointID,
		PointName:     point.PointName,
		PointAddress:  point.PointAddress,
		PetCount:      point.PetCount,
		CanCount:      point.CanCount,
		GlassCount:    point.GlassCount,
		Status:        string(point.Status),
		ValueGrosze:   point.PetCount*ratePet + point.CanCount*rateCan + point.GlassCount*rateGlass,
		ConfirmedAt:   point.ConfirmedAt,
		CreatedAt:     point.CreatedAt,
		UpdatedAt:     point.UpdatedAt,
	}
}

func ReturnPlanSuccessResponse(item entity.ReturnPlanWithPoints) ReturnPlanResponse {
	plan := item.Plan
	userID := uuid.Nil
	if plan.UserID != nil {
		userID = *plan.UserID
	}

	points := make([]ReturnPlanPointResponse, 0, len(item.Points))
	assignedPet, assignedCan, assignedGlass := 0, 0, 0
	totalValue := 0

	for _, point := range item.Points {
		resp := ReturnPlanPointSuccessResponse(point, plan.RatePet, plan.RateCan, plan.RateGlass)
		points = append(points, resp)
		assignedPet += point.PetCount
		assignedCan += point.CanCount
		assignedGlass += point.GlassCount
		totalValue += resp.ValueGrosze
	}

	return ReturnPlanResponse{
		ID:               plan.ID,
		UserID:           userID,
		Status:           string(plan.Status),
		TargetPet:        plan.TargetPet,
		TargetCan:        plan.TargetCan,
		TargetGlass:      plan.TargetGlass,
		RatePetGrosze:    plan.RatePet,
		RateCanGrosze:    plan.RateCan,
		RateGlassGrosze:  plan.RateGlass,
		AssignedPet:      assignedPet,
		AssignedCan:      assignedCan,
		AssignedGlass:    assignedGlass,
		TotalValueGrosze: totalValue,
		IsComplete:       assignedPet == plan.TargetPet && assignedCan == plan.TargetCan && assignedGlass == plan.TargetGlass,
		ConfirmedAt:      plan.ConfirmedAt,
		CreatedAt:        plan.CreatedAt,
		UpdatedAt:        plan.UpdatedAt,
		Points:           points,
	}
}

func ReturnPlanSummarySuccessResponse(plan entity.ReturnPlan) ReturnPlanSummaryResponse {
	userID := uuid.Nil
	if plan.UserID != nil {
		userID = *plan.UserID
	}

	return ReturnPlanSummaryResponse{
		ID:          plan.ID,
		UserID:      userID,
		Status:      string(plan.Status),
		TargetPet:   plan.TargetPet,
		TargetCan:   plan.TargetCan,
		TargetGlass: plan.TargetGlass,
		ConfirmedAt: plan.ConfirmedAt,
		CreatedAt:   plan.CreatedAt,
		UpdatedAt:   plan.UpdatedAt,
	}
}

func ListReturnPlanSummariesSuccessResponse(plans []entity.ReturnPlan) []ReturnPlanSummaryResponse {
	list := make([]ReturnPlanSummaryResponse, 0, len(plans))
	for _, plan := range plans {
		list = append(list, ReturnPlanSummarySuccessResponse(plan))
	}
	return list
}
