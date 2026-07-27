package usecase

import (
	"context"
	"time"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/dto"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/port"
	"github.com/google/uuid"
)

type IUserService interface {
	Create(ctx context.Context, in dto.CreateUserInput) (entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (entity.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (entity.User, error)
}

type userService struct {
	userRepo  IUserRepository
	hasherSvc port.IPasswordHasher
	tokenSvc  port.ITokenSigner
}

func NewUserService(userRepo IUserRepository, hasher port.IPasswordHasher, signer port.ITokenSigner) IUserService {
	return &userService{
		userRepo:  userRepo,
		hasherSvc: hasher,
		tokenSvc:  signer,
	}
}

func (s *userService) Create(ctx context.Context, in dto.CreateUserInput) (entity.User, error) {
	userByUsername, err := s.GetUserByUsername(ctx, in.Username)
	if err == nil && in.Username == userByUsername.Username {
		return entity.User{}, entity.NewConflict("username already taken")
	}

	passwordHash, err := s.hasherSvc.HashPassword(in.Password)
	if err != nil {
		return entity.User{}, entity.NewInternal(err)
	}

	user := entity.User{
		ID:           uuid.New(),
		Email:        in.Email,
		Username:     in.Username,
		PasswordHash: passwordHash,
		Reputation:   0,
		IsAdmin:      false,
		IsBanned:     false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	accessToken, err := s.tokenSvc.MakeJWT(user.ID, user.IsAdmin, user.Reputation, user.Username, user.Email)
	if err != nil {
		return entity.User{}, entity.NewInternal(err)
	}
	user.Token = accessToken

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
