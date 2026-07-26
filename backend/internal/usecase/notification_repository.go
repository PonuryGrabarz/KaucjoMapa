package usecase

import (
	"context"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/google/uuid"
)

type INotificationsRepository interface {
	Create(ctx context.Context, notification entity.Notification) (entity.Notification, error)
	List(ctx context.Context, id uuid.UUID) ([]entity.Notification, error)
	Update(ctx context.Context, id uuid.UUID, notification entity.Notification) (entity.Notification, error)
	GetByID(ctx context.Context, id uuid.UUID) (entity.Notification, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
