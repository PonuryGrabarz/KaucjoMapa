package handler

import (
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/adapter/http/presenter"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/port"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type BookmarkHandler struct {
	bookmarkSvc usecase.IBookmarkService
}

func NewBookmarkHandler(svc usecase.IBookmarkService) *BookmarkHandler {
	return &BookmarkHandler{bookmarkSvc: svc}
}

// Create godoc
// @Summary Dodaj zakładkę
// @Description Tworzy zakładkę punktu zwrotu dla zalogowanego użytkownika.
// @Tags bookmarks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body presenter.CreateBookmarkRequest true "Dane zakładki"
// @Success 201 {object} presenter.BookmarkResponse "Utworzona zakładka"
// @Failure 400 {object} presenter.ErrorResponse "Nieprawidłowe body żądania"
// @Failure 401 {object} presenter.ErrorResponse "Brak lub nieprawidłowy access token"
// @Failure 409 {object} presenter.ErrorResponse "Zakładka już istnieje"
// @Failure 500 {object} presenter.ErrorResponse "Nie udało się utworzyć zakładki"
// @Router /bookmarks [post]
func (h *BookmarkHandler) Create(c fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*port.JWTClaims)
	if !ok {
		return entity.NewUnauthorized("invalid or missing authentication")
	}

	var req presenter.CreateBookmarkRequest
	if err := c.Bind().Body(&req); err != nil {
		return err
	}

	returnPointID, err := uuid.Parse(req.ReturnPointID)
	if err != nil {
		return entity.NewBadRequest("invalid return point id")
	}

	bookmark, err := h.bookmarkSvc.Create(c.Context(), claims.UserID, returnPointID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(presenter.BookmarkSuccessResponse(bookmark))
}

// List godoc
// @Summary Lista zakładek
// @Description Zwraca zakładki zalogowanego użytkownika wraz z danymi punktów zwrotu.
// @Tags bookmarks
// @Produce json
// @Security BearerAuth
// @Success 200 {array} presenter.BookmarkWithReturnPointResponse "Lista zakładek"
// @Failure 401 {object} presenter.ErrorResponse "Brak lub nieprawidłowy access token"
// @Failure 500 {object} presenter.ErrorResponse "Nie udało się pobrać zakładek"
// @Router /bookmarks [get]
func (h *BookmarkHandler) List(c fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*port.JWTClaims)
	if !ok {
		return entity.NewUnauthorized("invalid or missing authentication")
	}

	bookmarks, err := h.bookmarkSvc.ListWithReturnPointsByUserID(c.Context(), claims.UserID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(presenter.ListBookmarksWithReturnPointsSuccessResponse(bookmarks))
}

// GetByID godoc
// @Summary Pobierz zakładkę
// @Description Zwraca zakładkę po identyfikatorze, jeśli należy do zalogowanego użytkownika.
// @Tags bookmarks
// @Produce json
// @Security BearerAuth
// @Param id path string true "Identyfikator zakładki" Format(uuid)
// @Success 200 {object} presenter.BookmarkResponse "Zakładka"
// @Failure 400 {object} presenter.ErrorResponse "Nieprawidłowy identyfikator zakładki"
// @Failure 401 {object} presenter.ErrorResponse "Brak lub nieprawidłowy access token"
// @Failure 403 {object} presenter.ErrorResponse "Zakładka nie należy do zalogowanego użytkownika"
// @Failure 404 {object} presenter.ErrorResponse "Nie znaleziono zakładki"
// @Failure 500 {object} presenter.ErrorResponse "Nie udało się pobrać zakładki"
// @Router /bookmarks/{id} [get]
func (h *BookmarkHandler) GetByID(c fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*port.JWTClaims)
	if !ok {
		return entity.NewUnauthorized("invalid or missing authentication")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return entity.NewBadRequest("invalid bookmark id")
	}

	bookmark, err := h.bookmarkSvc.GetByID(c.Context(), id)
	if err != nil {
		return err
	}

	if bookmark.UserID != claims.UserID {
		return entity.NewForbidden("bookmark does not belong to the authenticated user")
	}

	return c.Status(fiber.StatusOK).JSON(presenter.BookmarkSuccessResponse(bookmark))
}

