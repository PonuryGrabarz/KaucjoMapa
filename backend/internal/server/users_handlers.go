package server

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/auth"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/database"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

var JWTSecret = os.Getenv("JWT_SECRET")

// CreateUserHandler godoc
// @Summary      Rejestracja użytkownika
// @Description  Tworzy nowe konto użytkownika i zwraca dane utworzonego użytkownika.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body      server.CreateUserHandler.parameters  true  "Dane rejestracji użytkownika"
// @Success      200      {object}  User
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /users [post]
func (s *FiberServer) CreateUserHandler(c fiber.Ctx) error {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Username string `json:"username"`
	}

	params := parameters{}

	err := json.Unmarshal(c.Request().Body(), &params)
	if err != nil {
		//respondWithError(w, 400, "Something went wrong")
		return err
	}

	hashpassword, err := auth.HashPassword(params.Password)
	if err != nil {
		//respondWithError(w, 500, "something went wrong")
		return err
	}

	user, err := s.db.Queries.CreateUser(context.Background(), database.CreateUserParams{
		ID:           uuid.New(),
		Email:        params.Email,
		PasswordHash: hashpassword,
		Username:     params.Username,
		Reputation:   0,
		IsAdmin:      false,
		IsBanned:     false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	})
	if err != nil {
		//respondWithError(w, 500, "something went wrong")
		return err
	}

	newUser := User{
		ID:         user.ID,
		Username:   user.Username,
		Email:      user.Email,
		Reputation: user.Reputation,
		IsAdmin:    user.IsAdmin,
		IsBanned:   user.IsBanned,
		CreatedAt:  user.CreatedAt.String(),
		UpdatedAt:  user.UpdatedAt.String(),
	}

	//respondWithJson(w, 201, newUser)
	return c.JSON(newUser)
}
