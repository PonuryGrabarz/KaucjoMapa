package presenter

import "github.com/gofiber/fiber/v3"

type CreateUserRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func UserSuccessResponse(token string) *fiber.Map {
	return &fiber.Map{
		"token": token,
	}
}
