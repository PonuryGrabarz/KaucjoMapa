package entity

import "github.com/google/uuid"

type Notification struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        string
	Description string
	Seen        bool
	CreatedAt   string
}
