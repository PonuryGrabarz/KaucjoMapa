package presenter

import (
	"time"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/google/uuid"
)

type CreateBookmarkRequest struct {
	ReturnPointID string `json:"return_point_id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type UpdateBookmarkRequest struct {
	ReturnPointID string `json:"return_point_id" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type BookmarkResponse struct {
	ID            uuid.UUID `json:"id" format:"uuid"`
	UserID        uuid.UUID `json:"user_id" format:"uuid"`
	ReturnPointID uuid.UUID `json:"return_point_id" format:"uuid"`
	CreatedAt     time.Time `json:"created_at"`
}

type BookmarkWithReturnPointResponse struct {
	ID          uuid.UUID                  `json:"id" format:"uuid"`
	UserID      uuid.UUID                  `json:"user_id" format:"uuid"`
	CreatedAt   time.Time                  `json:"created_at"`
	ReturnPoint ReturnPointDetailsResponse `json:"return_point"`
}

func BookmarkSuccessResponse(bookmark entity.Bookmark) BookmarkResponse {
	return BookmarkResponse{
		ID:            bookmark.ID,
		UserID:        bookmark.UserID,
		ReturnPointID: bookmark.ReturnPointID,
		CreatedAt:     bookmark.CreatedAt,
	}
}

func ListBookmarksSuccessResponse(bookmarks []entity.Bookmark) []BookmarkResponse {
	list := make([]BookmarkResponse, 0, len(bookmarks))
	for _, bookmark := range bookmarks {
		list = append(list, BookmarkSuccessResponse(bookmark))
	}
	return list
}

func BookmarkWithReturnPointSuccessResponse(item entity.BookmarkWithReturnPoint) BookmarkWithReturnPointResponse {
	return BookmarkWithReturnPointResponse{
		ID:          item.Bookmark.ID,
		UserID:      item.Bookmark.UserID,
		CreatedAt:   item.Bookmark.CreatedAt,
		ReturnPoint: CreateReturnPointSuccessResponse(item.ReturnPoint),
	}
}

func ListBookmarksWithReturnPointsSuccessResponse(items []entity.BookmarkWithReturnPoint) []BookmarkWithReturnPointResponse {
	list := make([]BookmarkWithReturnPointResponse, 0, len(items))
	for _, item := range items {
		list = append(list, BookmarkWithReturnPointSuccessResponse(item))
	}
	return list
}
