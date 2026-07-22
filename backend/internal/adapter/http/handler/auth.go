package handler

import (
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/adapter/http/middleware"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/adapter/http/presenter"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/dto"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/port"
	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	userSvc         usecase.IUserService
	hasherSvc       usecase.IPasswordHasher
	tokenSvc        port.ITokenSigner
	refreshTokenSvc usecase.IAuthService
}

func NewAuthHandler(userSvc usecase.IUserService, hasherSvc port.IPasswordHasher, tokenSvc port.ITokenSigner, refreshTokenSvc usecase.IAuthService) *AuthHandler {
	return &AuthHandler{
		userSvc:         userSvc,
		hasherSvc:       hasherSvc,
		tokenSvc:        tokenSvc,
		refreshTokenSvc: refreshTokenSvc,
	}
}

// Login godoc
// @Summary Zaloguj użytkownika
// @Description Weryfikuje dane logowania, zwraca access token i ustawia cookie RefreshToken.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body presenter.LoginRequest true "Dane logowania"
// @Success 200 {object} presenter.TokenResponse "Access token"
// @Header 200 {string} Set-Cookie "Cookie RefreshToken; Path=/api/auth; HttpOnly; Secure; SameSite=Strict"
// @Failure 400 {object} presenter.ErrorResponse "Nieprawidłowe body żądania"
// @Failure 401 {object} presenter.ErrorResponse "Nieprawidłowa nazwa użytkownika lub hasło"
// @Failure 500 {object} presenter.ErrorResponse "Wewnętrzny błąd serwera"
// @Router /login [post]
func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req presenter.LoginRequest

	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	refreshToken, token, err := h.refreshTokenSvc.Login(c.Context(), dto.LoginUserInput{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid username or password")
	}

	cookie := fiber.Cookie{
		Name:     "RefreshToken",
		Value:    refreshToken.Token,
		Path:     "/api/auth",
		Expires:  refreshToken.ExpiresAt,
		HTTPOnly: true,
		SameSite: "strict",
		Secure:   true,
	}

	c.Cookie(&cookie)
	return c.JSON(presenter.UserSuccessResponse(token))
}

// RefreshToken godoc
// @Summary Odśwież access token
// @Description Generuje nowy access token na podstawie Bearer JWT i cookie RefreshToken.
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Param Cookie header string true "Cookie w formacie RefreshToken=<token>"
// @Success 200 {object} presenter.TokenResponse "Nowy access token"
// @Failure 401 {object} presenter.ErrorResponse "Brak lub nieprawidłowy access albo refresh token"
// @Failure 500 {object} presenter.ErrorResponse "Wewnętrzny błąd serwera"
// @Router /auth/token [post]
func (h *AuthHandler) RefreshToken(c fiber.Ctx) error {
	refreshToken := string(c.Request().Header.Cookie("RefreshToken"))
	if refreshToken == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized: Refresh token not provided")
	}

	token, err := middleware.GetBearerToken(c.GetHeaders())
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid Authorization header")
	}

	newToken, err := h.refreshTokenSvc.RefreshToken(c.Context(), dto.RefreshTokenInput{
		RefreshToken: refreshToken,
		AccessToken:  token,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid refresh token or access token")
	}

	return c.JSON(presenter.UserSuccessResponse(newToken))
}

// RevokeToken godoc
// @Summary Unieważnij refresh token
// @Description Unieważnia refresh token przekazany w cookie.
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Param Cookie header string true "Cookie w formacie RefreshToken=<token>"
// @Success 204 "No Content"
// @Failure 401 {object} presenter.ErrorResponse "Brak lub nieprawidłowy access albo refresh token"
// @Failure 500 {object} presenter.ErrorResponse "Wewnętrzny błąd serwera"
// @Router /auth/revoke [post]
func (h *AuthHandler) RevokeToken(c fiber.Ctx) error {
	refreshToken := string(c.Request().Header.Cookie("RefreshToken"))

	if refreshToken == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized: Refresh token not provided")
	}

	token, err := middleware.GetBearerToken(c.GetHeaders())
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid Authorization header")
	}

	err = h.refreshTokenSvc.RevokeToken(c.Context(), dto.RevokeTokenInput{
		RefreshToken: refreshToken,
		AccessToken:  token,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid refresh token or access token")
	}

	c.ClearCookie("RefreshToken")

	return c.SendStatus(fiber.StatusNoContent)
}
