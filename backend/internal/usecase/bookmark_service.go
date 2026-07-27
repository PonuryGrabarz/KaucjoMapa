package usecase

import (
	"context"
	"time"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/google/uuid"
)

type IBookmarkService interface {
	Create(ctx context.Context, userID, returnPointID uuid.UUID) (entity.Bookmark, error)
	GetByID(ctx context.Context, id uuid.UUID) (entity.Bookmark, error)
	GetByUserAndReturnPoint(ctx context.Context, userID, returnPointID uuid.UUID) (entity.Bookmark, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Bookmark, error)
	ListWithReturnPointsByUserID(ctx context.Context, userID uuid.UUID) ([]entity.BookmarkWithReturnPoint, error)
	Update(ctx context.Context, userID, id, returnPointID uuid.UUID) (entity.Bookmark, error)
	DeleteByID(ctx context.Context, userID, id uuid.UUID) error
	DeleteByUserAndReturnPoint(ctx context.Context, userID, returnPointID uuid.UUID) error
}

type bookmarkService struct {
	repo IBookmarkRepository
}

func NewBookmarkService(repo IBookmarkRepository) IBookmarkService {
	return &bookmarkService{repo: repo}
}

func (s *bookmarkService) Create(ctx context.Context, userID, returnPointID uuid.UUID) (entity.Bookmark, error) {
	if userID == uuid.Nil {
		return entity.Bookmark{}, entity.NewBadRequest("user id is required")
	}
	if returnPointID == uuid.Nil {
		return entity.Bookmark{}, entity.NewBadRequest("return point id is required")
	}

	bookmark, err := s.repo.Create(ctx, entity.Bookmark{
		ID:            uuid.New(),
		UserID:        userID,
		ReturnPointID: returnPointID,
		CreatedAt:     time.Now().UTC(),
	})
	if err != nil {
		return entity.Bookmark{}, err
	}

	return bookmark, nil
}

func (s *bookmarkService) GetByID(ctx context.Context, id uuid.UUID) (entity.Bookmark, error) {
	bookmark, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return entity.Bookmark{}, err
	}

	return bookmark, nil
}

func (s *bookmarkService) GetByUserAndReturnPoint(ctx context.Context, userID, returnPointID uuid.UUID) (entity.Bookmark, error) {
	bookmark, err := s.repo.GetByUserAndReturnPoint(ctx, userID, returnPointID)
	if err != nil {
		return entity.Bookmark{}, err
	}

	return bookmark, nil
}

func (s *bookmarkService) ListByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Bookmark, error) {
	bookmarks, err := s.repo.ListByUserID(ctx, userID)
	if err != nil {
		return []entity.Bookmark{}, err
	}

	return bookmarks, nil
}

func (s *bookmarkService) ListWithReturnPointsByUserID(ctx context.Context, userID uuid.UUID) ([]entity.BookmarkWithReturnPoint, error) {
	bookmarks, err := s.repo.ListWithReturnPointsByUserID(ctx, userID)
	if err != nil {
		return []entity.BookmarkWithReturnPoint{}, err
	}

	return bookmarks, nil
}

func (s *bookmarkService) Update(ctx context.Context, userID, id, returnPointID uuid.UUID) (entity.Bookmark, error) {
	if returnPointID == uuid.Nil {
		return entity.Bookmark{}, entity.NewBadRequest("return point id is required")
	}

	bookmark, err := s.repo.Update(ctx, userID, entity.Bookmark{
		ID:            id,
		UserID:        userID,
		ReturnPointID: returnPointID,
	})
	if err != nil {
		return entity.Bookmark{}, err
	}

	return bookmark, nil
}

func (s *bookmarkService) DeleteByID(ctx context.Context, userID, id uuid.UUID) error {
	return s.repo.DeleteByID(ctx, userID, id)
}

func (s *bookmarkService) DeleteByUserAndReturnPoint(ctx context.Context, userID, returnPointID uuid.UUID) error {
	return s.repo.DeleteByUserAndReturnPoint(ctx, userID, returnPointID)
}
