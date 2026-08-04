package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/dto"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/port"
	"github.com/google/uuid"
)

const (
	maxAvatarBytes     = 2 << 20 // 2 MiB
	avatarAPIClaimPath = "/api/users/avatar"
)

var allowedAvatarContentTypes = map[string]struct{}{
	"image/jpeg": {},
	"image/png":  {},
	"image/webp": {},
}

type IUserService interface {
	Create(ctx context.Context, in dto.CreateUserInput) (entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (entity.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (entity.User, error)
	ChangePassword(ctx context.Context, in dto.ChangePasswordInput) error
	ChangeAvatar(ctx context.Context, id uuid.UUID, contentType string, data []byte) error
	GetAvatar(ctx context.Context, id uuid.UUID) (data []byte, contentType string, err error)
}

type userService struct {
	userRepo  IUserRepository
	hasherSvc port.IPasswordHasher
	tokenSvc  port.ITokenSigner
	avatars   port.IAvatarStorage
}

func NewUserService(userRepo IUserRepository, hasher port.IPasswordHasher, signer port.ITokenSigner, avatars port.IAvatarStorage) IUserService {
	return &userService{
		userRepo:  userRepo,
		hasherSvc: hasher,
		tokenSvc:  signer,
		avatars:   avatars,
	}
}

func avatarClaim(path string) string {
	if path == "" {
		return ""
	}
	return avatarAPIClaimPath
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

	accessToken, err := s.tokenSvc.MakeJWT(user.ID, user.IsAdmin, user.Reputation, user.Username, user.Email, avatarClaim(user.AvatarPath))
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

func (s *userService) ChangePassword(ctx context.Context, in dto.ChangePasswordInput) error {
	user, err := s.userRepo.GetUserByID(ctx, in.UserID)
	if err != nil {
		return entity.NewUnauthorized("user not found")
	}

	if in.NewPassword != in.RepeatNewPassword {
		return entity.NewBadRequest("new password and repeat new password do not match")
	}

	isPasswordValid, err := s.hasherSvc.CheckPasswordHash(in.OldPassword, user.PasswordHash)
	if err != nil {
		return entity.NewUnauthorized("invalid old password")
	}

	if !isPasswordValid {
		return entity.NewUnauthorized("invalid old password")
	}

	hashedPassword, err := s.hasherSvc.HashPassword(in.NewPassword)
	if err != nil {
		return entity.NewInternal(err)
	}

	err = s.userRepo.UpdatePassword(ctx, in.UserID, hashedPassword)
	if err != nil {
		return entity.NewInternal(err)
	}

	return nil
}

func (s *userService) ChangeAvatar(ctx context.Context, id uuid.UUID, contentType string, data []byte) error {
	if len(data) == 0 {
		return entity.NewBadRequest("avatar body is required")
	}
	if len(data) > maxAvatarBytes {
		return entity.NewBadRequest("avatar exceeds maximum size of 2MiB")
	}

	contentType = strings.TrimSpace(strings.Split(contentType, ";")[0])
	if _, ok := allowedAvatarContentTypes[contentType]; !ok {
		return entity.NewBadRequest("unsupported content type; allowed: image/jpeg, image/png, image/webp")
	}

	user, err := s.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	path, err := s.avatars.Save(ctx, id, contentType, data)
	if err != nil {
		return err
	}

	if user.AvatarPath != "" && user.AvatarPath != path {
		_ = s.avatars.Delete(ctx, user.AvatarPath)
	}

	if err := s.userRepo.ChangeAvatar(ctx, id, path); err != nil {
		_ = s.avatars.Delete(ctx, path)
		return err
	}

	return nil
}

func (s *userService) GetAvatar(ctx context.Context, id uuid.UUID) ([]byte, string, error) {
	user, err := s.GetUserByID(ctx, id)
	if err != nil {
		return nil, "", err
	}

	if user.AvatarPath == "" {
		return nil, "", entity.NewNotFound("avatar not found")
	}

	return s.avatars.Load(ctx, user.AvatarPath)
}
