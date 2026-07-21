package handler

import (
	"strconv"
	"strings"
	"time"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/adapter/http/presenter"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/auth"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/returnpoint"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/user"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type ReturnPointHandler struct {
	svc     returnpoint.IReturnPointService
	userSvc user.IUserService
}

func NewReturnPointHandler(svc returnpoint.IReturnPointService, userSvc user.IUserService) *ReturnPointHandler {
	return &ReturnPointHandler{svc: svc, userSvc: userSvc}
}

func (h *ReturnPointHandler) Create(c fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*auth.JWTClaims)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid JWT claims")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid user ID in JWT claims")
	}

	var req presenter.CreateReturnPointRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	point, err := h.svc.Create(c.Context(), userID, entity.ReturnPoint{
		Name:          req.Name,
		Type:          entity.PointType(req.Type),
		Address:       req.Address,
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
		OpeningHours:  req.OpeningHours,
		Description:   req.Description,
		CurrentStatus: entity.PointStatus(req.CurrentStatus),
		IsVerified:    false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create return point")
	}
	return c.JSON(presenter.CreateReturnPointSuccessResponse(point))
}

func (h *ReturnPointHandler) GetReturnPoints(c fiber.Ctx) error {
	pointType := c.Query("type")
	pointStatus := c.Query("status")
	pageQuery := c.Query("page", "1")
	pageSizeQuery := c.Query("page_size", "10")
	query := strings.TrimSpace(c.Query("query", ""))

	page, err := strconv.Atoi(pageQuery)
	if err != nil || page < 1 {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid page number")
	}
	pageSize, err := strconv.Atoi(pageSizeQuery)
	if err != nil || pageSize < 1 {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid page size")
	}

	returnPoints, err := h.svc.List(c.Context(), returnpoint.Filter{
		Query:  query,
		Type:   entity.PointType(pointType),
		Status: entity.PointStatus(pointStatus),
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch return points")
	}

	mapPointsDetails := []*fiber.Map{}

	for _, point := range returnPoints {
		mapPointsDetails = append(mapPointsDetails, presenter.CreateReturnPointSuccessResponse(point))
	}

	maxPage := len(mapPointsDetails) / pageSize

	if maxPage > 0 && page <= maxPage {
		return c.JSON(mapPointsDetails[pageSize*(page-1) : min(pageSize*page, len(mapPointsDetails))])
	}

	return c.JSON(mapPointsDetails)
}

func (h *ReturnPointHandler) DeletePoint(c fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*auth.JWTClaims)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid JWT claims")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Invalid user ID in JWT claims")
	}

	user, err := h.userSvc.GetUserByID(c.Context(), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Faild to fetch user")
	}

	if !user.IsAdmin {
		return fiber.NewError(fiber.StatusForbidden, "User is not an admin")
	}
	mapPointID := c.Params("id")
	mapPointUUID, err := uuid.Parse(mapPointID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid map point ID")
	}

	err = h.svc.Delete(c.Context(), user, mapPointUUID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete return point")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *ReturnPointHandler) UpdatePoint(c fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*auth.JWTClaims)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid JWT claims")
	}

	mapPointID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid map point ID")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Invalid user ID in JWT claims")
	}

	user, err := h.userSvc.GetUserByID(c.Context(), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch user")
	}

	_, err = h.svc.GetByID(c.Context(), mapPointID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Return point not found")
	}

	var req presenter.UpdateKaucjomatRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	params := entity.ReturnPoint{
		ID: mapPointID,
	}

	if req.Name != nil {
		params.Name = *req.Name
	}

	if req.Type != nil {
		params.Type = entity.PointType(*req.Type)
	}

	if req.Address != nil {
		params.Address = *req.Address
	}

	if req.Latitude != nil {
		params.Latitude = *req.Latitude
	}

	if req.Longitude != nil {
		params.Longitude = *req.Longitude
	}

	if req.Pet != nil {
		params.Pet = *req.Pet
	}

	if req.Can != nil {
		params.Can = *req.Can
	}

	if req.Glass != nil {
		params.Glass = *req.Glass
	}

	if req.OpeningHours != nil {
		params.OpeningHours = *req.OpeningHours
	}

	if req.Description != nil {
		params.Description = *req.Description
	}

	updatedPoint, err := h.svc.Update(c.Context(), user, params)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update return point")
	}

	return c.Status(fiber.StatusOK).JSON(updatedPoint)
}
