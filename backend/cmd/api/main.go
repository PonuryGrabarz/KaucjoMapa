package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
// @schemes   https
func main() {
	server := server.New()

	if err := server.SeedAdminUser(context.Background()); err != nil {
		panic(fmt.Sprintf("seed admin user error: %s", err))
	}

	if err := server.ScrapPoints(); err != nil {
		fmt.Printf("scrap points error: %s", err)
	}

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

	server.RegisterFiberRoutes()

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
