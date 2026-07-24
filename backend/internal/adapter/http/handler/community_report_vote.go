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

type CommunityReportVoteHandler struct {
	communityReportVoteService usecase.ICommunityReportVoteService
}

func NewCommunityReportVoteHandler(communityReportVoteService usecase.ICommunityReportVoteService) *CommunityReportVoteHandler {
	return &CommunityReportVoteHandler{
		communityReportVoteService: communityReportVoteService,
	}
}

// Create godoc
// @Summary Oddaj głos na zgłoszenie społecznościowe
// @Description Tworzy lub aktualizuje głos zalogowanego użytkownika na wskazane zgłoszenie społecznościowe.
// @Tags community-report-votes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param report_id path string true "Identyfikator zgłoszenia społecznościowego" Format(uuid)
// @Param request body presenter.CreateCommunityReportVoteRequest true "Głos: confirm albo reject"
// @Success 201 {object} entity.CommunityReportVote "Utworzony lub zaktualizowany głos"
// @Failure 400 {object} presenter.ErrorResponse "Nieprawidłowy identyfikator zgłoszenia lub body żądania"
// @Failure 401 {object} presenter.ErrorResponse "Brak lub nieprawidłowy access token"
// @Failure 500 {object} presenter.ErrorResponse "Nie udało się zapisać głosu"
// @Router /votes/{report_id} [post]
func (h *CommunityReportVoteHandler) Create(c fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*port.JWTClaims)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid JWT claims")
	}

	reportIDParam := c.Params("report_id")
	reportID, err := uuid.Parse(reportIDParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid report ID")
	}

	var req presenter.CreateCommunityReportVoteRequest
	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	newVote, err := h.communityReportVoteService.Create(c.Context(), entity.CommunityReportVote{
		CommunityReportID: &reportID,
		UserID:            &claims.UserID,
		Vote:              entity.CommunityReportVoteValue(req.Vote),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create community report vote")
	}

	return c.Status(fiber.StatusCreated).JSON(
		newVote,
	)
}
