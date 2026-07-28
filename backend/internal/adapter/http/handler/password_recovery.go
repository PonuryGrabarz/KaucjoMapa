package handler

import (
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/adapter/http/presenter"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/dto"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type PasswordRecoveryHandler struct {
	recoverySvc usecase.IPasswordRecoveryService
}

func NewPasswordRecoveryHandler(recoverySvc usecase.IPasswordRecoveryService) *PasswordRecoveryHandler {
	return &PasswordRecoveryHandler{
		recoverySvc: recoverySvc,
	}
}

func challengeIDFromPath(c fiber.Ctx) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Params("challenge_id"))
	if err != nil {
		return uuid.Nil, entity.NewBadRequest("invalid challenge_id")
	}

	return id, nil
}

// StartRecovery godoc
// @Summary Rozpocznij odzyskiwanie hasła
// @Description Tworzy wyzwanie resetowania hasła i wysyła 8-cyfrowy kod na podany e-mail. Odpowiedź jest zawsze taka sama, niezależnie od tego, czy e-mail istnieje w systemie.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body presenter.StartRecoveryRequest true "Adres e-mail konta"
// @Success 202 {object} presenter.RecoveryChallengeEnvelope "Wyzwanie utworzone, kod wysłany"
// @Failure 400 {object} presenter.ErrorResponse "Nieprawidłowe body żądania"
// @Failure 500 {object} presenter.ErrorResponse "Wewnętrzny błąd serwera"
// @Router /auth/password-recovery [post]
func (h *PasswordRecoveryHandler) StartRecovery(c fiber.Ctx) error {
	var req presenter.StartRecoveryRequest

	if err := c.Bind().Body(&req); err != nil {
		return err
	}

	challenge, err := h.recoverySvc.StartRecovery(c.Context(), dto.StartRecoveryInput{
		Email: req.Email,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusAccepted).JSON(presenter.RecoveryChallengeSuccessResponse(challenge))
}

// VerifyCode godoc
// @Summary Zweryfikuj kod resetowania hasła
// @Description Sprawdza 8-cyfrowy kod i zwraca jednorazowy reset token ważny 10 minut. Maksymalnie 5 prób na wyzwanie.
// @Tags auth
// @Accept json
// @Produce json
// @Param challenge_id path string true "ID wyzwania (UUID)"
// @Param request body presenter.VerifyCodeRequest true "Kod z wiadomości e-mail"
// @Success 200 {object} presenter.ResetTokenEnvelope "Reset token"
// @Failure 400 {object} presenter.ErrorResponse "Nieprawidłowe body żądania"
// @Failure 401 {object} presenter.ErrorResponse "Nieprawidłowy kod"
// @Failure 403 {object} presenter.ErrorResponse "Przekroczono limit prób"
// @Failure 404 {object} presenter.ErrorResponse "Wyzwanie nie istnieje lub wygasło"
// @Failure 500 {object} presenter.ErrorResponse "Wewnętrzny błąd serwera"
// @Router /auth/password-recovery/{challenge_id}/verify [post]
func (h *PasswordRecoveryHandler) VerifyCode(c fiber.Ctx) error {
	challengeID, err := challengeIDFromPath(c)
	if err != nil {
		return err
	}

	var req presenter.VerifyCodeRequest
	if err := c.Bind().Body(&req); err != nil {
		return err
	}

	resetToken, expiresAt, err := h.recoverySvc.VerifyCode(c.Context(), dto.VerifyCodeInput{
		ChallengeID: challengeID,
		Code:        req.Code,
	})
	if err != nil {
		return err
	}

	return c.JSON(presenter.ResetTokenSuccessResponse(resetToken, expiresAt))
}

// ResendCode godoc
// @Summary Wyślij kod resetowania ponownie
// @Description Generuje nowy kod (stary traci ważność), zeruje licznik prób i przedłuża wyzwanie o 10 minut. Dostępne najwcześniej 60 sekund po poprzedniej wysyłce.
// @Tags auth
// @Produce json
// @Param challenge_id path string true "ID wyzwania (UUID)"
// @Success 202 {object} presenter.RecoveryChallengeEnvelope "Nowy kod wysłany"
// @Failure 400 {object} presenter.ErrorResponse "Nieprawidłowe ID wyzwania"
// @Failure 404 {object} presenter.ErrorResponse "Wyzwanie nie istnieje lub zostało zużyte"
// @Failure 409 {object} presenter.ErrorResponse "Ponowna wysyłka jeszcze niedostępna"
// @Failure 500 {object} presenter.ErrorResponse "Wewnętrzny błąd serwera"
// @Router /auth/password-recovery/{challenge_id}/resend [post]
func (h *PasswordRecoveryHandler) ResendCode(c fiber.Ctx) error {
	challengeID, err := challengeIDFromPath(c)
	if err != nil {
		return err
	}

	challenge, err := h.recoverySvc.ResendCode(c.Context(), dto.ResendCodeInput{
		ChallengeID: challengeID,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusAccepted).JSON(presenter.RecoveryChallengeSuccessResponse(challenge))
}

// ResetPassword godoc
// @Summary Ustaw nowe hasło
// @Description Ustawia nowe hasło na podstawie reset tokenu z weryfikacji kodu. Unieważnia wszystkie aktywne sesje użytkownika.
// @Tags auth
// @Accept json
// @Param challenge_id path string true "ID wyzwania (UUID)"
// @Param request body presenter.ResetPasswordRequest true "Reset token i nowe hasło"
// @Success 204 "No Content"
// @Failure 400 {object} presenter.ErrorResponse "Nieprawidłowe body żądania"
// @Failure 401 {object} presenter.ErrorResponse "Nieprawidłowy lub wygasły reset token"
// @Failure 404 {object} presenter.ErrorResponse "Wyzwanie nie istnieje"
// @Failure 500 {object} presenter.ErrorResponse "Wewnętrzny błąd serwera"
// @Router /auth/password-recovery/{challenge_id}/reset [post]
func (h *PasswordRecoveryHandler) ResetPassword(c fiber.Ctx) error {
	challengeID, err := challengeIDFromPath(c)
	if err != nil {
		return err
	}

	var req presenter.ResetPasswordRequest
	if err := c.Bind().Body(&req); err != nil {
		return err
	}

	err = h.recoverySvc.ResetPassword(c.Context(), dto.ResetPasswordInput{
		ChallengeID: challengeID,
		ResetToken:  req.ResetToken,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
