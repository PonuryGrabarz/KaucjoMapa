package dto

import "github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"

type CreateCommunityReportOutput struct {
	Report   entity.CommunityReport
	Username string
}
