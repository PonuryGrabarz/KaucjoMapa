package usecase

import (
	"context"
	"time"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/dto"
	"github.com/google/uuid"
)

type IReturnPointService interface {
	List(ctx context.Context, f Filter) ([]entity.ReturnPoint, error)
	ListSize(ctx context.Context, f Filter) (int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (entity.ReturnPoint, error)
	ExistsNearby(ctx context.Context, lat, lng float64) (bool, error)
	Create(ctx context.Context, in dto.CreateReturnPointInput) (entity.ReturnPoint, error)
	Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
	Update(ctx context.Context, actor entity.User, p entity.ReturnPoint) (entity.ReturnPoint, error)
}

type returnPointService struct {
	repo    IReturnPointRepository
	userSvc IUserService
}

func NewReturnPointService(repo IReturnPointRepository, userSvc IUserService) IReturnPointService {
	return &returnPointService{repo: repo, userSvc: userSvc}
}

func (s *returnPointService) List(ctx context.Context, f Filter) ([]entity.ReturnPoint, error) {
	retPoints, err := s.repo.List(ctx, f)
	if err != nil {
		return []entity.ReturnPoint{}, err
	}

	return retPoints, nil
}

func (s *returnPointService) ListSize(ctx context.Context, f Filter) (int64, error) {
	size, err := s.repo.ListSize(ctx, f)
	if err != nil {
		return 0, err
	}

	return size, nil
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

func (s *returnPointService) Create(ctx context.Context, in dto.CreateReturnPointInput) (entity.ReturnPoint, error) {
	nearby, err := s.repo.ExistsNearby(ctx, in.Latitude, in.Longitude)
	if err != nil {
		return entity.ReturnPoint{}, err
	}
	if nearby {
		return entity.ReturnPoint{}, entity.ErrPointExistsNearby
	}

	returnPoint := entity.ReturnPoint{
		ID:            uuid.New(),
		UserID:        &in.UserID,
		Name:          in.Name,
		Type:          entity.PointType(in.Type),
		Address:       in.Address,
		Latitude:      in.Latitude,
		Longitude:     in.Longitude,
		OpeningHours:  in.OpeningHours,
		Description:   in.Description,
		Pet:           in.Pet,
		Can:           in.Can,
		Glass:         in.Glass,
		CurrentStatus: entity.PointStatus(in.CurrentStatus),
		IsVerified:    false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	newReturnPoint, err := s.repo.Create(ctx, returnPoint)
	if err != nil {
		return entity.ReturnPoint{}, err
	}

	return newReturnPoint, nil
}

func (s *returnPointService) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	user, err := s.userSvc.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if !user.IsAdmin {
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