// GetByReturnPointID godoc
// @Summary Pobierz zakładkę po punkcie zwrotu
// @Description Zwraca zakładkę zalogowanego użytkownika dla wskazanego punktu zwrotu.
// @Tags bookmarks
// @Produce json
// @Security BearerAuth
// @Param return_point_id path string true "Identyfikator punktu zwrotu" Format(uuid)
// @Success 200 {object} presenter.BookmarkResponse "Zakładka"
// @Failure 400 {object} presenter.ErrorResponse "Nieprawidłowy identyfikator punktu zwrotu"
// @Failure 401 {object} presenter.ErrorResponse "Brak lub nieprawidłowy access token"
// @Failure 404 {object} presenter.ErrorResponse "Nie znaleziono zakładki"
// @Failure 500 {object} presenter.ErrorResponse "Nie udało się pobrać zakładki"
// @Router /bookmarks/return_point/{return_point_id} [get]
func (h *BookmarkHandler) GetByReturnPointID(c fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*port.JWTClaims)
	if !ok {
		return entity.NewUnauthorized("invalid or missing authentication")
	}

	returnPointID, err := uuid.Parse(c.Params("return_point_id"))
	if err != nil {
		return entity.NewBadRequest("invalid return point id")
	}

	bookmark, err := h.bookmarkSvc.GetByUserAndReturnPoint(c.Context(), claims.UserID, returnPointID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(presenter.BookmarkSuccessResponse(bookmark))
}

// Update godoc
// @Summary Zaktualizuj zakładkę
// @Description Zmienia punkt zwrotu przypisany do zakładki zalogowanego użytkownika.
// @Tags bookmarks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Identyfikator zakładki" Format(uuid)
// @Param request body presenter.UpdateBookmarkRequest true "Nowe dane zakładki"
// @Success 200 {object} presenter.BookmarkResponse "Zaktualizowana zakładka"
// @Failure 400 {object} presenter.ErrorResponse "Nieprawidłowe body lub identyfikator"
// @Failure 401 {object} presenter.ErrorResponse "Brak lub nieprawidłowy access token"
// @Failure 404 {object} presenter.ErrorResponse "Nie znaleziono zakładki"
// @Failure 409 {object} presenter.ErrorResponse "Konflikt unikalności zakładki"
// @Failure 500 {object} presenter.ErrorResponse "Nie udało się zaktualizować zakładki"
// @Router /bookmarks/{id} [put]
func (h *BookmarkHandler) Update(c fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*port.JWTClaims)
	if !ok {
		return entity.NewUnauthorized("invalid or missing authentication")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return entity.NewBadRequest("invalid bookmark id")
	}

	var req presenter.UpdateBookmarkRequest
	if err := c.Bind().Body(&req); err != nil {
		return err
	}

	returnPointID, err := uuid.Parse(req.ReturnPointID)
	if err != nil {
		return entity.NewBadRequest("invalid return point id")
	}

	bookmark, err := h.bookmarkSvc.Update(c.Context(), claims.UserID, id, returnPointID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(presenter.BookmarkSuccessResponse(bookmark))
}

// Delete godoc
// @Summary Usuń zakładkę
// @Description Usuwa zakładkę należącą do zalogowanego użytkownika.
// @Tags bookmarks
// @Security BearerAuth
// @Param id path string true "Identyfikator zakładki" Format(uuid)
// @Success 204 "No Content"
// @Failure 400 {object} presenter.ErrorResponse "Nieprawidłowy identyfikator zakładki"
// @Failure 401 {object} presenter.ErrorResponse "Brak lub nieprawidłowy access token"
// @Failure 500 {object} presenter.ErrorResponse "Nie udało się usunąć zakładki"
// @Router /bookmarks/{id} [delete]
func (h *BookmarkHandler) Delete(c fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*port.JWTClaims)
	if !ok {
		return entity.NewUnauthorized("invalid or missing authentication")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return entity.NewBadRequest("invalid bookmark id")
	}

	if err := h.bookmarkSvc.DeleteByID(c.Context(), claims.UserID, id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// DeleteByReturnPointID godoc
// @Summary Usuń zakładkę po punkcie zwrotu
// @Description Usuwa zakładkę zalogowanego użytkownika dla wskazanego punktu zwrotu.
// @Tags bookmarks
// @Security BearerAuth
// @Param return_point_id path string true "Identyfikator punktu zwrotu" Format(uuid)
// @Success 204 "No Content"
// @Failure 400 {object} presenter.ErrorResponse "Nieprawidłowy identyfikator punktu zwrotu"
// @Failure 401 {object} presenter.ErrorResponse "Brak lub nieprawidłowy access token"
// @Failure 500 {object} presenter.ErrorResponse "Nie udało się usunąć zakładki"
// @Router /bookmarks/return_point/{return_point_id} [delete]
func (h *BookmarkHandler) DeleteByReturnPointID(c fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*port.JWTClaims)
	if !ok {
		return entity.NewUnauthorized("invalid or missing authentication")
	}

	returnPointID, err := uuid.Parse(c.Params("return_point_id"))
	if err != nil {
		return entity.NewBadRequest("invalid return point id")
	}

	if err := h.bookmarkSvc.DeleteByUserAndReturnPoint(c.Context(), claims.UserID, returnPointID); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
