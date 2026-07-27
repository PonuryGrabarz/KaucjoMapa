package database

import (
	"database/sql"
	"log"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/infrastructure/config"
	_ "github.com/joho/godotenv/autoload"
)

type Service struct {
	DB      *sql.DB
	Queries *Queries
}

var dbInstance *Service

func NewConn(cfg *config.Config) Service {
	if dbInstance != nil {
		return *dbInstance
	}

	connStr := cfg.DBURL

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	dbInstance = &Service{
		DB:      db,
		Queries: New(db),
	}

	return *dbInstance
}
