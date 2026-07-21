package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	IsAdmin    bool   `json:"is_admin"`
	Reputation int32  `json:"reputation"`
	jwt.RegisteredClaims
}

type ITokenSigner interface {
	MakeJWT(userID uuid.UUID, isAdmin bool, reputation int32, username, email string) (string, error)
	MakeRefreshToken() string
	ValidateJWT(tokenString string) (JWTClaims, error)
	GetBearerToken(headers map[string][]string) (string, error)
}
