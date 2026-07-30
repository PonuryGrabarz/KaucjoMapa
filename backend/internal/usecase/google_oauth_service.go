package usecase

import (
	"context"
	"time"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/port"
	"github.com/google/uuid"
)

type IGoogleOAuthService interface {
	GetUserByGoogleID(ctx context.Context, googleID string) (entity.User, error)
	CreateUser(ctx context.Context, user entity.User) (entity.User, error)
	IssueSession(ctx context.Context, user entity.User) (entity.RefreshToken, string, error)
}

type googleOAuthService struct {
	oAuthRepo IGoogleOAuthRepository
	authRepo  IRefreshTokenRepository
	tokenSvc  port.ITokenSigner
}

func NewGoogleOAuthService(oAuthRepo IGoogleOAuthRepository, authRepo IRefreshTokenRepository, tokenSvc port.ITokenSigner) IGoogleOAuthService {
	return &googleOAuthService{
		oAuthRepo: oAuthRepo,
		authRepo:  authRepo,
		tokenSvc:  tokenSvc,
	}
}

func (s *googleOAuthService) GetUserByGoogleID(ctx context.Context, googleID string) (entity.User, error) {
	user, err := s.oAuthRepo.GetUserByGoogleID(ctx, googleID)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (s *googleOAuthService) CreateUser(ctx context.Context, user entity.User) (entity.User, error) {
	createdUser, err := s.oAuthRepo.CreateUser(ctx, user)
	if err != nil {
		return entity.User{}, err
	}

	return createdUser, nil
}

func (s *googleOAuthService) IssueSession(ctx context.Context, user entity.User) (entity.RefreshToken, string, error) {
	token, err := s.tokenSvc.MakeJWT(user.ID, user.IsAdmin, user.Reputation, user.Username, user.Email, avatarClaim(user.AvatarPath))
	if err != nil {
		return entity.RefreshToken{}, "", err
	}

	refreshToken := entity.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     s.tokenSvc.MakeRefreshToken(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		RevokedAt: nil,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	}

	newToken, err := s.authRepo.Create(ctx, refreshToken)
	if err != nil {
		return entity.RefreshToken{}, "", err
	}

	return newToken, token, err
}
