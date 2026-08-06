package handler

import (
	"time"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/adapter/http/presenter"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/port"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type DataErrorReportHandler struct {
	errorReportService usecase.IDataErrorReportsService
}

func NewDataErrorReportHandler(service usecase.IDataErrorReportsService) *DataErrorReportHandler {
	return &DataErrorReportHandler{
		errorReportService: service,
	}
}

// Create godoc
// @Summary Utwórz zgłoszenie błędu danych
// @Description Tworzy zgłoszenie błędu danych dotyczące wskazanego punktu zwrotu, przypisane do zalogowanego użytkownika.
// @Tags data-error-reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param return_point_id path string true "Identyfikator punktu zwrotu" Format(uuid)
// @Param request body presenter.CreateDataErrorReportRequest true "Dane zgłoszenia błędu"
// @Success 201 {object} presenter.DataErrorReportSuccessResponse "Utworzone zgłoszenie błędu danych"
// @Failure 400 {object} presenter.ErrorResponse "Nieprawidłowe body żądania lub identyfikator punktu zwrotu"
// @Failure 401 {object} presenter.ErrorResponse "Brak lub nieprawidłowy access token"
// @Failure 500 {object} presenter.ErrorResponse "Nie udało się utworzyć zgłoszenia"
// @Router /return_points/{return_point_id}/error_reports [post]
func (h *DataErrorReportHandler) Create(c fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*port.JWTClaims)
	if !ok {
		return entity.NewUnauthorized("invalid or missing authentication")
	}

	var req presenter.CreateDataErrorReportRequest
	if err := c.Bind().Body(&req); err != nil {
		return err
	}

	returnPointID, err := uuid.Parse(c.Params("return_point_id"))
	if err != nil {
		return entity.NewBadRequest("invalid return point id")
	}

	categories := []entity.Category{}
	for _, cat := range req.Categories {
		categories = append(categories, entity.Category(cat))
	}

	newReport, err := h.errorReportService.Create(c.Context(), entity.ErrorReport{
		ID:                  uuid.New(),
		UserID:              &claims.UserID,
		ReturnPointID:       &returnPointID,
		Categories:          categories,
		Description:         req.Description,
		SuggestedCorrection: req.SuggestedCorrection,
		Status:              entity.ErrorReportStatusPending,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(presenter.CreateDataErrorReportSuccessResponse(newReport))
}

// List godoc
// @Summary Pobierz zgłoszenia błędów danych
// @Description Zwraca listę wszystkich zgłoszeń błędów danych.
// @Tags data-error-reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} presenter.DataErrorReportSuccessResponse "Lista zgłoszeń błędów danych"
// @Failure 401 {object} presenter.ErrorResponse "Brak lub nieprawidłowy access token"
// @Failure 500 {object} presenter.ErrorResponse "Nie udało się pobrać zgłoszeń"
// @Router /error_reports [get]
func (h *DataErrorReportHandler) List(c fiber.Ctx) error {
	_, ok := c.Locals("claims").(*port.JWTClaims)
	if !ok {
		return entity.NewUnauthorized("invalid or missing authentication")
	}

	errorsList, err := h.errorReportService.List(c.Context())
	if err != nil {
		return err
	}

	listResponse := []presenter.DataErrorReportSuccessResponse{}
	for _, report := range errorsList {
		listResponse = append(listResponse, presenter.CreateDataErrorReportSuccessResponse(report.Report))
	}

	return c.Status(fiber.StatusOK).JSON(listResponse)
}

// ChangeStatus godoc
// @Summary Zmień status zgłoszenia błędu danych
// @Description Aktualizuje status wskazanego zgłoszenia błędu danych (oczekujace, rozpatrzone, odrzucone).
// @Tags data-error-reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Identyfikator zgłoszenia błędu danych" Format(uuid)
// @Param request body presenter.ChangeDataErrorReportStatusRequest true "Nowy status zgłoszenia"
// @Success 204 "No Content"
// @Failure 400 {object} presenter.ErrorResponse "Nieprawidłowe body żądania lub identyfikator zgłoszenia"
// @Failure 401 {object} presenter.ErrorResponse "Brak lub nieprawidłowy access token"
// @Failure 500 {object} presenter.ErrorResponse "Nie udało się zmienić statusu zgłoszenia"
// @Router /error_reports/{id}/status [put]
func (h *DataErrorReportHandler) ChangeStatus(c fiber.Ctx) error {
	_, ok := c.Locals("claims").(*port.JWTClaims)
	if !ok {
		return entity.NewUnauthorized("invalid or missing authentication")
	}

	dataErrorReportID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return entity.NewBadRequest("invalid data error report id")
	}

	var req presenter.ChangeDataErrorReportStatusRequest
	if err := c.Bind().Body(&req); err != nil {
		return err
	}

	status := entity.ErrorReportStatus(req.Status)
	err = h.errorReportService.ChangeStatus(c.Context(), dataErrorReportID, status)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// Delete godoc
// @Summary Usuń zgłoszenie błędu danych
// @Description Usuwa wskazane zgłoszenie błędu danych.
// @Tags data-error-reports
// @Security BearerAuth
// @Param id path string true "Identyfikator zgłoszenia błędu danych" Format(uuid)
// @Success 204 "No Content"
// @Failure 400 {object} presenter.ErrorResponse "Nieprawidłowy identyfikator zgłoszenia"
// @Failure 401 {object} presenter.ErrorResponse "Brak lub nieprawidłowy access token"
// @Failure 500 {object} presenter.ErrorResponse "Nie udało się usunąć zgłoszenia"
// @Router /error_reports/{id} [delete]
func (h *DataErrorReportHandler) Delete(c fiber.Ctx) error {
	_, ok := c.Locals("claims").(*port.JWTClaims)
	if !ok {
		return entity.NewUnauthorized("invalid or missing authentication")
	}

	dataErrorReportID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return entity.NewBadRequest("invalid data error report id")
	}

	err = h.errorReportService.Delete(c.Context(), dataErrorReportID)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
