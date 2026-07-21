package seed

import (
	"context"
	"fmt"
	"time"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/infrastructure/config"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/auth"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/returnpoint"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/user"
	"github.com/google/uuid"

	_ "github.com/joho/godotenv/autoload"
)

type Service struct {
	points        returnpoint.IReturnPointRepository
	users         user.IUserRepository
	source        DataSource
	adminUsername string
	hasher        auth.IPasswordHasher
	cfg           *config.Config
}

func NewService(returnPoint returnpoint.IReturnPointRepository, user user.IUserRepository, src DataSource, adminUsername string, hasher auth.IPasswordHasher, cfg *config.Config) *Service {
	return &Service{points: returnPoint, users: user, source: src, adminUsername: adminUsername, hasher: hasher, cfg: cfg}
}

func (s *Service) ScrapPoints(ctx context.Context) error {
	admin, err := s.users.GetUserByUsername(ctx, s.adminUsername)
	if err != nil {
		return err
	}

	pts, err := s.source.FetchReturnPoints(ctx)
	if err != nil {
		return err
	}

	for _, p := range pts {
		p.UserID = &admin.ID
		if _, err := s.points.Create(ctx, p); err != nil {
			fmt.Printf("Error creating point for kaucjomat %s: %v\n", p.Name, err)
			continue
		}
	}
	return nil
}

func (s *Service) SeedAdminUser(ctx context.Context) error {
	password := s.cfg.AdminPassword
	email := s.cfg.AdminEmail
	username := s.cfg.AdminUsername

	hashPassword, err := s.hasher.HashPassword(password)
	if err != nil {
		return err
	}

	user := entity.User{
		ID:           uuid.New(),
		Email:        email,
		Username:     username,
		PasswordHash: hashPassword,
		IsAdmin:      true,
		IsBanned:     false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	adminUser, err := s.users.Create(ctx, user)
	if err != nil {
		return err
	}
	fmt.Printf("Admin user created: %s\n", adminUser.Username)

	return nil
}
