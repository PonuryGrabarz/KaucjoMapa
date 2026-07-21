package seed

import (
	"context"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
)

type DataSource interface {
	FetchReturnPoints(ctx context.Context) ([]entity.ReturnPoint, error)
}
