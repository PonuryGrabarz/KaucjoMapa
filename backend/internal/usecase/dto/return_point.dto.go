package dto

import "github.com/google/uuid"

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
	OpeningHours  string
	Description   string
	CurrentStatus string
}
