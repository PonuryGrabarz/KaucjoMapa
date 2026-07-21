package mapakaucji

import (
	"context"
	"fmt"
	"time"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/kaucjomatapi"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/seed"
	"github.com/google/uuid"
)

type Scraper struct { /* baseURL, http.Client */
}

var _ seed.DataSource = (*Scraper)(nil) // gwarancja spełnienia portu

func (s *Scraper) FetchReturnPoints(ctx context.Context) ([]entity.ReturnPoint, error) {
	var returnPoints []entity.ReturnPoint

	kaucjomatyList, err := kaucjomatapi.GetKaucjomaty[[]kaucjomatapi.Kaucjomat](
		"data/kaucjomaty_polska.json",
	)
	if err != nil {
		return returnPoints, err
	}

	zabkaGeoJSON, err := kaucjomatapi.GetKaucjomaty[kaucjomatapi.FeatureCollection[kaucjomatapi.MapPoint]]("data/zabka_sklepy.geojson")
	if err != nil {
		return returnPoints, err
	}
	zabkaKaucjomatyList := zabkaGeoJSON.Features

	for _, retPoint := range kaucjomatyList {
		returnPoints = append(returnPoints, entity.ReturnPoint{
			ID:            uuid.New(),
			UserID:        nil,
			Name:          retPoint.Brand,
			Type:          entity.PointTypeButelkomat,
			Address:       fmt.Sprintf("%s %s", retPoint.City, retPoint.Address),
			Latitude:      retPoint.Lat,
			Longitude:     retPoint.Lng,
			Pet:           retPoint.Pet,
			Can:           retPoint.Can,
			Glass:         retPoint.Glass,
			OpeningHours:  "8:00 - 20:00",
			Description:   "Automat do zwrotu butelek i puszek",
			CurrentStatus: entity.PointStatusDostepny,
			IsVerified:    true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		})
	}

	for _, retPoint := range zabkaKaucjomatyList {
		openingHours := "8:00 - 20:00"
		if retPoint.Properties.OpeningHours.MonSun != "" {
			openingHours = retPoint.Properties.OpeningHours.MonSun
		}

		if retPoint.Properties.OpeningHours.MonSat != "" {
			openingHours = retPoint.Properties.OpeningHours.MonSat
		}

		returnPoints = append(returnPoints, entity.ReturnPoint{
			ID:            uuid.New(),
			UserID:        nil,
			Name:          "Żabka",
			Type:          entity.PointTypeButelkomat,
			Address:       fmt.Sprintf("%s %s", retPoint.Properties.Town, retPoint.Properties.Street),
			Latitude:      retPoint.Properties.Lat,
			Longitude:     retPoint.Properties.Lon,
			Pet:           true,
			Can:           true,
			Glass:         false,
			OpeningHours:  openingHours,
			Description:   "Automat do zwrotu butelek i puszek",
			CurrentStatus: entity.PointStatusDostepny,
			IsVerified:    true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		})
	}

	return returnPoints, nil
}
