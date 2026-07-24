package middleware

import (
	"errors"
	"strings"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/port"
	"github.com/gofiber/fiber/v3"
)

func GetBearerToken(headers map[string][]string) (string, error) {
	authHeaders := headers["Authorization"]
	if authHeaders == nil {
		return "", errors.New("invalid authorization header")
	}

	for _, auth := range authHeaders {
		if !strings.HasPrefix(auth, "Bearer ") {
			continue
		}

		token := strings.TrimPrefix(auth, "Bearer ")
		if token == "" {
			continue
		}

		return token, nil
	}

	return "", errors.New("invalid authorization header")
}

func RequireJWT(auth port.ITokenSigner, secret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		token, err := GetBearerToken(c.GetReqHeaders())
		if err != nil {
			return entity.NewUnauthorized("missing or malformed authorization header")
		}

		claims, err := auth.ValidateJWT(token)
		if err != nil {
			return entity.NewUnauthorized("invalid token")
		}

		c.Locals("claims", &claims)
		return c.Next()
	}
}
