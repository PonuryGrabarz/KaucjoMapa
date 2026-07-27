package token

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/port"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenSigner struct {
	JWTSecret      string
	ExpirationTime time.Duration
}

type claims struct {
	port.JWTClaims
	jwt.RegisteredClaims
}

var _ port.ITokenSigner = (*TokenSigner)(nil)

func NewTokenSigner(jwtSecret string, expTime time.Duration) *TokenSigner {
	return &TokenSigner{
		JWTSecret:      jwtSecret,
		ExpirationTime: expTime,
	}
}

func (t *TokenSigner) MakeJWT(userID uuid.UUID, isAdmin bool, reputation int32, username, email string) (string, error) {
	issuedAt := jwt.NewNumericDate(time.Now())
	expiresAt := jwt.NewNumericDate(time.Now().Add(t.ExpirationTime))

	claims := claims{
		port.JWTClaims{
			UserID:     userID,
			Username:   username,
			Email:      email,
			IsAdmin:    isAdmin,
			Reputation: reputation,
		},
		jwt.RegisteredClaims{
			Issuer:    "KaucjoMapa",
			IssuedAt:  issuedAt,
			ExpiresAt: expiresAt,
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(t.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (t *TokenSigner) MakeRefreshToken() string {
	refreshToken := [32]byte{}
	rand.Read(refreshToken[:])

	return hex.EncodeToString(refreshToken[:])
}

func (t *TokenSigner) ValidateJWT(tokenString string) (port.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (any, error) {
		return []byte(t.JWTSecret), nil
	})
	if err != nil {
		return port.JWTClaims{}, err
	}

	sub, err := token.Claims.GetSubject()
	if err != nil {
		return port.JWTClaims{}, err
	}

	_, err = uuid.Parse(sub)
	if err != nil {
		return port.JWTClaims{}, err
	}

	claims, ok := token.Claims.(*claims)
	if !ok {
		return port.JWTClaims{}, errors.New("invalid token claims")
	}

	return port.JWTClaims{
		UserID:     claims.UserID,
		Username:   claims.Username,
		IsAdmin:    claims.IsAdmin,
		Email:      claims.Email,
		Reputation: claims.Reputation,
	}, nil
}

func (t *TokenSigner) GetBearerToken(headers map[string][]string) (string, error) {
	authHeaders := headers["Authorization"]
	if authHeaders == nil {
		return "", errors.New("invalid authorization header")
	}

	for _, auth := range authHeaders {
		if !strings.Contains(auth, "Bearer ") {
			continue
		}

		token := strings.TrimPrefix(auth, "Bearer ")

		return token, nil
	}

	return "", errors.New("invalid authorization header")
}

func (t *TokenSigner) ExtractClaimsFromJWT(tokenString string) (port.JWTClaims, error) {
	jwtClaims := claims{}

	token, _, err := jwt.NewParser().ParseUnverified(tokenString, &jwtClaims)
	if err != nil {
		return port.JWTClaims{}, err
	}

	claims, ok := token.Claims.(*claims)
	if !ok {
		return port.JWTClaims{}, errors.New("invalid token claims")
	}

	return port.JWTClaims{
		UserID:     claims.UserID,
		Username:   claims.Username,
		Email:      claims.Email,
		IsAdmin:    claims.IsAdmin,
		Reputation: claims.Reputation,
	}, nil
}
