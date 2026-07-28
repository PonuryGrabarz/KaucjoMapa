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
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/dto"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

// @title KaucjoMapa API
// @version 1.0
// @description API do wyszukiwania i zarządzania punktami zwrotu opakowań kaucyjnych.
// @contact.name Wsparcie techniczne
// @contact.email support@kaucjomapa.pl
// @license.name Apache 2.0
// @license.url https://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3000
// @BasePath /api
// @schemes https
// @accept json
// @produce json
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Wartość w formacie: Bearer <JWT>
func main() {
	cfg := config.NewConfig()
	db := database.NewConn(cfg)

	scraper := mapakaucji.Scraper{}

	returnPointRepo := postgres.NewReturnPointRepo(db.Queries)
	userRepo := postgres.NewUserRepo(db.Queries)
	tokenRepo := postgres.NewRefreshTokenRepo(db.Queries)
	reportsRepo := postgres.NewCommunityReportsRepo(db.Queries)
	voteRepo := postgres.NewCommunityReportVoteRepo(db.Queries)
	notificationsRepo := postgres.NewNotificationsRepo(db.Queries)
	bookmarkRepo := postgres.NewBookmarkRepo(db.Queries)
	returnPlanRepo := postgres.NewReturnPlanRepo(db.Queries)
	oAuthRepo := postgres.NewOAuthRepo(db.Queries)
	tokenSigner := token.NewTokenSigner(cfg.JWTSecret, time.Second*3600)
	hasher := crypto.NewHasher()

	userSvc := usecase.NewUserService(userRepo, hasher, tokenSigner)
	returnPointSvc := usecase.NewReturnPointService(returnPointRepo, userSvc, tokenSigner)
	authSvc := usecase.NewAuthService(tokenRepo, userSvc, hasher, tokenSigner)
	reportSvc := usecase.NewCommunityReportsService(reportsRepo, userSvc)
	voteSvc := usecase.NewCommunityReportVoteService(voteRepo, userSvc, reportSvc)
	notificationsSvc := usecase.NewNotificationService(notificationsRepo, userSvc)
	bookmarkSvc := usecase.NewBookmarkService(bookmarkRepo)
	returnPlanSvc := usecase.NewReturnPlanService(returnPlanRepo, returnPointRepo)
	oAuthSvc := usecase.NewGoogleOAuthService(oAuthRepo, tokenRepo, tokenSigner)

	seedSvc := usecase.NewSeedService(returnPointRepo, userRepo, &scraper, cfg.AdminUsername, hasher)

	err := seedSvc.SeedAdminUser(context.Background(), dto.SeedAdminInput{
		Username: cfg.AdminUsername,
		Password: cfg.AdminPassword,
		Email:    cfg.AdminEmail,
	})
	if err != nil {
		log.Printf("Error seeding admin user: %v", err)
	}

	err = seedSvc.ScrapPoints(context.Background())
	if err != nil {
		fmt.Printf("Error scraping points: %v", err)
	}

	server := http.New(returnPointSvc, userSvc, authSvc, reportSvc, voteSvc, notificationsSvc, bookmarkSvc, returnPlanSvc, oAuthSvc, tokenSigner, hasher, cfg)

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
