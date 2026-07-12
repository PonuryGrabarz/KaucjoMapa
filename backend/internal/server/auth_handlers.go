package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/auth"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/database"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// LoginUserHandler godoc
// @Summary      Logowanie użytkownika
// @Description  Weryfikuje dane logowania i zwraca dane użytkownika z tokenem JWT oraz tokenem odświeżającym.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      server.LoginUserHandler.parameters  true  "Dane logowania"
// @Success      200      {object}  User
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /login [post]
func (s *FiberServer) LoginUserHandler(c fiber.Ctx) error {
	type parameters struct {
		Password string `json:"password"`
		Username string `json:"username"`
	}

	params := parameters{}

	err := json.Unmarshal(c.Request().Body(), &params)
	if err != nil {
		//respondWithError(w, 400, "Something went wrong")
		return err
	}

	user, err := s.db.Queries.GetUserByUsername(context.Background(), params.Username)
	if err != nil {
		//respondWithError(w, 401, "User not found")
		return err
	}

	isPasswordValid, err := auth.CheckPasswordHash(params.Password, user.PasswordHash)
	if err != nil {
		//respondWithError(w, 401, "Unauthorized")
		return err
	}

	if !isPasswordValid {
		//respondWithError(w, 401, "Unauthorized")
		return err
	}

	expTime := time.Second * 3600

	if isPasswordValid {
		loggedInUser := User{
			ID:         user.ID,
			CreatedAt:  user.CreatedAt.String(),
			UpdatedAt:  user.UpdatedAt.String(),
			Email:      user.Email,
			Username:   user.Username,
			Reputation: user.Reputation,
			IsAdmin:    user.IsAdmin,
			IsBanned:   user.IsBanned,
		}

		token, err := auth.MakeJWT(loggedInUser.ID, user.Username, user.Email, JWTSecret, expTime)
		if err != nil {
			//respondWithError(w, 500, "server error")
			return err
		}

		refreshToken := auth.MakeRefreshToken()

		_, err = s.db.Queries.CreateRefreshToken(c.Context(), database.CreateRefreshTokenParams{
			ID:        uuid.New(),
			Token:     refreshToken,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    loggedInUser.ID,
			ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
			RevokedAt: sql.NullTime{},
		})
		if err != nil {
			//respondWithError(w, 500, "server error")
			return err
		}

		loggedInUser.Token = token
		loggedInUser.RefreshToken = refreshToken

		//respondWithJSON(w, 200, loggedInUser)
		return c.JSON(loggedInUser)
	}

	return c.Err()
}

// refreshTokenHandler godoc
// @Summary      Odświeżenie tokenu JWT
// @Description  Przyjmuje refresh token w nagłówku Authorization Bearer i zwraca nowy token JWT.
// @Tags         auth
// @Produce      json
// @Param        Authorization  header    string  true  "Bearer refresh_token"
// @Success      200            {object}  map[string]string
// @Failure      401            {object}  map[string]string
// @Router       /auth/token [post]
func (s *FiberServer) refreshTokenHandler(c fiber.Ctx) error {
	token, err := auth.GetBearerToken(c.GetReqHeaders())
	if err != nil {
		//respondWithError(w, 401, err.Error())
		return err
	}

	userToken, err := s.db.Queries.GetUserToken(context.Background(), token)
	if err != nil {
		//respondWithError(w, 401, err.Error())
		return err
	}

	if userToken.RevokedAt.Valid || time.Now().Compare(userToken.ExpiresAt) != -1 {
		//respondWithError(w, 401, "Unauthorized")
		return err
	}

	user, err := s.db.Queries.GetUserById(c.Context(), userToken.UserID)
	if err != nil {
		return err
	}

	expTime := time.Second * 3600
	newToken, err := auth.MakeJWT(user.ID, user.Username, user.Email, JWTSecret, expTime)
	if err != nil {
		//respondWithError(w, 401, err.Error())
		return err
	}

	// respondWithJSON(w, 200, struct {
	// 	Token string `json:"token"`
	// }{
	// 	Token: newToken,
	// })

	return c.JSON(
		struct {
			Token string `json:"token"`
		}{
			Token: newToken,
		},
	)
}

// revokeTokenHandler godoc
// @Summary      Unieważnienie refresh tokenu
// @Description  Unieważnia refresh token przekazany w nagłówku Authorization Bearer.
// @Tags         auth
// @Param        Authorization  header  string  true  "Bearer refresh_token"
// @Success      204            "No Content"
// @Failure      401            {object}  map[string]string
// @Failure      500            {object}  map[string]string
// @Router       /auth/revoke [post]
func (s *FiberServer) revokeTokenHandler(c fiber.Ctx) error {
	token, err := auth.GetBearerToken(c.GetReqHeaders())
	if err != nil {
		//respondWithError(w, 401, err.Error())
		return err
	}

	userToken, err := s.db.Queries.GetUserToken(c.Context(), token)
	if err != nil {
		//respondWithError(w, 401, err.Error())
		return err
	}

	err = s.db.Queries.RevokeToken(c.Context(), userToken.Token)
	if err != nil {
		//respondWithError(w, 500, err.Error())
		return err
	}

	//respondWithJSON(w, 204, nil)
	return c.SendStatus(204)
}
