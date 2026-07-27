package server

import (
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/hex"
	"errors"
	"os"
	"time"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/auth"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/database"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

const (
	oauthStateCookie    = "GoogleOAuthState"
	oauthVerifierCookie = "GoogleOAuthVerifier"
)

// GoogleLoginHandler godoc
//
//	@Summary		Logowanie przez Google
//	@Description	Generuje state i PKCE, a następnie przekierowuje do Google OAuth2.
//	@Tags			auth
//	@Success		302
//	@Failure		500	{object}	map[string]string
//	@Router			/auth/google [get]
func (s *FiberServer) GoogleLoginHandler(c fiber.Ctx) error {
	config, err := auth.GoogleOAuthConfig()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Google OAuth is not configured")
	}

	state, err := randomOAuthState()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate OAuth state")
	}

	verifier := oauth2.GenerateVerifier()
	c.Cookie(newOAuthCookie(oauthStateCookie, state))
	c.Cookie(newOAuthCookie(oauthVerifierCookie, verifier))

	loginURL := config.AuthCodeURL(
		state,
		oauth2.S256ChallengeOption(verifier),
	)

	return c.Redirect().To(loginURL)
}

// GoogleCallbackHandler godoc
//
//	@Summary		Callback Google OAuth2
//	@Description	Weryfikuje odpowiedź Google, tworzy lub łączy konto i zwraca użytkownika z lokalnym JWT.
//	@Tags			auth
//	@Produce		json
//	@Param			code	query		string	true	"Authorization code"
//	@Param			state	query		string	true	"OAuth state"
//	@Success		200		{object}	User
//	@Header			200		{string}	Set-Cookie	"RefreshToken"
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Failure		403		{object}	map[string]string
//	@Failure		409		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/auth/google/callback [get]
func (s *FiberServer) GoogleCallbackHandler(c fiber.Ctx) error {
	if c.Query("error") != "" {
		clearOAuthCookie(c, oauthStateCookie)
		clearOAuthCookie(c, oauthVerifierCookie)
		return fiber.NewError(fiber.StatusUnauthorized, "Google login was cancelled")
	}

	cookieState := c.Cookies(oauthStateCookie)
	queryState := c.Query("state")
	verifier := c.Cookies(oauthVerifierCookie)

	clearOAuthCookie(c, oauthStateCookie)
	clearOAuthCookie(c, oauthVerifierCookie)

	if cookieState == "" || queryState == "" || subtle.ConstantTimeCompare([]byte(cookieState), []byte(queryState)) != 1 {
		return fiber.NewError(fiber.StatusForbidden, "Invalid OAuth state")
	}

	if verifier == "" {
		return fiber.NewError(fiber.StatusForbidden, "Missing OAuth verifier")
	}

	code := c.Query("code")
	if code == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing authorization code")
	}

	config, err := auth.GoogleOAuthConfig()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Google OAuth is not configured")
	}

	googleToken, err := config.Exchange(c.Context(), code, oauth2.VerifierOption(verifier))
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Failed to exchange Google authorization code")
	}

	googleUser, err := auth.GetGoogleUser(c.Context(), config, googleToken)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Failed to verify Google account")
	}

	googleSubject := sql.NullString{String: googleUser.Subject, Valid: true}
	user, err := s.db.Queries.GetUserByGoogleSubject(c.Context(), googleSubject)
	if errors.Is(err, sql.ErrNoRows) {
		randomPasswordHash, hashErr := auth.HashPassword(auth.MakeRefreshToken())
		if hashErr != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to prepare Google account")
		}

		user, err = s.db.Queries.UpsertGoogleUserByEmail(c.Context(), database.UpsertGoogleUserByEmailParams{
			ID:            uuid.New(),
			Email:         googleUser.Email,
			Username:      "google_" + uuid.NewString(),
			PasswordHash:  randomPasswordHash,
			GoogleSubject: googleSubject,
		})
		if errors.Is(err, sql.ErrNoRows) {
			return fiber.NewError(fiber.StatusConflict, "Email is already linked to another Google account")
		}
	}

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to load Google account")
	}

	if user.IsBanned {
		return fiber.NewError(fiber.StatusForbidden, "User account is banned")
	}

	return s.issueSession(c, user)
}

func randomOAuthState() (string, error) {
	value := make([]byte, 32)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}

	return hex.EncodeToString(value), nil
}

func newOAuthCookie(name, value string) *fiber.Cookie {
	return &fiber.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/auth/google",
		MaxAge:   600,
		Expires:  time.Now().Add(10 * time.Minute),
		HTTPOnly: true,
		Secure:   os.Getenv("APP_ENV") == "production",
		SameSite: "Lax",
	}
}

func clearOAuthCookie(c fiber.Ctx, name string) {
	c.Cookie(&fiber.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/auth/google",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HTTPOnly: true,
		Secure:   os.Getenv("APP_ENV") == "production",
		SameSite: "Lax",
	})
}
