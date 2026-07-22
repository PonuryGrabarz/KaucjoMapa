package usecase

import (
	"context"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/google/uuid"
)

type Filter struct {
	Query  string
	Type   entity.PointType
	Status entity.PointStatus
	Offset int32
	Limit  int32
}

type IReturnPointRepository interface {
	List(ctx context.Context, f Filter) ([]entity.ReturnPoint, error)
	ListSize(ctx context.Context, f Filter) (int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (entity.ReturnPoint, error)
	ExistsNearby(ctx context.Context, lat, lng float64) (bool, error)
	Create(ctx context.Context, p entity.ReturnPoint) (entity.ReturnPoint, error)
	Update(ctx context.Context, p entity.ReturnPoint) (entity.ReturnPoint, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
