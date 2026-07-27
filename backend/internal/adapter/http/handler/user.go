package handler

import (
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/adapter/http/presenter"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/dto"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/port"
	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	userSvc  usecase.IUserService
	tokenSvc port.ITokenSigner
	authSvc  usecase.IAuthService
}

func NewUserHandler(userSvc usecase.IUserService, tokenSvc port.ITokenSigner, authSvc usecase.IAuthService) *UserHandler {
	return &UserHandler{
		userSvc:  userSvc,
		tokenSvc: tokenSvc,
		authSvc:  authSvc,
	}
}

// Create godoc
// @Summary Zarejestruj użytkownika
// @Description Tworzy konto użytkownika i zwraca access token.
// @Tags users
// @Accept json
// @Produce json
// @Param request body presenter.CreateUserRequest true "Dane nowego użytkownika"
// @Success 200 {object} presenter.TokenResponse "Access token"
// @Failure 400 {object} presenter.ErrorResponse "Nieprawidłowe body żądania"
// @Failure 500 {object} presenter.ErrorResponse "Nie udało się utworzyć użytkownika"
// @Router /users [post]
func (h *UserHandler) Create(c fiber.Ctx) error {
	var req presenter.CreateUserRequest

	if err := c.Bind().Body(&req); err != nil {
		return err
	}

	user, err := h.userSvc.Create(c.Context(), dto.CreateUserInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return err
	}

	refreshToken, newToken, err := h.authSvc.Login(c.Context(), dto.LoginUserInput{
		Username: user.Username,
		Password: req.Password,
	})
	if err != nil {
		// Konto powstało, ale logowanie zawiodło — to błąd wewnętrzny.
		return entity.NewInternal(err)
	}

	cookie := fiber.Cookie{
		Name:     "RefreshToken",
		Value:    refreshToken.Token,
		Path:     "/api/auth",
		Expires:  refreshToken.ExpiresAt,
		HTTPOnly: true,
		SameSite: "strict",
		Secure:   true,
	}

	c.Cookie(&cookie)

	return c.JSON(presenter.UserSuccessResponse(newToken))
}
