package entity

import (
	"time"

	"github.com/google/uuid"
)

type Bookmark struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	ReturnPointID uuid.UUID
	CreatedAt     time.Time
}

// BookmarkWithReturnPoint łączy zakładkę z danymi punktu zwrotu.
type BookmarkWithReturnPoint struct {
	Bookmark    Bookmark
	ReturnPoint ReturnPoint
}
