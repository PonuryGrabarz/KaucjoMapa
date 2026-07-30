package port

import (
	"context"

	"github.com/google/uuid"
)

type IAvatarStorage interface {
	Save(ctx context.Context, userID uuid.UUID, contentType string, data []byte) (path string, err error)
	Load(ctx context.Context, path string) (data []byte, contentType string, err error)
	Delete(ctx context.Context, path string) error
}
