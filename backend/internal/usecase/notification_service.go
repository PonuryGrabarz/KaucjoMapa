package usecase

import (
	"context"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/google/uuid"
)

type INotificationsService interface {
	Create(ctx context.Context, notification entity.Notification) (entity.Notification, error)
	List(ctx context.Context, id uuid.UUID) ([]entity.Notification, error)
	GetByID(ctx context.Context, id uuid.UUID) (entity.Notification, error)
	Update(ctx context.Context, id uuid.UUID, notification entity.Notification) (entity.Notification, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type notificationsService struct {
	repo    INotificationsRepository
	userSvc IUserService
}

func NewNotificationService(repo INotificationsRepository, userSvc IUserService) INotificationsService {
	return &notificationsService{
		repo:    repo,
		userSvc: userSvc,
	}
}

func (s *notificationsService) Create(ctx context.Context, notification entity.Notification) (entity.Notification, error) {
	newNotification, err := s.repo.Create(ctx, notification)
	if err != nil {
		return entity.Notification{}, err
	}

	return newNotification, nil

}

func (s *notificationsService) List(ctx context.Context, id uuid.UUID) ([]entity.Notification, error) {
	notifications, err := s.repo.List(ctx, id)
	if err != nil {
		return []entity.Notification{}, err
	}

	return notifications, nil
}

func (s *notificationsService) GetByID(ctx context.Context, id uuid.UUID) (entity.Notification, error) {
	notification, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return entity.Notification{}, err
	}

	return notification, nil
}

func (s *notificationsService) Update(ctx context.Context, id uuid.UUID, notification entity.Notification) (entity.Notification, error) {
	updatedNotification, err := s.repo.Update(ctx, id, notification)
	if err != nil {
		return entity.Notification{}, err
	}

	return updatedNotification, nil

}

func (s *notificationsService) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
