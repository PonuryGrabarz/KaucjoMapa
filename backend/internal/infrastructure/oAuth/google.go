package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	_ "github.com/joho/godotenv/autoload"
)

const googleUserInfoURL = "https://openidconnect.googleapis.com/v1/userinfo"

type GoogleUser struct {
	Subject       string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

func GoogleOAuthConfig() (*oauth2.Config, error) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	if clientID == "" || clientSecret == "" || redirectURL == "" {
		return nil, fmt.Errorf("google OAuth configuration is incomplete")
	}

	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"openid",
			"email",
			"profile",
		},
		Endpoint: google.Endpoint,
	}, nil
}

func GetGoogleUser(ctx context.Context, config *oauth2.Config, token *oauth2.Token) (GoogleUser, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, googleUserInfoURL, nil)
	if err != nil {
		return GoogleUser{}, fmt.Errorf("create Google userinfo request: %w", err)
	}

	response, err := config.Client(ctx, token).Do(request)
	if err != nil {
		return GoogleUser{}, fmt.Errorf("fetch Google userinfo: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(response.Body, 4096))
		return GoogleUser{}, fmt.Errorf("google userinfo returned %s: %s", response.Status, body)
	}

	var user GoogleUser
	if err := json.NewDecoder(response.Body).Decode(&user); err != nil {
		return GoogleUser{}, fmt.Errorf("decode Google userinfo: %w", err)
	}

	if user.Subject == "" || user.Email == "" {
		return GoogleUser{}, fmt.Errorf("google response is missing user identity")
	}

	if !user.EmailVerified {
		return GoogleUser{}, fmt.Errorf("google email is not verified")
	}

	return user, nil
}
