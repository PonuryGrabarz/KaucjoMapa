package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"time"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/auth"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/database"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// LoginUserHandler godoc
//
//	@Summary		Logowanie użytkownika
//	@Description	Weryfikuje dane logowania, zwraca dane użytkownika z access tokenem i zapisuje refresh token w cookie.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		server.LoginUserHandler.parameters	true	"Dane logowania"
//	@Success		200		{object}	User
//	@Header			200		{string}	Set-Cookie	"RefreshToken=<token>; Path=/auth; HttpOnly; SameSite=Strict"
//	@Failure		401		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Failure		400		{object}	map[string]string
//	@Router			/login [post]
func (s *FiberServer) LoginUserHandler(c fiber.Ctx) error {
	type parameters struct {
		Password string `json:"password"`
		Username string `json:"username"`
	}

	params := parameters{}

	err := json.Unmarshal(c.Request().Body(), &params)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	user, err := s.db.Queries.GetUserByUsername(context.Background(), params.Username)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "User not found")
	}

	isPasswordValid, err := auth.CheckPasswordHash(params.Password, user.PasswordHash)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Invalid credentials")
	}

	if !isPasswordValid {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
	}

	return s.issueSession(c, user)
}

func (s *FiberServer) issueSession(c fiber.Ctx, user database.User) error {
	accessToken, err := auth.MakeJWT(user.ID, user.Username, user.Email, JWTSecret, time.Hour)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate JWT token")
	}

	now := time.Now()
	dbRefreshToken, err := s.db.Queries.CreateRefreshToken(c.Context(), database.CreateRefreshTokenParams{
		ID:        uuid.New(),
		Token:     auth.MakeRefreshToken(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		ExpiresAt: now.Add(60 * 24 * time.Hour),
		RevokedAt: sql.NullTime{},
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create refresh token")
	}

	c.Cookie(&fiber.Cookie{
		Name:     "RefreshToken",
		Value:    dbRefreshToken.Token,
		Path:     "/auth",
		Expires:  dbRefreshToken.ExpiresAt,
		HTTPOnly: true,
		Secure:   os.Getenv("APP_ENV") == "production",
		SameSite: "Lax",
	})

	return c.JSON(User{
		ID:         user.ID,
		Email:      user.Email,
		Username:   user.Username,
		Token:      accessToken,
		Reputation: user.Reputation,
		IsAdmin:    user.IsAdmin,
		IsBanned:   user.IsBanned,
		CreatedAt:  user.CreatedAt.String(),
		UpdatedAt:  user.UpdatedAt.String(),
	})
}

// refreshTokenHandler godoc
//
//	@Summary		Odświeżenie tokenu JWT
//	@Description	Weryfikuje access token z nagłówka Authorization i refresh token z cookie, a następnie zwraca nowy access token.
//	@Tags			auth
//	@Produce		json
//	@Param			Authorization	header		string				true	"Bearer access token"
//	@Param			Cookie			header		string				true	"RefreshToken=<refresh_token>"
//	@Success		200				{object}	map[string]string	"Nowy access token"
//	@Failure		401				{object}	map[string]string	"Brak, nieprawidłowy lub wygasły token"
//	@Failure		500				{object}	map[string]string	"Błąd generowania access tokenu"
//	@Router			/auth/token [post]
func (s *FiberServer) refreshTokenHandler(c fiber.Ctx) error {
	refreshToken := c.Cookies("RefreshToken")
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

	userToken, err := s.db.Queries.GetUserRefreshToken(c.Context(), database.GetUserRefreshTokenParams{
		UserID: userID,
		Token:  refreshToken,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized: Refresh token not found")
	}

	if userToken.RevokedAt.Valid || time.Now().Compare(userToken.ExpiresAt) != -1 {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized: Refresh token is revoked or expired")
	}

	user, err := s.db.Queries.GetUserById(c.Context(), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized: User not found")
	}

	expTime := time.Second * 3600
	newToken, err := auth.MakeJWT(user.ID, user.Username, user.Email, JWTSecret, expTime)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate new JWT token")
	}

	return c.JSON(
		struct {
			Token string `json:"token"`
		}{
			Token: newToken,
		},
	)
}

// revokeTokenHandler godoc
//
//	@Summary		Unieważnienie refresh tokenu
//	@Description	Weryfikuje access token z nagłówka Authorization i unieważnia refresh token przekazany w cookie.
//	@Tags			auth
//	@Produce		json
//	@Param			Authorization	header	string	true	"Bearer access token"
//	@Param			Cookie			header	string	true	"RefreshToken=<refresh_token>"
//	@Success		204				"No Content"
//	@Failure		401				{object}	map[string]string
//	@Failure		500				{object}	map[string]string
//	@Router			/auth/revoke [post]
func (s *FiberServer) revokeTokenHandler(c fiber.Ctx) error {
	refreshToken := c.Cookies("RefreshToken")
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

	userToken, err := s.db.Queries.GetUserRefreshToken(c.Context(), database.GetUserRefreshTokenParams{
		UserID: userID,
		Token:  refreshToken,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized: Refresh token not found")
	}

	err = s.db.Queries.RevokeToken(c.Context(), userToken.Token)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to revoke refresh token")
	}

	return c.SendStatus(204)
}
