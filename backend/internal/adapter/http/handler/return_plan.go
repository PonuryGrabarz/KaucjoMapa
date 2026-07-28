package handler

import (
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/adapter/http/presenter"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/port"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type ReturnPlanHandler struct {
	returnPlanSvc usecase.IReturnPlanService
}

func NewReturnPlanHandler(svc usecase.IReturnPlanService) *ReturnPlanHandler {
	return &ReturnPlanHandler{returnPlanSvc: svc}
}

func (h *ReturnPlanHandler) claimsUserID(c fiber.Ctx) (uuid.UUID, error) {
	claims, ok := c.Locals("claims").(*port.JWTClaims)
	if !ok {
		return uuid.Nil, entity.NewUnauthorized("invalid or missing authentication")
	}
	return claims.UserID, nil
}

// Create godoc
// @Summary Utwórz draft planu zwrotu
// @Description Tworzy nowy draft planu zwrotu dla zalogowanego użytkownika (max. jeden draft na użytkownika).
// @Tags return-plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body presenter.CreateReturnPlanRequest true "Cele planu"
// @Success 201 {object} presenter.ReturnPlanResponse "Utworzony plan"
// @Failure 400 {object} presenter.ErrorResponse "Nieprawidłowe body żądania"
// @Failure 401 {object} presenter.ErrorResponse "Brak lub nieprawidłowy access token"
// @Failure 409 {object} presenter.ErrorResponse "Draft już istnieje"
// @Failure 500 {object} presenter.ErrorResponse "Nie udało się utworzyć planu"
// @Router /return_plans [post]
func (h *ReturnPlanHandler) Create(c fiber.Ctx) error {
	userID, err := h.claimsUserID(c)
	if err != nil {
		return err
	}

	var req presenter.CreateReturnPlanRequest
	if err := c.Bind().Body(&req); err != nil {
		return err
	}

	plan, err := h.returnPlanSvc.CreateDraft(c.Context(), userID, req.TargetPet, req.TargetCan, req.TargetGlass)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(presenter.ReturnPlanSuccessResponse(plan))
}

// List godoc
// @Summary Lista planów zwrotu
// @Description Zwraca historię planów zwrotu zalogowanego użytkownika.
// @Tags return-plans
// @Produce json
// @Security BearerAuth
// @Success 200 {array} presenter.ReturnPlanSummaryResponse "Lista planów"
// @Failure 401 {object} presenter.ErrorResponse "Brak lub nieprawidłowy access token"
// @Failure 500 {object} presenter.ErrorResponse "Nie udało się pobrać planów"
// @Router /return_plans [get]
func (h *ReturnPlanHandler) List(c fiber.Ctx) error {
	userID, err := h.claimsUserID(c)
	if err != nil {
		return err
	}

	plans, err := h.returnPlanSvc.List(c.Context(), userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(presenter.ListReturnPlanSummariesSuccessResponse(plans))
}

// GetDraft godoc
// @Summary Pobierz draft planu
// @Description Zwraca aktualny draft planu zwrotu zalogowanego użytkownika wraz z punktami.
// @Tags return-plans
// @Produce json
// @Security BearerAuth
// @Success 200 {object} presenter.ReturnPlanResponse "Draft planu"
// @Failure 401 {object} presenter.ErrorResponse "Brak lub nieprawidłowy access token"
// @Failure 404 {object} presenter.ErrorResponse "Brak draftu"
// @Failure 500 {object} presenter.ErrorResponse "Nie udało się pobrać draftu"
// @Router /return_plans/draft [get]
func (h *ReturnPlanHandler) GetDraft(c fiber.Ctx) error {
	userID, err := h.claimsUserID(c)
	if err != nil {
		return err
	}

	plan, err := h.returnPlanSvc.GetDraft(c.Context(), userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(presenter.ReturnPlanSuccessResponse(plan))
}

// GetByID godoc
// @Summary Pobierz plan zwrotu
// @Description Zwraca plan zwrotu po ID wraz z alokacjami punktów.
// @Tags return-plans
// @Produce json
// @Security BearerAuth
// @Param id path string true "Identyfikator planu" Format(uuid)
// @Success 200 {object} presenter.ReturnPlanResponse "Plan zwrotu"
// @Failure 400 {object} presenter.ErrorResponse "Nieprawidłowy identyfikator"
// @Failure 401 {object} presenter.ErrorResponse "Brak lub nieprawidłowy access token"
// @Failure 404 {object} presenter.ErrorResponse "Nie znaleziono planu"
// @Failure 500 {object} presenter.ErrorResponse "Nie udało się pobrać planu"
// @Router /return_plans/{id} [get]
func (h *ReturnPlanHandler) GetByID(c fiber.Ctx) error {
	userID, err := h.claimsUserID(c)
	if err != nil {
		return err
	}

	planID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return entity.NewBadRequest("invalid plan id")
	}

	plan, err := h.returnPlanSvc.GetByID(c.Context(), userID, planID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(presenter.ReturnPlanSuccessResponse(plan))
}

// UpdateTargets godoc
// @Summary Zaktualizuj cele planu
// @Description Zmienia docelowe ilości opakowań w planie.
// @Tags return-plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Identyfikator planu" Format(uuid)
// @Param request body presenter.UpdateReturnPlanTargetsRequest true "Nowe cele"
// @Success 200 {object} presenter.ReturnPlanResponse "Zaktualizowany plan"
// @Failure 400 {object} presenter.ErrorResponse "Nieprawidłowe body lub cele poniżej alokacji"
// @Failure 401 {object} presenter.ErrorResponse "Brak lub nieprawidłowy access token"
// @Failure 404 {object} presenter.ErrorResponse "Nie znaleziono planu"
// @Failure 500 {object} presenter.ErrorResponse "Nie udało się zaktualizować planu"
// @Router /return_plans/{id}/targets [put]
func (h *ReturnPlanHandler) UpdateTargets(c fiber.Ctx) error {
	userID, err := h.claimsUserID(c)
	if err != nil {
		return err
	}

	planID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return entity.NewBadRequest("invalid plan id")
	}

	var req presenter.UpdateReturnPlanTargetsRequest
	if err := c.Bind().Body(&req); err != nil {
		return err
	}

	plan, err := h.returnPlanSvc.UpdateTargets(c.Context(), userID, planID, req.TargetPet, req.TargetCan, req.TargetGlass)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(presenter.ReturnPlanSuccessResponse(plan))
}

// Save godoc
// @Summary Zapisz plan
// @Description Zmienia status draftu na saved.
// @Tags return-plans
// @Produce json
// @Security BearerAuth
// @Param id path string true "Identyfikator planu" Format(uuid)
// @Success 200 {object} presenter.ReturnPlanResponse "Zapisany plan"
// @Failure 400 {object} presenter.ErrorResponse "Plan nie może zostać zapisany"
// @Failure 401 {object} presenter.ErrorResponse "Brak lub nieprawidłowy access token"
// @Failure 404 {object} presenter.ErrorResponse "Nie znaleziono planu"
// @Failure 500 {object} presenter.ErrorResponse "Nie udało się zapisać planu"
// @Router /return_plans/{id}/save [post]
func (h *ReturnPlanHandler) Save(c fiber.Ctx) error {
	userID, err := h.claimsUserID(c)
	if err != nil {
		return err
	}

	planID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return entity.NewBadRequest("invalid plan id")
	}

	plan, err := h.returnPlanSvc.Save(c.Context(), userID, planID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(presenter.ReturnPlanSuccessResponse(plan))
}

// ConfirmAll godoc
// @Summary Potwierdź cały zwrot
// @Description Potwierdza wszystkie pending punkty i zamyka plan, gdy alokacja jest kompletna.
// @Tags return-plans
// @Produce json
// @Security BearerAuth
// @Param id path string true "Identyfikator planu" Format(uuid)
// @Success 200 {object} presenter.ReturnPlanResponse "Potwierdzony plan"
// @Failure 400 {object} presenter.ErrorResponse "Plan niekompletny lub nieedytowalny"
// @Failure 401 {object} presenter.ErrorResponse "Brak lub nieprawidłowy access token"
// @Failure 404 {object} presenter.ErrorResponse "Nie znaleziono planu"
// @Failure 500 {object} presenter.ErrorResponse "Nie udało się potwierdzić planu"
// @Router /return_plans/{id}/confirm [post]
func (h *ReturnPlanHandler) ConfirmAll(c fiber.Ctx) error {
	userID, err := h.claimsUserID(c)
	if err != nil {
		return err
	}

	planID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return entity.NewBadRequest("invalid plan id")
	}

	plan, err := h.returnPlanSvc.ConfirmAll(c.Context(), userID, planID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(presenter.ReturnPlanSuccessResponse(plan))
}

// Delete godoc
// @Summary Usuń plan zwrotu
// @Description Usuwa niepotwierdzony plan należący do zalogowanego użytkownika.
// @Tags return-plans
// @Security BearerAuth
// @Param id path string true "Identyfikator planu" Format(uuid)
// @Success 204 "No Content"
// @Failure 400 {object} presenter.ErrorResponse "Nieprawidłowy identyfikator lub plan potwierdzony"
// @Failure 401 {object} presenter.ErrorResponse "Brak lub nieprawidłowy access token"
// @Failure 404 {object} presenter.ErrorResponse "Nie znaleziono planu"
// @Failure 500 {object} presenter.ErrorResponse "Nie udało się usunąć planu"
// @Router /return_plans/{id} [delete]
func (h *ReturnPlanHandler) Delete(c fiber.Ctx) error {
	userID, err := h.claimsUserID(c)
	if err != nil {
		return err
	}

	planID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return entity.NewBadRequest("invalid plan id")
	}

	if err := h.returnPlanSvc.Delete(c.Context(), userID, planID); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// AddPoint godoc
// @Summary Dodaj punkt do planu
// @Description Dodaje punkt zwrotu do planu i opcjonalnie ustawia początkowe ilości.
// @Tags return-plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Identyfikator planu" Format(uuid)
// @Param request body presenter.AddReturnPlanPointRequest true "Punkt i ilości"
// @Success 201 {object} presenter.ReturnPlanPointResponse "Dodany punkt planu"
// @Failure 400 {object} presenter.ErrorResponse "Nieprawidłowe body lub punkt nie przyjmuje typu"
// @Failure 401 {object} presenter.ErrorResponse "Brak lub nieprawidłowy access token"
// @Failure 404 {object} presenter.ErrorResponse "Nie znaleziono planu lub punktu zwrotu"
// @Failure 409 {object} presenter.ErrorResponse "Punkt już jest w planie"
// @Failure 500 {object} presenter.ErrorResponse "Nie udało się dodać punktu"
// @Router /return_plans/{id}/points [post]
func (h *ReturnPlanHandler) AddPoint(c fiber.Ctx) error {
	userID, err := h.claimsUserID(c)
	if err != nil {
		return err
	}

	planID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return entity.NewBadRequest("invalid plan id")
	}

	var req presenter.AddReturnPlanPointRequest
	if err := c.Bind().Body(&req); err != nil {
		return err
	}

	returnPointID, err := uuid.Parse(req.ReturnPointID)
	if err != nil {
		return entity.NewBadRequest("invalid return point id")
	}

	point, err := h.returnPlanSvc.AddPoint(c.Context(), userID, planID, returnPointID, req.PetCount, req.CanCount, req.GlassCount)
	if err != nil {
		return err
	}

	plan, err := h.returnPlanSvc.GetByID(c.Context(), userID, planID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(presenter.ReturnPlanPointSuccessResponse(point, plan.Plan.RatePet, plan.Plan.RateCan, plan.Plan.RateGlass))
}

// UpdatePointQuantities godoc
// @Summary Zaktualizuj ilości w punkcie planu
// @Description Zmienia ilości opakowań przypisane do pending punktu w planie.
// @Tags return-plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Identyfikator planu" Format(uuid)
// @Param point_id path string true "Identyfikator punktu planu" Format(uuid)
// @Param request body presenter.UpdateReturnPlanPointQuantitiesRequest true "Nowe ilości"
// @Success 200 {object} presenter.ReturnPlanPointResponse "Zaktualizowany punkt planu"
// @Failure 400 {object} presenter.ErrorResponse "Nieprawidłowe body lub przekroczony cel"
// @Failure 401 {object} presenter.ErrorResponse "Brak lub nieprawidłowy access token"
// @Failure 404 {object} presenter.ErrorResponse "Nie znaleziono planu lub punktu"
// @Failure 500 {object} presenter.ErrorResponse "Nie udało się zaktualizować punktu"
// @Router /return_plans/{id}/points/{point_id} [put]
func (h *ReturnPlanHandler) UpdatePointQuantities(c fiber.Ctx) error {
	userID, err := h.claimsUserID(c)
	if err != nil {
		return err
	}

	planID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return entity.NewBadRequest("invalid plan id")
	}

	planPointID, err := uuid.Parse(c.Params("point_id"))
	if err != nil {
		return entity.NewBadRequest("invalid plan point id")
	}

	var req presenter.UpdateReturnPlanPointQuantitiesRequest
	if err := c.Bind().Body(&req); err != nil {
		return err
	}

	point, err := h.returnPlanSvc.UpdatePointQuantities(c.Context(), userID, planID, planPointID, req.PetCount, req.CanCount, req.GlassCount)
	if err != nil {
		return err
	}

	plan, err := h.returnPlanSvc.GetByID(c.Context(), userID, planID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(presenter.ReturnPlanPointSuccessResponse(point, plan.Plan.RatePet, plan.Plan.RateCan, plan.Plan.RateGlass))
}

// ConfirmPoint godoc
// @Summary Potwierdź oddanie w punkcie
// @Description Potwierdza pending punkt planu i aktualizuje status planu.
// @Tags return-plans
// @Produce json
// @Security BearerAuth
// @Param id path string true "Identyfikator planu" Format(uuid)
// @Param point_id path string true "Identyfikator punktu planu" Format(uuid)
// @Success 200 {object} presenter.ReturnPlanResponse "Zaktualizowany plan"
// @Failure 400 {object} presenter.ErrorResponse "Punkt nie może zostać potwierdzony"
// @Failure 401 {object} presenter.ErrorResponse "Brak lub nieprawidłowy access token"
// @Failure 404 {object} presenter.ErrorResponse "Nie znaleziono planu lub punktu"
// @Failure 500 {object} presenter.ErrorResponse "Nie udało się potwierdzić punktu"
// @Router /return_plans/{id}/points/{point_id}/confirm [post]
func (h *ReturnPlanHandler) ConfirmPoint(c fiber.Ctx) error {
	userID, err := h.claimsUserID(c)
	if err != nil {
		return err
	}

	planID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return entity.NewBadRequest("invalid plan id")
	}

	planPointID, err := uuid.Parse(c.Params("point_id"))
	if err != nil {
		return entity.NewBadRequest("invalid plan point id")
	}

	plan, err := h.returnPlanSvc.ConfirmPoint(c.Context(), userID, planID, planPointID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(presenter.ReturnPlanSuccessResponse(plan))
}

// DeletePoint godoc
// @Summary Usuń punkt z planu
// @Description Usuwa pending punkt z planu zwrotu.
// @Tags return-plans
// @Security BearerAuth
// @Param id path string true "Identyfikator planu" Format(uuid)
// @Param point_id path string true "Identyfikator punktu planu" Format(uuid)
// @Success 204 "No Content"
// @Failure 400 {object} presenter.ErrorResponse "Punkt potwierdzony lub nieprawidłowy identyfikator"
// @Failure 401 {object} presenter.ErrorResponse "Brak lub nieprawidłowy access token"
// @Failure 404 {object} presenter.ErrorResponse "Nie znaleziono planu lub punktu"
// @Failure 500 {object} presenter.ErrorResponse "Nie udało się usunąć punktu"
// @Router /return_plans/{id}/points/{point_id} [delete]
func (h *ReturnPlanHandler) DeletePoint(c fiber.Ctx) error {
	userID, err := h.claimsUserID(c)
	if err != nil {
		return err
	}

	planID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return entity.NewBadRequest("invalid plan id")
	}

	planPointID, err := uuid.Parse(c.Params("point_id"))
	if err != nil {
		return entity.NewBadRequest("invalid plan point id")
	}

	if err := h.returnPlanSvc.DeletePoint(c.Context(), userID, planID, planPointID); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
