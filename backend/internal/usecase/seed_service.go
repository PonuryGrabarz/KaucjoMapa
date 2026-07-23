package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase/dto"
	"github.com/google/uuid"

	_ "github.com/joho/godotenv/autoload"
)

type Service struct {
	points        IReturnPointRepository
	users         IUserRepository
	source        DataSource
	adminUsername string
	hasher        IPasswordHasher
}

func NewSeedService(returnPoint IReturnPointRepository, user IUserRepository, src DataSource, adminUsername string, hasher IPasswordHasher) *Service {
	return &Service{points: returnPoint, users: user, source: src, adminUsername: adminUsername, hasher: hasher}
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
			continue
		}
	}
	return nil
}

func (s *Service) SeedAdminUser(ctx context.Context, in dto.SeedAdminInput) error {
	hashPassword, err := s.hasher.HashPassword(in.Password)
	if err != nil {
		return err
	}

	user := entity.User{
		ID:           uuid.New(),
		Email:        in.Email,
		Username:     in.Username,
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
