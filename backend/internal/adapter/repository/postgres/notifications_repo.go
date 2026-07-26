package postgres

import (
	"context"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/infrastructure/database"
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/usecase"
	"github.com/google/uuid"
)

type NotificationsRepo struct {
	q *database.Queries
}

func NewNotificationsRepo(q *database.Queries) *NotificationsRepo {
	return &NotificationsRepo{
		q: q,
	}
}

var _ usecase.INotificationsRepository = (*NotificationsRepo)(nil)

func (r *NotificationsRepo) Create(ctx context.Context, notification entity.Notification) (entity.Notification, error) {
	row, err := r.q.CreateUserNotification(ctx, toCreateNotificationParams(notification))
	if err != nil {
		return entity.Notification{}, mapError(err)
	}

	return toNotificationEntity(row), nil

}

func (r *NotificationsRepo) List(ctx context.Context, id uuid.UUID) ([]entity.Notification, error) {
	rows, err := r.q.ListUserNotificationsByUserId(ctx, toNullUUID(&id))
	if err != nil {
		return []entity.Notification{}, mapError(err)
	}

	notifications := []entity.Notification{}

	for _, n := range rows {
		notifications = append(notifications, toNotificationEntity(n))
	}

	return notifications, nil
}

func (r *NotificationsRepo) GetByID(ctx context.Context, id uuid.UUID) (entity.Notification, error) {
	row, err := r.q.GetNotificationById(ctx, id)
	if err != nil {
		return entity.Notification{}, mapError(err)
	}

	return toNotificationEntity(row), nil
}

func (r *NotificationsRepo) Update(ctx context.Context, id uuid.UUID, notification entity.Notification) (entity.Notification, error) {
	row, err := r.q.UpdateUserNotification(ctx, toUpdateNotificationParams(notification, id))
	if err != nil {
		return entity.Notification{}, mapError(err)
	}

	return toNotificationEntity(row), nil
}

func (r *NotificationsRepo) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeleteUserNotification(ctx, id)
	if err != nil {
		return mapError(err)
	}

	return nil
}
