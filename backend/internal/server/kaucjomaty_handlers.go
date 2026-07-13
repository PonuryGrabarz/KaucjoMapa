package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/auth"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/database"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
)

// GetKaucjomatyHandler godoc
// @Summary      Lista kaucjomatów
// @Description  Zwraca punkty zwrotu do wyświetlenia na mapie.
// @Tags         kaucjomaty
// @Produce      json
// @Success      200  {array}   MapPoint
// @Failure      500  {object}  map[string]string
// @Router       /kaucjomaty [get]
func (s *FiberServer) GetKaucjomatyHandler(c fiber.Ctx) error {
	pointType := c.Query("type")
	pointStatus := c.Query("status")
	// pet := c.Query("pet")
	// can := c.Query("can")
	// glass := c.Query(
	// 	"glass",
	// )

	kaucjomatyList, err := s.db.Queries.GetAllReturnPoints(context.Background())
	if err != nil {
		return err
	}

	mapPoints := []MapPoint{}

	for _, kaucjomat := range kaucjomatyList {
		style := ResolvePointStyle(kaucjomat.Name)
		if pointType != "" {
			fmt.Println(kaucjomat.Type)
			if kaucjomat.Type != database.PointType(pointType) {
				continue
			}
		}

		if pointStatus != "" {
			fmt.Println(kaucjomat.CurrentStatus)
			if kaucjomat.CurrentStatus != database.PointStatus(pointStatus) {
				continue
			}
		}

		mapPoint := MapPoint{
			ReturnPointID: uuid.New(),
			ShortName:     style.ShortName,
			BgColor:       style.BgColor,
			TextColor:     style.TextColor,
			Latitude:      kaucjomat.Latitude,
			Longitude:     kaucjomat.Longitude,
		}

		mapPoints = append(mapPoints, mapPoint)

	}

	return c.JSON(mapPoints)
}

// CreateKaucjomatHandler godoc
// @Summary      Dodanie kaucjomatu
// @Description  Tworzy nowy punkt zwrotu/kaucjomat dla zalogowanego użytkownika.
// @Tags         kaucjomaty
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization  header  string                   true  "Bearer JWT token"
// @Param        request        body    server.CreateKaucjomatRequest  true  "Dane kaucjomatu"
// @Success      200            {object}  MapPoint       "Utworzony kaucjomat"
// @Failure      400            {object}  map[string]string "Nieprawidłowe dane wejściowe"
// @Failure      401            {object}  map[string]string               "Brak lub nieprawidłowy token JWT"
// @Failure      500            {object}  map[string]string               "Błąd serwera"
// @Router       /kaucjomaty [post]
func (s *FiberServer) CreateKaucjomatHandler(c fiber.Ctx) error {
	var params CreateKaucjomatRequest

	err := json.Unmarshal(c.Body(), &params)
	if err != nil {
		return err
	}

	token, err := auth.GetBearerToken(c.GetReqHeaders())
	if err != nil {
		return err
	}

	claims, err := auth.ValidateJWT(token, os.Getenv("JWT_SECRET"))
	if err != nil {
		return err
	}

	point, err := s.db.Queries.GetPointByDistance(c.Context(), database.GetPointByDistanceParams{
		LlToEarth:   params.Latitude,
		LlToEarth_2: params.Longitude,
	})
	fmt.Println("point", point)
	if errors.Is(err, pgx.ErrNoRows) || point.ID != uuid.Nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Point already exists at this location",
		})
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return err
	}
	address := sql.NullString{
		String: params.Address,
	}
	openingHours := sql.NullString{
		String: params.OpeningHours,
	}
	description := sql.NullString{
		String: params.OpeningHours,
	}

	mapPoint, err := s.db.Queries.CreatePoint(context.Background(), database.CreatePointParams{
		ID: uuid.New(),
		UserID: uuid.NullUUID{
			UUID:  userID,
			Valid: true,
		},
		Name:               params.Name,
		Type:               database.PointType(params.Type),
		Address:            address,
		Latitude:           params.Latitude,
		Longitude:          params.Longitude,
		OpeningHours:       openingHours,
		Description:        description,
		CurrentStatus:      database.PointStatus(params.CurrentStatus),
		ConfirmationsCount: 0,
		RejectionsCount:    0,
		IsVerified:         true,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	})
	if err != nil {
		return err
	}

	style := ResolvePointStyle(mapPoint.Name)

	return c.JSON(MapPoint{
		ReturnPointID: mapPoint.ID,
		ShortName:     style.ShortName,
		BgColor:       style.BgColor,
		TextColor:     style.TextColor,
		Latitude:      mapPoint.Latitude,
		Longitude:     mapPoint.Longitude,
	})
}
