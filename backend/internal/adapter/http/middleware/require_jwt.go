package middleware

import (
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/auth"
	"github.com/gofiber/fiber/v3"
)

func RequireJWT(auth auth.ITokenSigner, secret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		token, err := auth.GetBearerToken(c.GetReqHeaders())
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Missing token")
		}

		claims, err := auth.ValidateJWT(token)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
		}

		c.Locals("claims", &claims)
		return c.Next()
	}
}
