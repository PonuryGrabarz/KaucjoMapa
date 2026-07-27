package middleware

import (
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/port"
	"github.com/gofiber/fiber/v3"
)

func ExtractJWT(auth port.ITokenSigner, secret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		token, err := GetBearerToken(c.GetReqHeaders())
		if err != nil {
			return entity.NewUnauthorized("missing or malformed authorization header")
		}

		claims, err := auth.ExtractClaimsFromJWT(token)
		if err != nil {
			return entity.NewUnauthorized("invalid token")
		}

		c.Locals("claims", &claims)
		return c.Next()
	}
}
