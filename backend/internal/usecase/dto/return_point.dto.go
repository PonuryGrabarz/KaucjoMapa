package dto

import (
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/google/uuid"
)

type CreateReturnPointInput struct {
	Name          string
	UserID        uuid.UUID
	Type          string
	Address       string
	Latitude      float64
	Longitude     float64
	Pet           bool
	Can           bool
	Glass         bool
	ShortName     string
	BgColor       string
	TextColor     string
	RefundMethod  []entity.RefundMethod
	OpeningHours  entity.OpeningHours
	Description   string
	CurrentStatus string
}
