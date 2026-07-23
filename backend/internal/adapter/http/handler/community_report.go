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
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid JWT claims")
	}

	var req presenter.CreateCommunityReportRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	rID, err := uuid.Parse(req.ReturnPointID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid report ID")
	}

	newReport, err := h.communityReportSvc.Create(c.Context(), entity.CommunityReport{
		ID:            uuid.New(),
		UserID:        &claims.UserID,
		ReturnPointID: &rID,
		Status:        entity.PointStatus(req.Status),
		Description:   req.Description,
	})
	if err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, "Failed creating report")
	}

	user, err := h.userSvc.GetUserByID(c.Context(), claims.UserID)
	if err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, "Failed creating report")
	}

	return c.Status(fiber.StatusCreated).JSON(presenter.CommunityReportSuccessResponse(newReport, user))
}

func (h *CommunityReportHandler) List(c fiber.Ctx) error {
	var req presenter.ListReportsByReturnPointIDRequest

	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	rID, err := uuid.Parse(req.ReturnPointID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid report ID")
	}

	reports, err := h.communityReportSvc.List(c.Context(), rID)
	if err != nil {
		return fiber.NewError(fiber.ErrInternalServerError.Code, "Failed to list reports")
	}
}
