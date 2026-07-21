package handler

import (
	"time"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/adapter/http/presenter"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/auth"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/user"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type AuthHandler struct {
	userSvc         user.IUserService
	hasherSvc       auth.IPasswordHasher
	tokenSvc        auth.ITokenSigner
	refreshTokenSvc auth.IRefreshTokenService
}

func NewAuthHandler(userSvc user.IUserService, hasherSvc auth.IPasswordHasher, tokenSvc auth.ITokenSigner, refreshTokenSvc auth.IRefreshTokenService) *AuthHandler {
	return &AuthHandler{
		userSvc:         userSvc,
		hasherSvc:       hasherSvc,
		tokenSvc:        tokenSvc,
		refreshTokenSvc: refreshTokenSvc,
	}
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req presenter.LoginRequest

	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	user, err := h.userSvc.GetUserByUsername(c.Context(), req.Username)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
	}

	isPasswordValid, err := h.hasherSvc.CheckPasswordHash(req.Password, user.PasswordHash)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Server error")
	}

	if !isPasswordValid {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
	}

	token, err := h.tokenSvc.MakeJWT(user.ID, user.IsAdmin, user.Reputation, user.Username, user.Email)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate JWT token")
	}

	refreshTokenStr := h.tokenSvc.MakeRefreshToken()

	refreshToken, err := h.refreshTokenSvc.Create(c.Context(), entity.RefreshToken{
		ID:        uuid.New(),
		Token:     refreshTokenStr,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create refresh token")
	}

	cookie := new(fiber.Cookie)
	cookie.Name = "RefreshToken"
	cookie.Value = refreshToken.Token
	cookie.Path = "/auth"
	cookie.Expires = refreshToken.ExpiresAt
	cookie.HTTPOnly = true
	cookie.SameSite = "strict"
	cookie.Secure = true

	c.Cookie(cookie)
	return c.JSON(presenter.UserSuccessResponse(token))
}

func (h *AuthHandler) RefreshToken(c fiber.Ctx) error {
	refreshToken := string(c.Request().Header.Cookie("RefreshToken"))

	if refreshToken == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized: Refresh token not provided")
	}

	claims, ok := c.Locals("claims").(*auth.JWTClaims)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid JWT claims")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid user ID in JWT claims")
	}

	userToken, err := h.refreshTokenSvc.GetToken(c.Context(), userID, refreshToken)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid refresh token")
	}

	if time.Now().Compare(userToken.ExpiresAt) != -1 {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized: Refresh token is revoked or expired")
	}

	user, err := h.userSvc.GetUserByID(c.Context(), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "User not found")
	}

	newToken, err := h.tokenSvc.MakeJWT(user.ID, user.IsAdmin, user.Reputation, user.Username, user.Email)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate new JWT token")
	}

	return c.JSON(presenter.UserSuccessResponse(newToken))
}

func (h *AuthHandler) RevokeToken(c fiber.Ctx) error {
	refreshToken := string(c.Request().Header.Cookie("RefreshToken"))

	if refreshToken == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized: Refresh token not provided")
	}

	claims, ok := c.Locals("claims").(*auth.JWTClaims)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid JWT claims")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid user ID in JWT claims")
	}

	userToken, err := h.refreshTokenSvc.GetToken(c.Context(), userID, refreshToken)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid refresh token")
	}

	err = h.refreshTokenSvc.RevokeToken(c.Context(), userToken.ID, userToken.Token)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to revoke refresh token")
	}

	return c.SendStatus(fiber.StatusNoContent)
}
