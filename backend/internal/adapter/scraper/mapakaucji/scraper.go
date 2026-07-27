package mapakaucji

import (
	"context"
	"fmt"
	"time"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase"
	"github.com/google/uuid"
)

type Scraper struct{}

type FeatureCollection[T any] struct {
	Type     string `json:"type"`
	Features []T    `json:"features"`
}

var _ usecase.DataSource = (*Scraper)(nil) // gwarancja spełnienia portu

func (s *Scraper) FetchReturnPoints(ctx context.Context) ([]entity.ReturnPoint, error) {
	var returnPoints []entity.ReturnPoint

	kaucjomatyList, err := GetKaucjomaty[[]Kaucjomat](
		"data/kaucjomaty_polska.json",
	)
	if err != nil {
		return returnPoints, err
	}

	zabkaGeoJSON, err := GetKaucjomaty[FeatureCollection[MapPoint]]("data/zabka_sklepy.geojson")
	if err != nil {
		return returnPoints, err
	}
	zabkaKaucjomatyList := zabkaGeoJSON.Features

	for _, retPoint := range kaucjomatyList {
		returnPoints = append(returnPoints, entity.ReturnPoint{
			ID:           uuid.New(),
			UserID:       nil,
			Name:         retPoint.Brand,
			Type:         entity.PointTypeButelkomat,
			Address:      fmt.Sprintf("%s %s", retPoint.City, retPoint.Address),
			Latitude:     retPoint.Lat,
			Longitude:    retPoint.Lng,
			Pet:          retPoint.Pet,
			Can:          retPoint.Can,
			Glass:        retPoint.Glass,
			RefundMethod: []entity.RefundMethod{entity.RefundMethodCash, entity.RefundMethodVoucher},
			OpeningHours: entity.OpeningHours{
				MonFri: "06:00:00 - 23:00:00",
			},
			Description:   "Automat do zwrotu butelek i puszek",
			CurrentStatus: entity.PointStatusDostepny,
			IsVerified:    true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		})
	}

	for _, retPoint := range zabkaKaucjomatyList {
		openingHours := entity.OpeningHours{}
		if retPoint.Properties.OpeningHours.MonSun != "" {
			openingHours.MonSun = retPoint.Properties.OpeningHours.MonSun
		}

		if retPoint.Properties.OpeningHours.MonSat != "" {
			openingHours.MonSat = retPoint.Properties.OpeningHours.MonSat
		}

		if retPoint.Properties.OpeningHours.MonFri != "" {
			openingHours.MonFri = retPoint.Properties.OpeningHours.MonFri
		}

		if retPoint.Properties.OpeningHours.Sat != "" {
			openingHours.Sat = retPoint.Properties.OpeningHours.Sat
		}

		if retPoint.Properties.OpeningHours.Sun != "" {
			openingHours.Sun = retPoint.Properties.OpeningHours.Sun
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
			RefundMethod:  []entity.RefundMethod{entity.RefundMethodCash, entity.RefundMethodVoucher},
			Description:   "Automat do zwrotu butelek i puszek",
			CurrentStatus: entity.PointStatusDostepny,
			IsVerified:    true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		})
	}

	return returnPoints, nil
}
