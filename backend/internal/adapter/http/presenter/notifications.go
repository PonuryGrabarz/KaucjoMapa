package presenter

import (
	"time"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/google/uuid"
)

type CreateNotificationRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type NotificationResponse struct {
	ID          uuid.UUID  `json:"id"`
	UserID      *uuid.UUID `json:"user_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Seen        bool       `json:"seen"`
	CreatedAt   time.Time  `json:"created_at"`
}

func CreateNotificationSuccessReponse(notification entity.Notification) NotificationResponse {
	return NotificationResponse{
		ID:          notification.ID,
		UserID:      notification.UserID,
		Name:        notification.Name,
		Description: notification.Description,
		Seen:        notification.Seen,
		CreatedAt:   notification.CreatedAt,
	}

}

func ListNotificationsSuccessReponse(notifications []entity.Notification) []NotificationResponse {
	list := []NotificationResponse{}
	for _, n := range notifications {
		list = append(list, NotificationResponse{
			ID:          n.ID,
			UserID:      n.UserID,
			Name:        n.Name,
			Description: n.Description,
			Seen:        n.Seen,
			CreatedAt:   n.CreatedAt,
		})

	}

	return list
}
