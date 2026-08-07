package handler

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	auth "github.com/PonuryGrabarz/KaucjoMapa/backend/internal/infrastructure/oAuth"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/port"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/encryptcookie"
	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
)

const (
	oauthStateCookie    = "GoogleOAuthState"
	oauthVerifierCookie = "GoogleOAuthVerifier"

	HandoffAccessCookie  = "OAuthAccess"
	HandoffRefreshCookie = "OAuthRefresh"

	refreshCookieName = "RefreshToken"

	handoffCookieTTL = 60 * time.Second

	loginLandingPath    = "/logowanie/google"
	registerLandingPath = "/rejestracja/google"
	loginErrorPath      = "/logowanie?blad=google"
	registerErrorPath   = "/rejestracja?blad=google"
)

type OAuthHandler struct {
	userSvc      usecase.IUserService
	hasherSvc    usecase.IPasswordHasher
	oAuthSvc     usecase.IGoogleOAuthService
	tokenSvc     port.ITokenSigner
	frontendURL  string
	cookieSecret string
}

func NewOAuthHandler(userSvc usecase.IUserService, hasherSvc usecase.IPasswordHasher, oAuthSvc usecase.IGoogleOAuthService, tokenSvc port.ITokenSigner, frontendURL, cookieSecret string) *OAuthHandler {
	return &OAuthHandler{
		userSvc:      userSvc,
		hasherSvc:    hasherSvc,
		oAuthSvc:     oAuthSvc,
		tokenSvc:     tokenSvc,
		frontendURL:  frontendURL,
		cookieSecret: cookieSecret,
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
	c.Cookie(newOAuthCookie(c, oauthStateCookie, state))
	c.Cookie(newOAuthCookie(c, oauthVerifierCookie, verifier))

	loginURL := config.AuthCodeURL(
		state,
		oauth2.S256ChallengeOption(verifier),
	)

	return c.Redirect().To(loginURL)
}

// GoogleCallbackHandler godoc
//
//	@Summary		Callback Google OAuth2
//	@Description	Weryfikuje odpowiedź Google, tworzy lub łączy konto, po czym przekierowuje na frontend.
//	@Description	Tokeny wędrują w krótkożyjących cookies OAuthAccess i OAuthRefresh (Path=/, SameSite=Lax),
//	@Description	które strona lądowania frontendu odczytuje raz i kasuje. Każdy błąd również kończy się
//	@Description	przekierowaniem — na /logowanie?blad=google lub /rejestracja?blad=google.
//	@Tags			auth
//	@Param			code	query	string	true	"Authorization code"
//	@Param			state	query	string	true	"OAuth state"
//	@Success		303		"Przekierowanie na /logowanie/google lub /rejestracja/google"
//	@Header			303		{string}	Set-Cookie	"OAuthAccess i OAuthRefresh; Path=/; HttpOnly; SameSite=Lax; Max-Age=60"
//	@Router			/auth/google/callback [get]
func (h *OAuthHandler) GoogleCallbackHandler(c fiber.Ctx) error {
	cookieState := c.Cookies(oauthStateCookie)
	queryState := c.Query("state")
	verifier := c.Cookies(oauthVerifierCookie)

	clearOAuthCookie(c, oauthStateCookie)
	clearOAuthCookie(c, oauthVerifierCookie)

	if authError := c.Query("error"); authError != "" {
		return h.redirectFailure(c, loginErrorPath, fmt.Errorf("google oauth error: %s", authError))
	}

	if cookieState == "" || queryState == "" || subtle.ConstantTimeCompare([]byte(cookieState), []byte(queryState)) != 1 {
		return h.redirectFailure(c, loginErrorPath, errors.New("invalid oauth state"))
	}

	if verifier == "" {
		return h.redirectFailure(c, loginErrorPath, errors.New("missing auth verifier"))
	}

	code := c.Query("code")
	if code == "" {
		return h.redirectFailure(c, loginErrorPath, errors.New("missing authorization code"))
	}

	config, err := auth.GoogleOAuthConfig()
	if err != nil {
		return h.redirectFailure(c, loginErrorPath, err)
	}

	googleToken, err := config.Exchange(c.Context(), code, oauth2.VerifierOption(verifier))
	if err != nil {
		return h.redirectFailure(c, loginErrorPath, fmt.Errorf("failed to exchange code for token: %w", err))
	}

	googleUser, err := auth.GetGoogleUser(c.Context(), config, googleToken)
	if err != nil {
		return h.redirectFailure(c, loginErrorPath, fmt.Errorf("failed to get google user info: %w", err))
	}

	user, err := h.oAuthSvc.GetUserByGoogleID(c.Context(), googleUser.Subject)

	isRegistration := entity.IsKind(err, entity.KindNotFound)

	errorPath := loginErrorPath
	if isRegistration {
		errorPath = registerErrorPath
	}

	if isRegistration {
		randomPasswordHash, hashErr := h.hasherSvc.HashPassword(uuid.NewString())
		if hashErr != nil {
			return h.redirectFailure(c, errorPath, hashErr)
		}

		user, err = h.oAuthSvc.CreateUser(c.Context(), entity.User{
			ID:           uuid.New(),
			Email:        googleUser.Email,
			Username:     googleUser.Name,
			PasswordHash: randomPasswordHash,
			GoogleID:     googleUser.Subject,
		})
	}

	if err != nil {
		return h.redirectFailure(c, errorPath, err)
	}

	if user.IsBanned {
		return h.redirectFailure(c, errorPath, fmt.Errorf("banned user %s attempted google login", user.ID))
	}

	rToken, accessToken, err := h.oAuthSvc.IssueSession(c.Context(), user)
	if err != nil {
		return h.redirectFailure(c, errorPath, fmt.Errorf("failed to issue session: %w", err))
	}

	if err := h.issueHandoff(c, accessToken, rToken.Token); err != nil {
		return h.redirectFailure(c, errorPath, err)
	}

	if isRegistration {
		return h.redirectToFrontend(c, registerLandingPath)
	}

	return h.redirectToFrontend(c, loginLandingPath)
}

func (h *OAuthHandler) redirectToFrontend(c fiber.Ctx, path string) error {
	return c.Redirect().Status(fiber.StatusSeeOther).To(h.frontendURL + path)
}

func (h *OAuthHandler) redirectFailure(c fiber.Ctx, path string, cause error) error {
	log.Printf("google oauth callback failed: path=%s error=%v", c.Path(), cause)

	return h.redirectToFrontend(c, path)
}

func randomOAuthState() (string, error) {
	value := make([]byte, 32)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}

	return hex.EncodeToString(value), nil
}

