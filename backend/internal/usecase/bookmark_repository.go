package usecase

import (
	"context"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/google/uuid"
)

type IBookmarkRepository interface {
	Create(ctx context.Context, bookmark entity.Bookmark) (entity.Bookmark, error)
	GetByID(ctx context.Context, id uuid.UUID) (entity.Bookmark, error)
	GetByUserAndReturnPoint(ctx context.Context, userID, returnPointID uuid.UUID) (entity.Bookmark, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Bookmark, error)
	ListWithReturnPointsByUserID(ctx context.Context, userID uuid.UUID) ([]entity.BookmarkWithReturnPoint, error)
	Update(ctx context.Context, userID uuid.UUID, bookmark entity.Bookmark) (entity.Bookmark, error)
	DeleteByID(ctx context.Context, userID, id uuid.UUID) error
	DeleteByUserAndReturnPoint(ctx context.Context, userID, returnPointID uuid.UUID) error
}
