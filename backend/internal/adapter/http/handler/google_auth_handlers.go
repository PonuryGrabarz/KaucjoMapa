package handler

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"os"
	"time"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/adapter/http/presenter"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/auth"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
)

const (
	oauthStateCookie    = "GoogleOAuthState"
	oauthVerifierCookie = "GoogleOAuthVerifier"
)

type OAuthHandler struct {
	userSvc   usecase.IUserService
	hasherSvc usecase.IPasswordHasher
	oAuthSvc  usecase.IGoogleOAuthService
}

func NewOAuthHandler(userSvc usecase.IUserService, hasherSvc usecase.IPasswordHasher, oAuthSvc usecase.IGoogleOAuthService) *OAuthHandler {
	return &OAuthHandler{
		userSvc:   userSvc,
		hasherSvc: hasherSvc,
		oAuthSvc:  oAuthSvc,
	}
}

// GoogleLoginHandler godoc
//
//	@Summary		Logowanie przez Google
//	@Description	Generuje state i PKCE, a następnie przekierowuje do Google OAuth2.
//	@Tags			auth
//	@Success		302
//	@Failure		500	{object}	presenter.ErrorResponse	"Wewnętrzny błąd serwera"
//	@Router			/auth/google [get]
func (h *OAuthHandler) GoogleLoginHandler(c fiber.Ctx) error {
	config, err := auth.GoogleOAuthConfig()
	if err != nil {
		return entity.NewInternal(err)
	}

	state, err := randomOAuthState()
	if err != nil {
		return entity.NewInternal(err)
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
//	@Description	Weryfikuje odpowiedź Google, tworzy lub łączy konto, zwraca access token i ustawia cookie RefreshToken.
//	@Tags			auth
//	@Produce		json
//	@Param			code	query		string	true	"Authorization code"
//	@Param			state	query		string	true	"OAuth state"
//	@Success		200		{object}	presenter.TokenResponse	"Access token"
//	@Header			200		{string}	Set-Cookie				"Cookie RefreshToken; Path=/api/auth; HttpOnly; Secure; SameSite=Strict"
//	@Failure		400		{object}	presenter.ErrorResponse	"Brak authorization code lub email jest już używany"
//	@Failure		401		{object}	presenter.ErrorResponse	"Błąd OAuth Google lub wymiany kodu"
//	@Failure		403		{object}	presenter.ErrorResponse	"Nieprawidłowy state/verifier lub konto zbanowane"
//	@Failure		500		{object}	presenter.ErrorResponse	"Wewnętrzny błąd serwera"
//	@Router			/auth/google/callback [get]
func (h *OAuthHandler) GoogleCallbackHandler(c fiber.Ctx) error {
	if c.Query("error") != "" {
		clearOAuthCookie(c, oauthStateCookie)
		clearOAuthCookie(c, oauthVerifierCookie)
		return entity.NewUnauthorized("Google OAuth error: " + c.Query("error"))
	}

	cookieState := c.Cookies(oauthStateCookie)
	queryState := c.Query("state")
	verifier := c.Cookies(oauthVerifierCookie)

	clearOAuthCookie(c, oauthStateCookie)
	clearOAuthCookie(c, oauthVerifierCookie)

	if cookieState == "" || queryState == "" || subtle.ConstantTimeCompare([]byte(cookieState), []byte(queryState)) != 1 {
		return entity.NewForbidden("Invalid OAuth state")
	}

	if verifier == "" {
		return entity.NewForbidden("Missing auth verifier")
	}

	code := c.Query("code")
	if code == "" {
		return entity.NewBadRequest("Missing authorization code")
	}

	config, err := auth.GoogleOAuthConfig()
	if err != nil {
		return entity.NewInternal(err)
	}

	googleToken, err := config.Exchange(c.Context(), code, oauth2.VerifierOption(verifier))
	if err != nil {
		return entity.NewUnauthorized("Failed to exchange code for token: " + err.Error())
	}

	googleUser, err := auth.GetGoogleUser(c.Context(), config, googleToken)
	if err != nil {
		return entity.NewUnauthorized("Failed to get Google user info: " + err.Error())
	}

	user, err := h.oAuthSvc.GetUserByGoogleID(c.Context(), googleUser.Subject)
	if entity.IsKind(err, entity.KindNotFound) {
		randomPasswordHash, hashErr := h.hasherSvc.HashPassword(uuid.NewString())
		if hashErr != nil {
			return entity.NewInternal(hashErr)
		}

		user, err = h.oAuthSvc.CreateUser(c.Context(), entity.User{
			ID:           uuid.New(),
			Email:        googleUser.Email,
			Username:     googleUser.Name,
			PasswordHash: randomPasswordHash,
			GoogleID:     googleUser.Subject,
		})
		if entity.IsKind(err, entity.KindNotFound) {
			return entity.NewBadRequest("email is already used")
		}
	}

	if err != nil {
		return entity.NewInternal(err)
	}

	if user.IsBanned {
		return entity.NewForbidden("You are banned")
	}

	rToken, token, err := h.oAuthSvc.IssueSession(c.Context(), user)

	cookie := fiber.Cookie{
		Name:     "RefreshToken",
		Value:    rToken.Token,
		Path:     "/api/auth",
		Expires:  rToken.ExpiresAt,
		HTTPOnly: true,
		SameSite: "strict",
		Secure:   true,
	}

	c.Cookie(&cookie)
	return c.JSON(presenter.UserSuccessResponse(token))
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
		Path:     "/api/auth/google",
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
		Path:     "/api/auth/google",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HTTPOnly: true,
		Secure:   os.Getenv("APP_ENV") == "production",
		SameSite: "Lax",
	})
}
