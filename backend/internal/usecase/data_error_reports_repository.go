package usecase

import (
	"context"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/google/uuid"
)

type IDataErrorReportsRepository interface {
	Create(ctx context.Context, report entity.ErrorReport) (entity.ErrorReport, error)
	GetByID(ctx context.Context, id uuid.UUID) (entity.UserErrorReport, error)
	List(ctx context.Context) ([]entity.UserErrorReport, error)
	Update(ctx context.Context, report entity.ErrorReport) (entity.ErrorReport, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
