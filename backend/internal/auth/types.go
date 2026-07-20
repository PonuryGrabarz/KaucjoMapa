package auth

import "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	IsAdmin    bool   `json:"is_admin"`
	Reputation int32  `json:"reputation"`
	jwt.RegisteredClaims
}