func (h *OAuthHandler) issueHandoff(c fiber.Ctx, accessToken, refreshToken string) error {
	encryptedRefresh, err := encryptcookie.EncryptCookie(refreshCookieName, refreshToken, h.cookieSecret)
	if err != nil {
		return fmt.Errorf("failed to encrypt handoff refresh token: %w", err)
	}

	c.Cookie(newHandoffCookie(c, HandoffAccessCookie, accessToken))
	c.Cookie(newHandoffCookie(c, HandoffRefreshCookie, encryptedRefresh))

	return nil
}

func newHandoffCookie(c fiber.Ctx, name, value string) *fiber.Cookie {
	return &fiber.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Expires:  time.Now().Add(handoffCookieTTL),
		MaxAge:   int(handoffCookieTTL.Seconds()),
		HTTPOnly: true,
		Secure:   secureCookies(c),
		SameSite: "Lax",
	}
}

func secureCookies(c fiber.Ctx) bool {
	return os.Getenv("APP_ENV") == "production" || c.Scheme() == "https"
}

func newOAuthCookie(c fiber.Ctx, name, value string) *fiber.Cookie {
	return &fiber.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/api/auth/google",
		MaxAge:   600,
		Expires:  time.Now().Add(10 * time.Minute),
		HTTPOnly: true,
		Secure:   secureCookies(c),
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
		Secure:   secureCookies(c),
		SameSite: "Lax",
	})
}
