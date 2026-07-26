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

type NotificationsHandler struct {
	notificationsSvc usecase.INotificationsService
}

func NewNotificationsHandler(svc usecase.INotificationsService) *NotificationsHandler {
	return &NotificationsHandler{
		notificationsSvc: svc,
	}
}

// Create godoc
// @Summary Utwórz powiadomienie
// @Description Tworzy powiadomienie przypisane do zalogowanego użytkownika.
// @Tags notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body presenter.CreateNotificationRequest true "Dane powiadomienia"
// @Success 201 {object} presenter.NotificationResponse "Utworzone powiadomienie"
// @Failure 400 {object} presenter.ErrorResponse "Nieprawidłowe body żądania"
// @Failure 401 {object} presenter.ErrorResponse "Brak lub nieprawidłowy access token"
// @Failure 500 {object} presenter.ErrorResponse "Nie udało się utworzyć powiadomienia"
// @Router /notifications [post]
func (h *NotificationsHandler) Create(c fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*port.JWTClaims)
	if !ok {
		return entity.NewUnauthorized("invalid or missing authentication")
	}

	var req presenter.CreateNotificationRequest
	if err := c.Bind().Body(&req); err != nil {
		return err
	}

	newNotification, err := h.notificationsSvc.Create(c.Context(), entity.Notification{
		ID:          uuid.New(),
		UserID:      &claims.UserID,
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now(),
	})

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(presenter.CreateNotificationSuccessReponse(newNotification))
}

// List godoc
// @Summary Pobierz powiadomienia
// @Description Zwraca listę powiadomień zalogowanego użytkownika.
// @Tags notifications
// @Produce json
// @Security BearerAuth
// @Success 200 {array} presenter.NotificationResponse "Lista powiadomień"
// @Failure 401 {object} presenter.ErrorResponse "Brak lub nieprawidłowy access token"
// @Failure 500 {object} presenter.ErrorResponse "Nie udało się pobrać powiadomień"
// @Router /notifications [get]
func (h *NotificationsHandler) List(c fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*port.JWTClaims)
	if !ok {
		return entity.NewUnauthorized("invalid or missing authentication")
	}

	notifications, err := h.notificationsSvc.List(c.Context(), claims.UserID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(presenter.ListNotificationsSuccessReponse(notifications))
}

// MarkAsSeen godoc
// @Summary Oznacz powiadomienie jako przeczytane
// @Description Oznacza wskazane powiadomienie zalogowanego użytkownika jako przeczytane.
// @Tags notifications
// @Produce json
// @Security BearerAuth
// @Param id path string true "Identyfikator powiadomienia" Format(uuid)
// @Success 200 {object} presenter.NotificationResponse "Zaktualizowane powiadomienie"
// @Failure 400 {object} presenter.ErrorResponse "Nieprawidłowy identyfikator powiadomienia"
// @Failure 401 {object} presenter.ErrorResponse "Brak lub nieprawidłowy access token"
// @Failure 403 {object} presenter.ErrorResponse "Powiadomienie nie należy do zalogowanego użytkownika"
// @Failure 500 {object} presenter.ErrorResponse "Nie udało się zaktualizować powiadomienia"
// @Router /notifications/{id} [put]
func (h *NotificationsHandler) MarkAsSeen(c fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*port.JWTClaims)
	if !ok {
		return entity.NewUnauthorized("invalid or missing authentication")
	}

	ID := c.Params("id")
	notificationID, err := uuid.Parse(ID)
	if err != nil {
		return entity.NewBadRequest("invalid notification id")
	}

	notification, err := h.notificationsSvc.GetByID(c.Context(), notificationID)
	if err != nil {
		return err
	}

	if *notification.UserID != claims.UserID {
		return entity.NewForbidden("invalid user id")
	}

	seenNotification := entity.Notification{
		ID:          notification.ID,
		UserID:      &claims.UserID,
		Name:        notification.Name,
		Description: notification.Description,
		Seen:        true,
		CreatedAt:   notification.CreatedAt,
	}

	updatedNotification, err := h.notificationsSvc.Update(c.Context(), notification.ID, seenNotification)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(presenter.CreateNotificationSuccessReponse(updatedNotification))

}

// Delete godoc
// @Summary Usuń powiadomienie
// @Description Usuwa wskazane powiadomienie należące do zalogowanego użytkownika.
// @Tags notifications
// @Security BearerAuth
// @Param id path string true "Identyfikator powiadomienia" Format(uuid)
// @Success 204 "No Content"
// @Failure 400 {object} presenter.ErrorResponse "Nieprawidłowy identyfikator powiadomienia"
// @Failure 401 {object} presenter.ErrorResponse "Brak lub nieprawidłowy access token"
// @Failure 403 {object} presenter.ErrorResponse "Powiadomienie nie należy do zalogowanego użytkownika"
// @Failure 500 {object} presenter.ErrorResponse "Nie udało się usunąć powiadomienia"
// @Router /notifications/{id} [delete]
func (h *NotificationsHandler) Delete(c fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*port.JWTClaims)
	if !ok {
		return entity.NewUnauthorized("invalid or missing authentication")
	}

	ID := c.Params("id")
	notificationID, err := uuid.Parse(ID)
	if err != nil {
		return entity.NewBadRequest("invalid notification id")
	}

	notification, err := h.notificationsSvc.GetByID(c.Context(), notificationID)
	if err != nil {
		return err
	}

	if *notification.UserID != claims.UserID {
		return entity.NewForbidden("invalid user id")
	}

	err = h.notificationsSvc.Delete(c.Context(), notificationID)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
