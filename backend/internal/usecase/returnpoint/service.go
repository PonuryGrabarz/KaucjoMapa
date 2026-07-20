package returnpoint

import (
	"context"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/google/uuid"
)

type IReturnPointService interface {
	List(ctx context.Context, f Filter) ([]entity.ReturnPoint, error)
	GetByID(ctx context.Context, id uuid.UUID) (entity.ReturnPoint, error)
	ExistsNearby(ctx context.Context, lat, lng float64) (bool, error)
	Create(ctx context.Context, userID uuid.UUID, p entity.ReturnPoint) (entity.ReturnPoint, error)
	Delete(ctx context.Context, actor entity.User, id uuid.UUID) error
	Update(ctx context.Context, actor entity.User, p entity.ReturnPoint) (entity.ReturnPoint, error)
}

type returnPointService struct {
	repo IReturnPointRepository
}

func NewService(repo IReturnPointRepository) IReturnPointService {
	return &returnPointService{repo: repo}
}

func (s *returnPointService) List(ctx context.Context, f Filter) ([]entity.ReturnPoint, error) {
	retPoints, err := s.repo.List(ctx, f)
	if err != nil {
		return []entity.ReturnPoint{}, err
	}

	return retPoints, nil
}

func (s *returnPointService) GetByID(ctx context.Context, id uuid.UUID) (entity.ReturnPoint, error) {
	retPoint, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return entity.ReturnPoint{}, err
	}

	return retPoint, nil
}

func (s *returnPointService) ExistsNearby(ctx context.Context, lat, lng float64) (bool, error) {
	nearby, err := s.repo.ExistsNearby(ctx, lat, lng)
	if err != nil {
		return false, err
	}

	return nearby, nil
}

func (s *returnPointService) Create(ctx context.Context, userID uuid.UUID, p entity.ReturnPoint) (entity.ReturnPoint, error) {
	nearby, err := s.repo.ExistsNearby(ctx, p.Latitude, p.Longitude)
	if err != nil {
		return entity.ReturnPoint{}, err
	}
	if nearby {
		return entity.ReturnPoint{}, entity.ErrPointExistsNearby
	}

	p.ID = uuid.New()
	p.UserID = &userID
	p.IsVerified = true

	return s.repo.Create(ctx, p)
}

func (s *returnPointService) Delete(ctx context.Context, actor entity.User, id uuid.UUID) error {
	if !actor.IsAdmin {
		return entity.ErrForbidden
	}

	return s.repo.Delete(ctx, id)
}

func (s *returnPointService) Update(ctx context.Context, actor entity.User, p entity.ReturnPoint) (entity.ReturnPoint, error) {
	_, err := s.repo.GetByID(ctx, p.ID)
	if err != nil {
		return entity.ReturnPoint{}, err
	}

	p.IsVerified = false
	return s.repo.Update(ctx, p)
}
