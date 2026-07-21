package handler

import (
	"time"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/adapter/http/presenter"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/auth"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/user"
	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	userSvc   user.IUserService
	hasherSvc auth.IPasswordHasher
	tokenSvc  auth.ITokenSigner
}

func NewUserHandler(userSvc user.IUserService, hasherSvc auth.IPasswordHasher, tokenSvc auth.ITokenSigner) *UserHandler {
	return &UserHandler{
		userSvc:   userSvc,
		hasherSvc: hasherSvc,
		tokenSvc:  tokenSvc,
	}
}

func (h *UserHandler) Create(c fiber.Ctx) error {
	var req presenter.CreateUserRequest

	if err := c.Bind().Body(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	passwordHash, err := h.hasherSvc.HashPassword(req.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to hash password")
	}

	user, err := h.userSvc.Create(c.Context(), entity.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
		Reputation:   0,
		IsAdmin:      false,
		IsBanned:     false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create user")
	}

	accessToken, err := h.tokenSvc.MakeJWT(user.ID, user.IsAdmin, user.Reputation, user.Username, user.Email)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate access token")
	}

	return c.JSON(presenter.UserSuccessResponse(accessToken))
}
