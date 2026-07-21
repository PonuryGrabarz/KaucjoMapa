package presenter

import "github.com/gofiber/fiber/v3"

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginSuccessResponse(token string) *fiber.Map {
	return &fiber.Map{
		"token": token,
	}
}

func RefreshTokenSuccessResponse(token string) *fiber.Map {
	return &fiber.Map{
		"token": token,
	}
}
