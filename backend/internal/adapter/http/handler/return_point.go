package handler

import (
	"math"
	"strconv"
	"strings"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/adapter/http/middleware"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/adapter/http/presenter"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/dto"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/port"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type ReturnPointHandler struct {
	svc     usecase.IReturnPointService
	userSvc usecase.IUserService
}

func NewReturnPointHandler(svc usecase.IReturnPointService, userSvc usecase.IUserService) *ReturnPointHandler {
	return &ReturnPointHandler{svc: svc, userSvc: userSvc}
}

// Create godoc
// @Summary Dodaj punkt zwrotu
// @Description Tworzy punkt zwrotu przypisany do zalogowanego użytkownika.
// @Tags return-points
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body presenter.CreateReturnPointRequest true "Dane punktu zwrotu"
// @Success 200 {object} presenter.ReturnPointDetailsResponse "Utworzony punkt zwrotu"
// @Failure 400 {object} presenter.ErrorResponse "Nieprawidłowe body żądania"
// @Failure 401 {object} presenter.ErrorResponse "Brak lub nieprawidłowy access token"
// @Failure 500 {object} presenter.ErrorResponse "Nie udało się utworzyć punktu zwrotu"
// @Router /return_points [post]
func (h *ReturnPointHandler) Create(c fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*port.JWTClaims)
	if !ok {
		return entity.NewUnauthorized("invalid or missing authentication")
	}

	var req presenter.CreateReturnPointRequest
	if err := c.Bind().Body(&req); err != nil {
		return err
	}

	refundMethods := []entity.RefundMethod{}
	for _, method := range req.RefundMethod {
		refundMethods = append(refundMethods, entity.RefundMethod(method))
	}

	defaultPointStyle := presenter.ResolvePointStyle(req.Name)

	pointStyle := presenter.PointStyle{}
	if req.BgColor != nil {
		pointStyle.BgColor = *req.BgColor
	} else {
		pointStyle.BgColor = defaultPointStyle.BgColor
	}

	if req.TextColor != nil {
		pointStyle.TextColor = *req.TextColor
	} else {
		pointStyle.TextColor = defaultPointStyle.TextColor
	}

	if req.ShortName != nil {
		pointStyle.ShortName = *req.ShortName
	} else {
		pointStyle.ShortName = defaultPointStyle.ShortName
	}

	returnPoint, err := h.svc.Create(c.Context(), dto.CreateReturnPointInput{
		Name:         req.Name,
		ShortName:    pointStyle.ShortName,
		BgColor:      pointStyle.BgColor,
		TextColor:    pointStyle.TextColor,
		UserID:       claims.UserID,
		Type:         req.Type,
		Address:      req.Address,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		RefundMethod: refundMethods,
		OpeningHours: entity.OpeningHours{
			MonSun: req.OpeningHours.MonSun,
			MonSat: req.OpeningHours.MonSat,
			MonFri: req.OpeningHours.MonFri,
			Sat:    req.OpeningHours.Sat,
			Sun:    req.OpeningHours.Sun,
		},
		Description:   req.Description,
		CurrentStatus: req.CurrentStatus,
		Pet:           req.Pet,
		Can:           req.Can,
		Glass:         req.Glass,
	})
	if err != nil {
		return err
	}

	return c.JSON(presenter.CreateReturnPointSuccessResponse(returnPoint))
}

