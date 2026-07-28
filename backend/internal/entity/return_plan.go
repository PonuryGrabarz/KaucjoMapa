package entity

import (
	"time"

	"github.com/google/uuid"
)

type PlanStatus string
type PlanPointStatus string

const (
	PlanStatusDraft              PlanStatus = "draft"
	PlanStatusSaved              PlanStatus = "saved"
	PlanStatusPartiallyConfirmed PlanStatus = "partially_confirmed"
	PlanStatusConfirmed          PlanStatus = "confirmed"
)

const (
	PlanPointStatusPending   PlanPointStatus = "pending"
	PlanPointStatusConfirmed PlanPointStatus = "confirmed"
)

type ReturnPlan struct {
	ID          uuid.UUID
	UserID      *uuid.UUID
	Status      PlanStatus
	TargetPet   int
	TargetCan   int
	TargetGlass int
	RatePet     int
	RateCan     int
	RateGlass   int
	ConfirmedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ReturnPlanPoint struct {
	ID            uuid.UUID
	ReturnPlanID  uuid.UUID
	ReturnPointID uuid.UUID
	PointName     string
	PointAddress  string
	PetCount      int
	CanCount      int
	GlassCount    int
	Status        PlanPointStatus
	ConfirmedAt   *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type ReturnPlanWithPoints struct {
	Plan   ReturnPlan
	Points []ReturnPlanPoint
}

const (
	DefaultRatePetGrosze   = 50
	DefaultRateCanGrosze   = 50
	DefaultRateGlassGrosze = 100
)
