package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/auth"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/database"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
)

// GetKaucjomatyHandler godoc
//
// @Summary      Lista kaucjomatów
// @Description  Zwraca punkty zwrotu do wyświetlenia na mapie. Parametr details=true rozszerza odpowiedź o szczegóły punktów.
// @Tags         kaucjomaty
// @Produce      json
// @Param        type       query  string  false  "Typ punktu"     Enums(butelkomat, sklep, punkt_skupu, inny)
// @Param        status     query  string  false  "Status punktu"  Enums(dostepny, zapelniony, nieczynny, awaria, brak_informacji, do_weryfikacji)
// @Param        details    query  bool    false  "Zwróć szczegółowe dane punktów"  default(false)
// @Param        page       query  int     false  "Numer strony (używany dla details=true)"  minimum(1)  default(1)
// @Param        page_size  query  int     false  "Liczba punktów na stronie (używana dla details=true)"  minimum(1)  default(10)
// @Param        query      query  string  false  "Nazwa sklepu lub adres"
// @Success      200  {array}   MapPoint           "Lista punktów; dla details=true zwracany jest model MapPointDetails"
// @Failure      400  {object}  map[string]string  "Nieprawidłowy numer strony lub rozmiar strony"
// @Failure      500  {object}  map[string]string  "Błąd pobierania punktów zwrotu"
// @Router       /kaucjomaty [get]
func (s *FiberServer) GetKaucjomatyHandler(c fiber.Ctx) error {
	pointType := c.Query("type")
	pointStatus := c.Query("status")
	details := c.Query("details")
	pageQuery := c.Query("page", "1")
	pageSizeQuery := c.Query("page_size", "10")
	query := strings.TrimSpace(c.Query("query", ""))

	page, err := strconv.Atoi(pageQuery)
	if err != nil || page < 1 {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid page number")
	}
	pageSize, err := strconv.Atoi(pageSizeQuery)
	if err != nil || pageSize < 1 {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid page size")
	}

	kaucjomatyList, err := s.db.Queries.ListReturnPoints(c.Context(), database.ListReturnPointsParams{
		Query: sql.NullString{
			String: query,
			Valid:  query != "",
		},
		PointType: database.NullPointType{
			PointType: database.PointType(pointType),
			Valid:     pointType != "",
		},
		PointStatus: database.NullPointStatus{
			PointStatus: database.PointStatus(pointStatus),
			Valid:       pointStatus != "",
		},
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch return points")
	}

	mapPoints := []MapPoint{}
	mapPointsDetails := []MapPointDetails{}

	for _, kaucjomat := range kaucjomatyList {
		style := ResolvePointStyle(kaucjomat.Name)

		if details == "true" {
			mapPointDetails := MapPointDetails{
				ReturnPointID: kaucjomat.ID,
				Name:          kaucjomat.Name,
				ShortName:     style.ShortName,
				BgColor:       style.BgColor,
				TextColor:     style.TextColor,
				Latitude:      kaucjomat.Latitude,
				Longitude:     kaucjomat.Longitude,
				Address:       kaucjomat.Address.String,
				Pet:           kaucjomat.Pet,
				Can:           kaucjomat.Can,
				Glass:         kaucjomat.Glass,
				OpeningHours:  kaucjomat.OpeningHours.String,
				Description:   kaucjomat.Description.String,
				CurrentStatus: string(kaucjomat.CurrentStatus),
			}

			mapPointsDetails = append(mapPointsDetails, mapPointDetails)
		} else {
			mapPoint := MapPoint{
				ReturnPointID: kaucjomat.ID,
				ShortName:     style.ShortName,
				BgColor:       style.BgColor,
				TextColor:     style.TextColor,
				Latitude:      kaucjomat.Latitude,
				Longitude:     kaucjomat.Longitude,
			}

			mapPoints = append(mapPoints, mapPoint)
		}

	}

	maxPage := len(mapPointsDetails) / pageSize

	if details == "true" && maxPage > 0 && page <= maxPage {
		return c.JSON(mapPointsDetails[pageSize*(page-1) : min(pageSize*page, len(mapPointsDetails))])
	}

	return c.JSON(mapPoints)
}

// CreateKaucjomatHandler godoc
//
//	@Summary		Dodanie kaucjomatu
//	@Description	Tworzy nowy punkt zwrotu/kaucjomat dla zalogowanego użytkownika.
//	@Description	Dozwolone wartości type: butelkomat, sklep, punkt_skupu, inny.
//	@Description	Dozwolone wartości current_status: dostepny, zapelniony, nieczynny, awaria, brak_informacji, do_weryfikacji.
//	@Tags			kaucjomaty
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							true	"Bearer access token"
//	@Param			request			body		server.CreateKaucjomatRequest	true	"Dane kaucjomatu"
//	@Success		200				{object}	MapPoint						"Utworzony kaucjomat"
//	@Failure		400				{object}	map[string]string				"Nieprawidłowe dane wejściowe"
//	@Failure		401				{object}	map[string]string				"Brak lub nieprawidłowy access token"
//	@Failure		500				{object}	map[string]string				"Błąd serwera"
//	@Router			/kaucjomaty [post]
func (s *FiberServer) CreateKaucjomatHandler(c fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*auth.JWTClaims)
	fmt.Println(ok)
	fmt.Println(claims)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid JWT claims")
	}

	var params CreateKaucjomatRequest

	err := json.Unmarshal(c.Body(), &params)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	point, err := s.db.Queries.GetPointByDistance(c.Context(), database.GetPointByDistanceParams{
		LlToEarth:   params.Latitude,
		LlToEarth_2: params.Longitude,
	})
	if errors.Is(err, pgx.ErrNoRows) || point.ID != uuid.Nil {
		return fiber.NewError(fiber.StatusBadRequest, "Point already exists nearby")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Invalid user ID in JWT claims")
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

	mapPoint, err := s.db.Queries.CreatePoint(c.Context(), database.CreatePointParams{
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
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create return point")
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

// DeleteKaucjomatHandler godoc
//
//	@Summary		Usunięcie kaucjomatu
//	@Description	Usuwa punkt zwrotu. Operacja jest dostępna wyłącznie dla administratora.
//	@Tags			kaucjomaty
//	@Produce		json
//	@Param			Authorization	header	string	true	"Bearer access token"
//	@Param			id				path	string	true	"Identyfikator kaucjomatu"	Format(uuid)
//	@Success		204				"No Content"
//	@Failure		400				{object}	map[string]string	"Nieprawidłowy identyfikator kaucjomatu"
//	@Failure		401				{object}	map[string]string	"Brak lub nieprawidłowy access token"
//	@Failure		403				{object}	map[string]string	"Użytkownik nie jest administratorem"
//	@Failure		500				{object}	map[string]string	"Błąd serwera"
//	@Router			/kaucjomaty/{id} [delete]
func (s *FiberServer) DeleteKaucjomatHandler(c fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*auth.JWTClaims)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid JWT claims")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Invalid user ID in JWT claims")
	}

	user, err := s.db.Queries.GetUserById(c.Context(), userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Faild to fetch user")
	}

	if !user.IsAdmin {
		return fiber.NewError(fiber.StatusForbidden, "User is not an admin")
	}
	mapPointID := c.Params("id")
	mapPointUUID, err := uuid.Parse(mapPointID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid map point ID")
	}

	err = s.db.Queries.DeletePointById(c.Context(), mapPointUUID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete return point")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// UpdateKaucjomatHandler godoc
//
//	@Summary		Aktualizacja kaucjomatu
//	@Description	Aktualizuje dane punktu zwrotu dla zalogowanego użytkownika.
//	@Tags			kaucjomaty
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string													true	"Bearer access token"
//	@Param			id				path		string													true	"Identyfikator kaucjomatu"	Format(uuid)
//	@Param			request			body		server.UpdateKaucjomatHandler.UpdateKaucjomatRequest	true	"Dane kaucjomatu"
//	@Success		200				{object}	map[string]interface{}									"Zaktualizowany kaucjomat"
//	@Failure		400				{object}	map[string]string										"Nieprawidłowy identyfikator lub body"
//	@Failure		401				{object}	map[string]string										"Brak lub nieprawidłowy access token"
//	@Failure		500				{object}	map[string]string										"Nie znaleziono kaucjomatu lub błąd serwera"
//	@Router			/kaucjomaty/{id} [put]
func (s *FiberServer) UpdateKaucjomatHandler(c fiber.Ctx) error {
	claims, ok := c.Locals("claims").(*auth.JWTClaims)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid JWT claims")
	}

	mapPointID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid map point ID")
	}

	_, err = uuid.Parse(claims.Subject)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Invalid user ID in JWT claims")
	}

	_, err = s.db.Queries.GetPointById(c.Context(), mapPointID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Return point not found")
	}

	type UpdateKaucjomatRequest struct {
		Name         *string  `json:"name"`
		Type         *string  `json:"type"`
		Address      *string  `json:"address"`
		Latitude     *float64 `json:"latitude"`
		Longitude    *float64 `json:"longitude"`
		Pet          *bool    `json:"pet"`
		Can          *bool    `json:"can"`
		Glass        *bool    `json:"glass"`
		OpeningHours *string  `json:"opening_hours"`
		Description  *string  `json:"description"`
	}

	var reqBody UpdateKaucjomatRequest
	if err := json.Unmarshal(c.Body(), &reqBody); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	params := database.UpdatePointByIdParams{
		ID: mapPointID,
	}

	if reqBody.Name != nil {
		params.Name = sql.NullString{String: *reqBody.Name, Valid: true}
	}

	if reqBody.Type != nil {
		params.Type = database.NullPointType{PointType: database.PointType(*reqBody.Type), Valid: true}
	}

	if reqBody.Address != nil {
		params.Address = sql.NullString{String: *reqBody.Address, Valid: true}
	}

	if reqBody.Latitude != nil {
		params.Latitude = sql.NullFloat64{Float64: *reqBody.Latitude, Valid: true}
	}

	if reqBody.Longitude != nil {
		params.Longitude = sql.NullFloat64{Float64: *reqBody.Longitude, Valid: true}
	}

	if reqBody.Pet != nil {
		params.Pet = sql.NullBool{Bool: *reqBody.Pet, Valid: true}
	}

	if reqBody.Can != nil {
		params.Can = sql.NullBool{Bool: *reqBody.Can, Valid: true}
	}

	if reqBody.Glass != nil {
		params.Glass = sql.NullBool{Bool: *reqBody.Glass, Valid: true}
	}

	if reqBody.OpeningHours != nil {
		params.OpeningHours = sql.NullString{String: *reqBody.OpeningHours, Valid: true}
	}

	if reqBody.Description != nil {
		params.Description = sql.NullString{String: *reqBody.Description, Valid: true}
	}

	updatedPoint, err := s.db.Queries.UpdatePointById(c.Context(), params)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update return point")
	}

	return c.Status(fiber.StatusOK).JSON(updatedPoint)
}
