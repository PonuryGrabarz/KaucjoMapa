package server

import (
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/database"
	"github.com/gofiber/fiber/v3"
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
		}),

		db: database.NewConn(),
	}

	return server
}
