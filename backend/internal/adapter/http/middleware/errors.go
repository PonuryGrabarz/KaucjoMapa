package middleware

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/gofiber/fiber/v3"
)

// domainErrorStatus mapuje błędy domenowe na kody HTTP.
func domainErrorStatus(err error) (int, string, bool) {
	switch {
	case errors.Is(err, entity.ErrNotFound):
		return fiber.StatusNotFound, "Resource not found", true
	case errors.Is(err, entity.ErrForbidden):
		return fiber.StatusForbidden, "Forbidden", true
	case errors.Is(err, entity.ErrUnauthorized):
		return fiber.StatusUnauthorized, "Unauthorized", true
	case errors.Is(err, entity.ErrInvalidCredentials):
		return fiber.StatusUnauthorized, "Invalid credentials", true
	case errors.Is(err, entity.ErrBadRequest):
		return fiber.StatusBadRequest, "Bad request", true
	case errors.Is(err, entity.ErrPointExistsNearby):
		return fiber.StatusConflict, "A return point already exists nearby", true
	default:
		return 0, "", false
	}
}

func ErrorHandler(c fiber.Ctx, err error) error {
	status := fiber.StatusInternalServerError
	message := "Internal server error"

	// Błędy domenowe mają pierwszeństwo nad domyślnym 500.
	if s, m, ok := domainErrorStatus(err); ok {
		status = s
		message = m
	}

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
