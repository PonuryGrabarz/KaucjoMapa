package usecase

import (
	"context"
	"time"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/dto"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/port"
	"github.com/google/uuid"
)

type IAuthService interface {
	Login(ctx context.Context, in dto.LoginUserInput) (entity.RefreshToken, string, error)
	RefreshToken(ctx context.Context, in dto.RefreshTokenInput) (string, error)
	RevokeToken(ctx context.Context, in dto.RevokeTokenInput) error
}

type authService struct {
	refreshTokenRepo IRefreshTokenRepository
	userSvc          IUserService
	hasherSvc        IPasswordHasher
	tokenSvc         port.ITokenSigner
}

func NewAuthService(refreshTokenRepo IRefreshTokenRepository, userSvc IUserService, hasherSvc IPasswordHasher, tokenSvc port.ITokenSigner) IAuthService {
	return &authService{
		refreshTokenRepo: refreshTokenRepo,
		userSvc:          userSvc,
		hasherSvc:        hasherSvc,
		tokenSvc:         tokenSvc,
	}
}

func (s *authService) Login(ctx context.Context, in dto.LoginUserInput) (entity.RefreshToken, string, error) {
	user, err := s.userSvc.GetUserByUsername(ctx, in.Username)
	if err != nil {
		return entity.RefreshToken{}, "", err
	}

	isPasswordValid, err := s.hasherSvc.CheckPasswordHash(in.Password, user.PasswordHash)
	if err != nil {
		return entity.RefreshToken{}, "", entity.NewInvalidCredentials()
	}

	if !isPasswordValid {
		return entity.RefreshToken{}, "", entity.NewInvalidCredentials()
	}

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

	newToken, err := s.refreshTokenRepo.Create(ctx, refreshToken)
	if err != nil {
		return entity.RefreshToken{}, "", err
	}

	return newToken, token, err
}

func (s *authService) RefreshToken(ctx context.Context, in dto.RefreshTokenInput) (string, error) {
	claims, err := s.tokenSvc.ExtractClaimsFromJWT(in.AccessToken)
	if err != nil {
		return "", err
	}

	userToken, err := s.refreshTokenRepo.GetToken(ctx, claims.UserID, in.RefreshToken)
	if err != nil {
		return "", err
	}

	if time.Now().Compare(userToken.ExpiresAt) != -1 || userToken.RevokedAt != nil {
		return "", entity.NewUnauthorized("refresh token expired or revoked")
	}

	user, err := s.userSvc.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return "", err
	}

	newToken, err := s.tokenSvc.MakeJWT(user.ID, user.IsAdmin, user.Reputation, user.Username, user.Email, avatarClaim(user.AvatarPath))
	if err != nil {
		return "", err
	}

	return newToken, nil
}

func (s *authService) RevokeToken(ctx context.Context, in dto.RevokeTokenInput) error {
	claims, err := s.tokenSvc.ValidateJWT(in.AccessToken)
	if err != nil {
		return err
	}

	refreshToken, err := s.refreshTokenRepo.GetToken(ctx, claims.UserID, in.RefreshToken)
	if err != nil {
		return err
	}

	err = s.refreshTokenRepo.RevokeToken(ctx, refreshToken.UserID, refreshToken.Token)
	if err != nil {
		return err
	}

	return nil
}
