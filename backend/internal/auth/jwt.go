package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
)

func MakeJWT(userID uuid.UUID, username string, email string, tokenSecret string, expiresIn time.Duration) (string, error) {
	issuedAt := jwt.NewNumericDate(time.Now())
	expiresAt := jwt.NewNumericDate(time.Now().Add(expiresIn))

	claims := JWTClaims{
		Username: username,
		Email:    email,
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

func ValidateJWT(tokenString, tokenSecret string) (JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		fmt.Println("Error parsing token:", err)
		return JWTClaims{}, err
	}

	sub, err := token.Claims.GetSubject()
	if err != nil {
		return JWTClaims{}, err
	}

	_, err = uuid.Parse(sub)
	if err != nil {
		return JWTClaims{}, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return JWTClaims{}, errors.New("invalid token claims")
	}

	fmt.Println(claims.Username)

	return *claims, nil
}

func GetBearerToken(headers map[string][]string) (string, error) {
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

func MakeRefreshToken() string {
	refreshToken := [32]byte{}
	rand.Read(refreshToken[:])

	return hex.EncodeToString(refreshToken[:])
}
