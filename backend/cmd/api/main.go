package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/adapter/http"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/adapter/repository/postgres"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/adapter/scraper/mapakaucji"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/infrastructure/config"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/infrastructure/crypto"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/infrastructure/database"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/infrastructure/token"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/auth"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/returnpoint"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/seed"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/user"
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
// @BasePath  /api
// @schemes   https
func main() {
	cfg := config.NewConfig()
	db := database.NewConn(cfg)

	scraper := mapakaucji.Scraper{}

	returnPointRepo := postgres.NewReturnPointRepo(db.Queries)
	userRepo := postgres.NewUserRepo(db.Queries)
	tokenRepo := postgres.NewRefreshTokenRepo(db.Queries)
	tokenSigner := token.NewTokenSigner(cfg.JWTSecret, time.Second*3600)
	hasher := crypto.NewHasher()

	returnPointSvc := returnpoint.NewService(returnPointRepo)
	userSvc := user.NewService(userRepo)
	authSvc := auth.NewRefreshTokenService(tokenRepo)

	seedSvc := seed.NewService(returnPointRepo, userRepo, &scraper, cfg.AdminUsername, hasher, cfg)

	err := seedSvc.SeedAdminUser(context.Background())
	if err != nil {
		log.Printf("Error seeding admin user: %v", err)
	}

	err = seedSvc.ScrapPoints(context.Background())
	if err != nil {
		log.Printf("Error scraping points: %v", err)
	}

	server := http.New(returnPointSvc, userSvc, authSvc, tokenSigner, hasher, cfg)

	certsDir, err := filepath.Abs("certs")
	if err != nil {
		log.Fatalf("resolve certs directory: %v", err)
	}

	log.Printf("certs directory: %s", certsDir)

	certPath := filepath.Join(certsDir, "ssl.cert")
	keyPath := filepath.Join(certsDir, "ssl.key")

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		fmt.Printf("Error loading TLS certificate: %v\n", err)
	}
	config := &tls.Config{Certificates: []tls.Certificate{cert}}

	port, _ := strconv.Atoi(os.Getenv("PORT"))

	ln, err := tls.Listen("tcp", fmt.Sprintf(":%v", port), config)
	if err != nil {
		fmt.Println(err)
	}

	err = server.Listener(ln)
	if err != nil {
		panic(fmt.Sprintf("http server error: %s", err))
	}
}
