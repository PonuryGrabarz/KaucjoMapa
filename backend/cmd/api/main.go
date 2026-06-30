package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/server"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

// @title           KaucjoMapa API
// @version         1.0
// @description     To jest serwer API dla aplikacji KaucjoMapa.
// @termsOfService  http://swagger.io/terms/
// @contact.name   Wsparcie Techniczne
// @contact.email  support@kaucjomapa.pl
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:3000
// @BasePath  /
func main() {
	server := server.New()

	if err := server.SeedAdminUser(context.Background()); err != nil {
		panic(fmt.Sprintf("seed admin user error: %s", err))
	}

	if err := server.ScrapPoints(); err != nil {
		fmt.Printf("scrap points error: %s", err)
	}

	server.RegisterFiberRoutes()

	port, _ := strconv.Atoi(os.Getenv("PORT"))

	err := server.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		panic(fmt.Sprintf("http server error: %s", err))
	}
}
