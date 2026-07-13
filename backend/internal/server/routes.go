package server

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
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
			"Set-Cookie",
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
	s.Post("/kaucjomaty", s.CreateKaucjomatHandler)
	s.Post("/users", s.CreateUserHandler)
	s.Post("/login", s.LoginUserHandler)
	s.Post("/auth/token", s.refreshTokenHandler)
	s.Post("/auth/revoke", s.revokeTokenHandler)
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
