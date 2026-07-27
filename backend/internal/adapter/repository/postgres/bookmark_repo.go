package postgres

import (
	"context"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/infrastructure/database"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase"
	"github.com/google/uuid"
)

type BookmarkRepo struct {
	q *database.Queries
}

func NewBookmarkRepo(q *database.Queries) *BookmarkRepo {
	return &BookmarkRepo{q: q}
}

var _ usecase.IBookmarkRepository = (*BookmarkRepo)(nil)

func (r *BookmarkRepo) Create(ctx context.Context, bookmark entity.Bookmark) (entity.Bookmark, error) {
	row, err := r.q.CreateBookmark(ctx, toCreateBookmarkParams(bookmark))
	if err != nil {
		return entity.Bookmark{}, mapError(err)
	}

	return toBookmarkEntity(row), nil
}

func (r *BookmarkRepo) GetByID(ctx context.Context, id uuid.UUID) (entity.Bookmark, error) {
	row, err := r.q.GetBookmarkById(ctx, id)
	if err != nil {
		return entity.Bookmark{}, mapError(err)
	}

	return toBookmarkEntity(row), nil
}

func (r *BookmarkRepo) GetByUserAndReturnPoint(ctx context.Context, userID, returnPointID uuid.UUID) (entity.Bookmark, error) {
	row, err := r.q.GetBookmarkByUserAndReturnPoint(ctx, database.GetBookmarkByUserAndReturnPointParams{
		UserID:        userID,
		ReturnPointID: returnPointID,
	})
	if err != nil {
		return entity.Bookmark{}, mapError(err)
	}

	return toBookmarkEntity(row), nil
}

func (r *BookmarkRepo) ListByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Bookmark, error) {
	rows, err := r.q.ListBookmarksByUserId(ctx, userID)
	if err != nil {
		return []entity.Bookmark{}, mapError(err)
	}

	bookmarks := make([]entity.Bookmark, 0, len(rows))
	for _, row := range rows {
		bookmarks = append(bookmarks, toBookmarkEntity(row))
	}

	return bookmarks, nil
}

func (r *BookmarkRepo) ListWithReturnPointsByUserID(ctx context.Context, userID uuid.UUID) ([]entity.BookmarkWithReturnPoint, error) {
	rows, err := r.q.ListBookmarksWithReturnPointsByUserId(ctx, userID)
	if err != nil {
		return []entity.BookmarkWithReturnPoint{}, mapError(err)
	}

	return toBookmarksWithReturnPoints(rows), nil
}

func (r *BookmarkRepo) Update(ctx context.Context, userID uuid.UUID, bookmark entity.Bookmark) (entity.Bookmark, error) {
	row, err := r.q.UpdateBookmark(ctx, database.UpdateBookmarkParams{
		ID:            bookmark.ID,
		ReturnPointID: bookmark.ReturnPointID,
		UserID:        userID,
	})
	if err != nil {
		return entity.Bookmark{}, mapError(err)
	}

	return toBookmarkEntity(row), nil
}

func (r *BookmarkRepo) DeleteByID(ctx context.Context, userID, id uuid.UUID) error {
	err := r.q.DeleteBookmarkById(ctx, database.DeleteBookmarkByIdParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		return mapError(err)
	}

	return nil
}

func (r *BookmarkRepo) DeleteByUserAndReturnPoint(ctx context.Context, userID, returnPointID uuid.UUID) error {
	err := r.q.DeleteBookmarkByUserAndReturnPoint(ctx, database.DeleteBookmarkByUserAndReturnPointParams{
		UserID:        userID,
		ReturnPointID: returnPointID,
	})
	if err != nil {
		return mapError(err)
	}

	return nil
}
