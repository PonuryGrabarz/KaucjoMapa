package usecase

import (
	"context"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/google/uuid"
)

type ICommunityReportsRepository interface {
	Create(ctx context.Context, report entity.CommunityReport) (entity.CommunityReport, error)
	GetByID(ctx context.Context, id uuid.UUID) (entity.CommunityReport, error)
	List(ctx context.Context, id uuid.UUID) ([]entity.UserCommunityReport, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
