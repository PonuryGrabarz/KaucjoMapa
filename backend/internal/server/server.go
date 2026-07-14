package server

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/database"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/kaucjomatapi"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/server/middleware"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/encryptcookie"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "KaucjoMapa",
			AppName:      "KaucjoMapa",
			ErrorHandler: middleware.ErrorHandler,
		}),

		db: database.NewConn(),
	}

	server.Use(recover.New())
	server.Use(encryptcookie.New(encryptcookie.Config{
		Key: os.Getenv("COOKIE_SECRET"),
	}))
	server.Use(helmet.New())
	server.Use(requestid.New())

	server.Use(logger.New(logger.Config{
		Format: "${requestid} [${ip}]:${port} ${status} - ${method} ${path}\n",
	}))
	return server
}

func (s *FiberServer) SeedAdminUser(ctx context.Context) error {
	return s.db.SeedAdminUser(ctx)
}

func (s *FiberServer) ScrapPoints() error {
	kaucjomatyList, err := kaucjomatapi.GetKaucjomaty("data/kaucjomaty_polska.json")
	if err != nil {
		return err
	}

	user, err := s.db.Queries.GetUserByUsername(context.Background(), adminUsername)
	if err != nil {
		return err
	}

	userID := uuid.NullUUID{
		UUID: user.ID,
	}

	if user.ID == uuid.Nil {
		userID.Valid = false
	} else {
		userID.Valid = true
	}

	for _, kaucjomat := range kaucjomatyList {
		address := sql.NullString{
			String: kaucjomat.Address,
		}

		if kaucjomat.Address == "" {
			address.Valid = false
		} else {
			address.Valid = true
		}

		_, err := s.db.Queries.CreatePoint(context.Background(), database.CreatePointParams{
			ID:        uuid.New(),
			UserID:    userID,
			Name:      kaucjomat.Brand,
			Type:      database.PointTypeButelkomat,
			Address:   address,
			Latitude:  kaucjomat.Lat,
			Longitude: kaucjomat.Lng,
			Pet:       kaucjomat.Pet,
			Can:       kaucjomat.Can,
			Glass:     kaucjomat.Glass,
			OpeningHours: sql.NullString{
				String: "8:00 - 20:00",
				Valid:  true,
			},
			Description: sql.NullString{
				String: "Automat do zwrotu butelek i puszek",
				Valid:  true,
			},
			CurrentStatus:      database.PointStatusBrakInformacji,
			ConfirmationsCount: 0,
			RejectionsCount:    0,
			IsVerified:         true,
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		})
		if err != nil {
			fmt.Printf("Error creating point for kaucjomat %s: %v\n", kaucjomat.Brand, err)
		}

	}

	return nil
}
