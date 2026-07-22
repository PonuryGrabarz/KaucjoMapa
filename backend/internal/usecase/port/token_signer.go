package port

import "github.com/google/uuid"

type JWTClaims struct {
	UserID     uuid.UUID `json:"user_id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	IsAdmin    bool      `json:"is_admin"`
	Reputation int32     `json:"reputation"`
}

type ITokenSigner interface {
	MakeJWT(userID uuid.UUID, isAdmin bool, reputation int32, username, email string) (string, error)
	MakeRefreshToken() string
	ValidateJWT(tokenString string) (JWTClaims, error)
}
