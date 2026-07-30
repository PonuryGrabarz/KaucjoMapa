package config

import (
	"fmt"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	DBURL              string
	Port               int
	AdminUsername      string
	AdminEmail         string
	AdminPassword      string
	JWTSecret          string
	CookieSecret       string
	GoogleClientID     string
	GoogleClientSecret string
	GootleRedirectURL  string
	SMTPHost           string
	SMTPPort           string
	SMTPUsername       string
	SMTPPassword       string
	SMTPFrom           string
	AvatarStoragePath  string
}

func NewConfig() *Config {
	var (
		database = os.Getenv("DB_DATABASE")
		password = os.Getenv("DB_PASSWORD")
		username = os.Getenv("DB_USERNAME")
		port     = os.Getenv("DB_PORT")
		host     = os.Getenv("DB_HOST")
		schema   = os.Getenv("DB_SCHEMA")
	)

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)
	p, err := strconv.Atoi(port)
	if err != nil {
		p = 3000
	}

	return &Config{
		DBURL:              dbURL,
		Port:               p,
		AdminUsername:      os.Getenv("ADMIN_USERNAME"),
		AdminEmail:         os.Getenv("ADMIN_EMAIL"),
		AdminPassword:      os.Getenv("ADMIN_PASSWORD"),
		JWTSecret:          os.Getenv("JWT_SECRET"),
		CookieSecret:       os.Getenv("COOKIE_SECRET"),
		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GootleRedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		SMTPHost:           os.Getenv("SMTP_HOST"),
		SMTPPort:           os.Getenv("SMTP_PORT"),
		SMTPUsername:       os.Getenv("SMTP_USERNAME"),
		SMTPPassword:       os.Getenv("SMTP_PASSWORD"),
		SMTPFrom:           os.Getenv("SMTP_FROM"),
		AvatarStoragePath:  avatarStoragePath(),
	}
}

func avatarStoragePath() string {
	if path := os.Getenv("AVATAR_STORAGE_PATH"); path != "" {
		return path
	}
	return "./uploads/avatars"
}
