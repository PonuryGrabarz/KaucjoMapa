package database

import (
	"context"
	"fmt"
	"os"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/auth"
	"github.com/google/uuid"
)

func (s Service) SeedAdminUser(ctx context.Context) error {
	username := os.Getenv("ADMIN_USERNAME")
	email := os.Getenv("ADMIN_EMAIL")
	password := os.Getenv("ADMIN_PASSWORD")

	if username == "" {
		username = "admin"
	}

	if email == "" {
		email = "admin@localhost"
	}

	if password == "" {
		return fmt.Errorf("ADMIN_PASSWORD is required to seed admin user")
	}

	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return fmt.Errorf("hash admin password: %w", err)
	}

	err = s.Queries.UpsertAdminUser(ctx, UpsertAdminUserParams{
		ID:           uuid.New(),
		Email:        email,
		Username:     username,
		PasswordHash: string(passwordHash),
	})
	if err != nil {
		return fmt.Errorf("upsert admin user: %w", err)
	}

	return nil
}
