package server

import (
	"context"
	"encoding/json"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
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
	kaucjomatyList, err := s.db.Queries.GetAllReturnPoints(context.Background())
	if err != nil {
		return err
	}

	mapPoints := []MapPoint{}

	for _, kaucjomat := range kaucjomatyList {
		style := ResolvePointStyle(kaucjomat.Name)
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

func (s *FiberServer) CreateKaucjomatHandler(c fiber.Ctx) error {
	var mapPoint MapPoint

	err := json.Unmarshal(c.Body(), &mapPoint)
	if err != nil {
		return nil
	}

	// mapPoint, err := s.db.Queries.CreatePoint(context.Background(), database.CreatePointParams{
	// 	ID:        uuid.New(),
	// 	UserID:    userID,
	// 	Name:      kaucjomat.Brand,
	// 	Type:      database.PointTypeButelkomat,
	// 	Address:   address,
	// 	Latitude:  kaucjomat.Lat,
	// 	Longitude: kaucjomat.Lng,
	// 	OpeningHours: sql.NullString{
	// 		String: "8:00 - 20:00",
	// 		Valid:  true,
	// 	},
	// 	Description: sql.NullString{
	// 		String: "Automat do zwrotu butelek i puszek",
	// 		Valid:  true,
	// 	},
	// 	CurrentStatus:      database.PointStatusBrakInformacji,
	// 	ConfirmationsCount: 0,
	// 	RejectionsCount:    0,
	// 	IsVerified:         true,
	// 	CreatedAt:          time.Now(),
	// 	UpdatedAt:          time.Now(),
	// })

	return nil
}
