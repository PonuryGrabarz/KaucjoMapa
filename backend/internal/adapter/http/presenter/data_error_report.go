package presenter

import (
	"time"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/google/uuid"
)

type CreateDataErrorReportRequest struct {
	Categories          []string `json:"categories" validate:"required"`
	Description         string   `json:"description" validate:"required"`
	SuggestedCorrection string   `json:"suggested_correction"`
}

type ChangeDataErrorReportStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=oczekujace rozpatrzone odrzucone"`
}

type DataErrorReportSuccessResponse struct {
	ID                  uuid.UUID                `json:"id"`
	UserID              *uuid.UUID               `json:"user_id"`
	ReturnPointID       *uuid.UUID               `json:"return_point_id"`
	Categories          []entity.Category        `json:"categories"`
	Description         string                   `json:"description"`
	SuggestedCorrection string                   `json:"suggested_correction"`
	Status              entity.ErrorReportStatus `json:"status"`
	CreatedAt           time.Time                `json:"created_at"`
	UpdatedAt           time.Time                `json:"updated_at"`
}

func CreateDataErrorReportSuccessResponse(report entity.ErrorReport) DataErrorReportSuccessResponse {
	return DataErrorReportSuccessResponse{
		ID:                  report.ID,
		UserID:              report.UserID,
		ReturnPointID:       report.ReturnPointID,
		Categories:          report.Categories,
		Description:         report.Description,
		SuggestedCorrection: report.SuggestedCorrection,
		Status:              report.Status,
		CreatedAt:           report.CreatedAt,
		UpdatedAt:           report.UpdatedAt,
	}
}
