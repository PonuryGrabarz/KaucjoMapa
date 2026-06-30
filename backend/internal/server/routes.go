package server

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/google/uuid"
	"github.com/yokeTH/gofiber-scalar/scalar/v3"
)

var adminUsername = os.Getenv("ADMIN_USERNAME")

func (s *FiberServer) RegisterFiberRoutes() {
	s.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
			"PATCH",
		},
		AllowHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
		},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	s.Get("/", s.HelloWorldHandler)

	openAPISpec, err := os.ReadFile("./docs/swagger.json")
	if err != nil {
		log.Fatalf("Nie udało się wczytać pliku OpenAPI: %v", err)
	}

	s.Get("/docs/*", scalar.New(
		scalar.Config{
			BasePath:          "/",
			Path:              "/docs",
			Title:             "KaucjoMapa API Documentation",
			FileContentString: string(openAPISpec),
		},
	))

	s.Get("/kaucjomaty", s.GetKaucjomatyHandler)
}

// HelloWorldHandler godoc
// @Summary      Sprawdzenie działania API
// @Description  Zwraca prostą wiadomość potwierdzającą, że API działa.
// @Tags         health
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       / [get]
func (s *FiberServer) HelloWorldHandler(c fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}

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
