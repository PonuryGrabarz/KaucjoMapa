package handler

import (
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/adapter/http/presenter"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/port"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type CommunityReportHandler struct {
	communityReportSvc usecase.ICommunityReportsService
	userSvc            usecase.IUserService
}

func NewCommunityReportHandler(svc usecase.ICommunityReportsService, userSvc usecase.IUserService) *CommunityReportHandler {
	return &CommunityReportHandler{
		communityReportSvc: svc,
		userSvc:            userSvc,
	}
}

// Create godoc
// @Summary Utwórz zgłoszenie społecznościowe
// @Description Tworzy zgłoszenie dotyczące statusu punktu zwrotu, przypisane do zalogowanego użytkownika.
// @Tags community-reports
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body presenter.CreateCommunityReportRequest true "Dane zgłoszenia społecznościowego"
// @Success 201 {object} presenter.CreateCommunityReportResponse "Utworzone zgłoszenie społecznościowe"
// @Failure 400 {object} presenter.ErrorResponse "Nieprawidłowe body żądania lub identyfikator punktu zwrotu"
// @Failure 401 {object} presenter.ErrorResponse "Brak lub nieprawidłowy access token"
// @Failure 500 {object} presenter.ErrorResponse "Nie udało się utworzyć zgłoszenia"
// @Router /reports [post]
func (h *CommunityReportHandler) Create(c fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*port.JWTClaims)
	if !ok {
		return entity.NewUnauthorized("invalid or missing authentication")
	}

	var req presenter.CreateCommunityReportRequest
	if err := c.Bind().Body(&req); err != nil {
		return err
	}

	rID, err := uuid.Parse(req.ReturnPointID)
	if err != nil {
		return entity.NewBadRequest("invalid identifier")
	}

	newReport, err := h.communityReportSvc.Create(c.Context(), entity.CommunityReport{
		ID:            uuid.New(),
		UserID:        &claims.UserID,
		ReturnPointID: &rID,
		Status:        entity.PointStatus(req.Status),
		Description:   req.Description,
	})
	if err != nil {
		return err
	}

	user, err := h.userSvc.GetUserByID(c.Context(), claims.UserID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(presenter.CommunityReportSuccessResponse(newReport, user))
}

// List godoc
// @Summary Pobierz zgłoszenia społecznościowe
// @Description Zwraca zgłoszenia społecznościowe dotyczące wskazanego punktu zwrotu.
// @Tags community-reports
// @Accept json
// @Produce json
// @Param return_point_id path string true "Identyfikator punktu zwrotu" Format(uuid)
// @Success 200 {array} presenter.CreateCommunityReportResponse "Lista zgłoszeń społecznościowych"
// @Failure 400 {object} presenter.ErrorResponse "Nieprawidłowe body żądania lub identyfikator punktu zwrotu"
// @Failure 500 {object} presenter.ErrorResponse "Nie udało się pobrać zgłoszeń"
// @Router /reports/{return_point_id} [get]
func (h *CommunityReportHandler) List(c fiber.Ctx) error {
	id := c.Params("return_point_id")

	rID, err := uuid.Parse(id)
	if err != nil {
		return entity.NewBadRequest("invalid identifier")
	}

	reports, err := h.communityReportSvc.List(c.Context(), rID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(presenter.ListCommunityReportsSuccessReponse(reports))
}

// Delete godoc
// @Summary Usuń zgłoszenie społecznościowe
// @Description Usuwa wskazane zgłoszenie społecznościowe należące do zalogowanego użytkownika.
// @Tags community-reports
// @Security BearerAuth
// @Param id path string true "Identyfikator zgłoszenia społecznościowego" Format(uuid)
// @Success 204 "No Content"
// @Failure 400 {object} presenter.ErrorResponse "Nieprawidłowy identyfikator zgłoszenia"
// @Failure 401 {object} presenter.ErrorResponse "Brak lub nieprawidłowy access token"
// @Failure 500 {object} presenter.ErrorResponse "Nie udało się usunąć zgłoszenia"
// @Router /reports/{id} [delete]
func (h *CommunityReportHandler) Delete(c fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*port.JWTClaims)
	if !ok {
		return entity.NewUnauthorized("invalid or missing authentication")
	}

	id := c.Params("id")

	reportID, err := uuid.Parse(id)
	if err != nil {
		return entity.NewBadRequest("invalid identifier")
	}

	err = h.communityReportSvc.Delete(c.Context(), claims.UserID, reportID)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *CommunityReportHandler) GetByID(c fiber.Ctx) error {
	id := c.Params("id")

	reportID, err := uuid.Parse(id)
	if err != nil {
		return entity.NewBadRequest("invalid identifier")
	}

	report, err := h.communityReportSvc.GetByID(c.Context(), reportID)
	if err != nil {
		return err
	}
	user, err := h.userSvc.GetUserByID(c.Context(), *report.UserID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(presenter.GetCommunityReportSuccessResponse(report, user))
}