// GetReturnPoints godoc
// @Summary Pobierz punkty zwrotu
// @Description Zwraca punkty zwrotu filtrowane po typie, statusie, wyszukiwanej frazie, obsługiwanych opakowaniach i odległości.
// @Tags return-points
// @Produce json
// @Security BearerAuth
// @Param type query string false "Typ punktu" Enums(butelkomat, sklep, punkt_skupu, inny)
// @Param status query string false "Status punktu" Enums(dostepny, zapelniony, nieczynny, awaria, brak_informacji, do_weryfikacji)
// @Param query query string false "Fragment nazwy lub adresu"
// @Param min_distance query number false "Minimalna odległość od użytkownika w km" minimum(0)
// @Param max_distance query number false "Maksymalna odległość od użytkownika w km" minimum(0)
// @Param lat query number false "Szerokość geograficzna użytkownika" minimum(-90) maximum(90)
// @Param lon query number false "Długość geograficzna użytkownika" minimum(-180) maximum(180)
// @Param pet query boolean false "Czy punkt przyjmuje opakowania PET"
// @Param can query boolean false "Czy punkt przyjmuje puszki"
// @Param glass query boolean false "Czy punkt przyjmuje szkło"
// @Param refund_method query string false "Metody rozliczenia rozdzielone przecinkami" example(gotowka,karta)
// @Param page query int false "Numer strony" default(1) minimum(1)
// @Param page_size query int false "Liczba elementów na stronie" default(10) minimum(1)
// @Success 200 {object} presenter.ReturnPointPage "Strona punktów zwrotu"
// @Failure 400 {object} presenter.ErrorResponse "Nieprawidłowy numer lub rozmiar strony"
// @Failure 500 {object} presenter.ErrorResponse "Nie udało się pobrać punktów zwrotu"
// @Router /return_points [get]
func (h *ReturnPointHandler) GetReturnPoints(c fiber.Ctx) error {
	token, _ := middleware.GetBearerToken(c.GetHeaders())

	pointType := c.Query("type")
	pointStatus := c.Query("status")
	query := strings.TrimSpace(c.Query("query", ""))
	minDistance, err := parseOptionalFloat(c.Query("min_distance", ""))
	if err != nil || (minDistance != nil && *minDistance < 0) {
		return entity.NewBadRequest("invalid minimum distance")
	}

	maxDistance, err := parseOptionalFloat(c.Query("max_distance", ""))
	if err != nil || (maxDistance != nil && *maxDistance < 0) {
		return entity.NewBadRequest("invalid maximum distance")
	}

	latitude, err := parseOptionalFloat(c.Query("lat", ""))
	if err != nil || (latitude != nil && (*latitude < -90 || *latitude > 90)) {
		return entity.NewBadRequest("invalid latitude")
	}

	longitude, err := parseOptionalFloat(c.Query("lon", ""))
	if err != nil || (longitude != nil && (*longitude < -180 || *longitude > 180)) {
		return entity.NewBadRequest("invalid longitude")
	}

	if (latitude == nil) != (longitude == nil) {
		return entity.NewBadRequest("latitude and longitude must be provided together")
	}
	if (minDistance != nil || maxDistance != nil) && latitude == nil {
		return entity.NewBadRequest("latitude and longitude are required for distance filters")
	}
	if minDistance != nil && maxDistance != nil && *minDistance > *maxDistance {
		return entity.NewBadRequest("minimum distance cannot exceed maximum distance")
	}

	pet, err := parseOptionalBool(c.Query("pet", ""))
	if err != nil {
		return entity.NewBadRequest("invalid pet filter")
	}
	can, err := parseOptionalBool(c.Query("can", ""))
	if err != nil {
		return entity.NewBadRequest("invalid can filter")
	}
	glass, err := parseOptionalBool(c.Query("glass", ""))
	if err != nil {
		return entity.NewBadRequest("invalid glass filter")
	}
	refundMethods, err := parseRefundMethods(c.Query("refund_method", ""))
	if err != nil {
		return entity.NewBadRequest("invalid refund method filter")
	}

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || (page < 1 && page != -1) {
		return entity.NewBadRequest("invalid page number")
	}

	pageSize, err := strconv.Atoi(c.Query("page_size", "10"))
	if err != nil || (pageSize < 1 && pageSize != -1) {
		return entity.NewBadRequest("invalid page size")
	}

	filters := usecase.Filter{
		Query:         query,
		Type:          entity.PointType(pointType),
		Status:        entity.PointStatus(pointStatus),
		Latitude:      latitude,
		Longitude:     longitude,
		MinDistanceKm: minDistance,
		MaxDistanceKm: maxDistance,
		Pet:           pet,
		Can:           can,
		Glass:         glass,
		RefundMethods: refundMethods,
	}

	listSize, err := h.svc.ListSize(c.Context(), filters)
	if err != nil {
		return err
	}

	// -1 na page lub page_size oznacza "zwróć wszystko".
	returnAll := page == -1 || pageSize == -1

	var maxPage int64
	if returnAll {
		filters.Offset = 0
		filters.Limit = int32(listSize)
		maxPage = 1
		page = 1
		pageSize = int(listSize)
	} else {
		filters.Offset = int32((page - 1) * pageSize)
		filters.Limit = int32(pageSize)
		maxPage = (listSize + int64(pageSize) - 1) / int64(pageSize)
	}

	returnPoints, err := h.svc.List(c.Context(), token, filters)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(presenter.GetReturnPointSuccessResponse(returnPoints, maxPage, int64(page), int64(pageSize), listSize))
}

