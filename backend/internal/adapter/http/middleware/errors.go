package middleware

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/adapter/http/presenter"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/gofiber/fiber/v3"
)

// kindToStatus mapuje rodzaj błędu domenowego na kod HTTP.
func kindToStatus(k entity.ErrorKind) int {
	switch k {
	case entity.KindNotFound:
		return fiber.StatusNotFound
	case entity.KindBadRequest:
		return fiber.StatusBadRequest
	case entity.KindUnauthorized, entity.KindInvalidCredentials:
		return fiber.StatusUnauthorized
	case entity.KindForbidden:
		return fiber.StatusForbidden
	case entity.KindConflict:
		return fiber.StatusConflict
	default:
		return fiber.StatusInternalServerError
	}
}

func toFieldResponses(fields []entity.FieldError) []presenter.FieldErrorResponse {
	if len(fields) == 0 {
		return nil
	}

	out := make([]presenter.FieldErrorResponse, 0, len(fields))
	for _, f := range fields {
		out = append(out, presenter.FieldErrorResponse{
			Field:   f.Field,
			Rule:    f.Rule,
			Param:   f.Param,
			Message: f.Message,
		})
	}
	return out
}

func ErrorHandler(c fiber.Ctx, err error) error {
	status := fiber.StatusInternalServerError
	resp := presenter.ErrorResponse{Error: "Internal server error"}

	var appErr *entity.AppError
	var fiberErr *fiber.Error
	var syntaxErr *json.SyntaxError
	var typeErr *json.UnmarshalTypeError

	switch {
	case errors.As(err, &appErr):
		status = kindToStatus(appErr.Kind)
		resp.Error = appErr.Message
		resp.Details = toFieldResponses(appErr.Fields)

	case errors.As(err, &fiberErr):
		// np. 404 z routera Fibera, błędy z middleware / dokumentacji.
		status = fiberErr.Code
		resp.Error = fiberErr.Message

	case errors.As(err, &syntaxErr), errors.As(err, &typeErr):
		// Błędne JSON w body to błąd klienta.
		status = fiber.StatusBadRequest
		resp.Error = "Invalid request body"
	}

	if status >= fiber.StatusInternalServerError {
		// Pełny błąd tylko w logach; klient dostaje generyk bez szczegółów.
		log.Printf(
			"request failed: method=%s path=%s error=%v",
			c.Method(),
			c.Path(),
			err,
		)
		resp.Error = "Internal server error"
		resp.Details = nil
	}

	return c.Status(status).JSON(resp)
}
