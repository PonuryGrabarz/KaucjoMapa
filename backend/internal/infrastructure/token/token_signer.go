package token

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenSigner struct{}

var _ auth.ITokenSigner = (*TokenSigner)(nil)

func (t *TokenSigner) MakeJWT(userID uuid.UUID, expiresIn time.Duration, isAdmin bool, reputation int32, username, email, tokenSecret string) (string, error) {
	issuedAt := jwt.NewNumericDate(time.Now())
	expiresAt := jwt.NewNumericDate(time.Now().Add(expiresIn))

	claims := auth.JWTClaims{
		Username:   username,
		Email:      email,
		IsAdmin:    isAdmin,
		Reputation: reputation,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "KaucjoMapa",
			IssuedAt:  issuedAt,
			ExpiresAt: expiresAt,
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(tokenSecret))
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

func (t *TokenSigner) ValidateJWT(tokenString, tokenSecret string) (auth.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &auth.JWTClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return auth.JWTClaims{}, err
	}

	sub, err := token.Claims.GetSubject()
	if err != nil {
		return auth.JWTClaims{}, err
	}

	_, err = uuid.Parse(sub)
	if err != nil {
		return auth.JWTClaims{}, err
	}

	claims, ok := token.Claims.(*auth.JWTClaims)
	if !ok {
		return auth.JWTClaims{}, errors.New("invalid token claims")
	}

	fmt.Println(claims.Username)

	return *claims, nil
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
