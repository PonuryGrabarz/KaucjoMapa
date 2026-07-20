package user

import (
	"context"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/google/uuid"
)

type IUserService interface {
	Create(ctx context.Context, user entity.User) (entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (entity.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (entity.User, error)
}

type userService struct {
	userRepo IUserRepository
}

func NewService(userRepo IUserRepository) IUserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) Create(ctx context.Context, user entity.User) (entity.User, error) {
	user.ID = uuid.New()

	newUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return entity.User{}, err
	}

	return newUser, nil
}

func (s *userService) GetUserByUsername(ctx context.Context, username string) (entity.User, error) {
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (s *userService) GetUserByID(ctx context.Context, id uuid.UUID) (entity.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}