func parseOptionalFloat(value string) (*float64, error) {
	if value == "" {
		return nil, nil
	}

	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil || math.IsNaN(parsed) || math.IsInf(parsed, 0) {
		if err == nil {
			err = strconv.ErrSyntax
		}
		return nil, err
	}

	return &parsed, nil
}

func parseOptionalBool(value string) (*bool, error) {
	if value == "" {
		return nil, nil
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}

func parseRefundMethods(value string) ([]entity.RefundMethod, error) {
	if value == "" {
		return nil, nil
	}

	methods := make([]entity.RefundMethod, 0)
	seen := make(map[entity.RefundMethod]struct{})
	for _, value := range strings.Split(value, ",") {
		method := entity.RefundMethod(strings.TrimSpace(value))
		switch method {
		case entity.RefundMethodCash, entity.RefundMethodVoucher, entity.RefundMethodCard:
		default:
			return nil, strconv.ErrSyntax
		}

		if _, exists := seen[method]; exists {
			continue
		}
		seen[method] = struct{}{}
		methods = append(methods, method)
	}

	return methods, nil
}

// DeletePoint godoc
// @Summary Usuń punkt zwrotu
// @Description Usuwa punkt zwrotu. Operacja jest dostępna tylko dla administratora.
// @Tags return-points
// @Produce json
// @Security BearerAuth
// @Param id path string true "Identyfikator punktu" format(uuid)
// @Success 204 "No Content"
// @Failure 400 {object} presenter.ErrorResponse "Nieprawidłowy identyfikator punktu"
// @Failure 401 {object} presenter.ErrorResponse "Brak lub nieprawidłowy access token"
// @Failure 500 {object} presenter.ErrorResponse "Nie udało się usunąć punktu zwrotu"
// @Router /return_points/{id} [delete]
func (h *ReturnPointHandler) DeletePoint(c fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*port.JWTClaims)
	if !ok {
		return entity.NewUnauthorized("invalid or missing authentication")
	}

	mapPointID := c.Params("id")
	mapPointUUID, err := uuid.Parse(mapPointID)
	if err != nil {
		return entity.NewBadRequest("invalid return point ID")
	}

	err = h.svc.Delete(c.Context(), claims.UserID, mapPointUUID)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// UpdatePoint godoc
// @Summary Zaktualizuj punkt zwrotu
// @Description Częściowo aktualizuje dane istniejącego punktu zwrotu.
// @Tags return-points
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Identyfikator punktu" format(uuid)
// @Param request body presenter.UpdateKaucjomatRequest true "Pola do aktualizacji"
// @Success 200 {object} presenter.ReturnPointDetailsResponse "Zaktualizowany punkt zwrotu"
// @Failure 400 {object} presenter.ErrorResponse "Nieprawidłowy identyfikator lub body żądania"
// @Failure 401 {object} presenter.ErrorResponse "Brak lub nieprawidłowy access token"
// @Failure 500 {object} presenter.ErrorResponse "Nie udało się pobrać użytkownika lub zaktualizować punktu"
// @Router /return_points/{id} [put]
func (h *ReturnPointHandler) UpdatePoint(c fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*port.JWTClaims)
	if !ok {
		return entity.NewUnauthorized("invalid or missing authentication")
	}

	mapPointID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return entity.NewBadRequest("invalid return point ID")
	}

	user, err := h.userSvc.GetUserByID(c.Context(), claims.UserID)
	if err != nil {
		return err
	}

	point, err := h.svc.GetByID(c.Context(), mapPointID)
	if err != nil {
		return err
	}

	var req presenter.UpdateKaucjomatRequest
	if err := c.Bind().Body(&req); err != nil {
		return err
	}

	params := entity.ReturnPoint{
		ID:    mapPointID,
		Pet:   point.Pet,
		Can:   point.Can,
		Glass: point.Glass,
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

	if req.RefundMethod != nil {
		refundMethods := []entity.RefundMethod{}
		for _, method := range *req.RefundMethod {
			refundMethods = append(refundMethods, entity.RefundMethod(method))
		}
		params.RefundMethod = refundMethods
	}

	if req.OpeningHours != nil {
		params.OpeningHours = entity.OpeningHours{
			MonSun: req.OpeningHours.MonSun,
			MonSat: req.OpeningHours.MonSat,
			MonFri: req.OpeningHours.MonFri,
			Sat:    req.OpeningHours.Sat,
			Sun:    req.OpeningHours.Sun,
		}
	}

	if req.Description != nil {
		params.Description = *req.Description
	}

	updatedPoint, err := h.svc.Update(c.Context(), user, params)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(presenter.CreateReturnPointSuccessResponse(updatedPoint))
}
