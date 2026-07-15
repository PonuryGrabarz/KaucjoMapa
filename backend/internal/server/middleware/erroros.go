package middleware

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/gofiber/fiber/v3"
)

func ErrorHandler(c fiber.Ctx, err error) error {
	status := fiber.StatusInternalServerError
	message := "Internal server error"

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		status = fiberErr.Code
		message = fiberErr.Message
	}

	// Błędy parsowania JSON są błędami klienta.
	var syntaxErr *json.SyntaxError
	var typeErr *json.UnmarshalTypeError

	if errors.As(err, &syntaxErr) || errors.As(err, &typeErr) {
		status = fiber.StatusBadRequest
		message = "Invalid request body"
	}

	if status >= fiber.StatusInternalServerError {
		// Pełny błąd tylko w logach, nie w odpowiedzi.
		log.Printf(
			"request failed: method=%s path=%s error=%v",
			c.Method(),
			c.Path(),
			err,
		)
	}

	return c.Status(status).JSON(fiber.Map{
		"error": message,
	})
}
