package handler

import (
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/adapter/http/presenter"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/returnpoint"
	"github.com/gofiber/fiber/v3"
)

type ReturnPointHandler struct {
	svc returnpoint.IReturnPointService
}

func NewReturnPointHandler(svc returnpoint.IReturnPointService) *ReturnPointHandler {
	return &ReturnPointHandler{svc: svc}
}

func (h *ReturnPointHandler) Create(c fiber.Ctx) error {
	userID, err := userIDFromCtx(c) // z claims
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid JWT claims")
	}
	var req presenter.CreateKaucjomatRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	point, err := h.svc.Create(c.Context(), userID, req.ToEntity())
	if err != nil {
		return presenter.MapError(err) // entity.ErrPointExistsNearby -> 400 itd.
	}
	return c.JSON(presenter.ToMapPoint(point))
}
