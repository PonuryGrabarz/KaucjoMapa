package presenter

import "github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"

type CreateCommunityReportRequest struct {
	ReturnPointID string `json:"report_id" validate:"required,uuid"`
	Status        string `json:"status" validate:"required,oneof=dostepny nieczynny bark_informacji"`
	Description   string `json:"description" validate:"required"`
}

type CreateCommunityReportResponse struct {
	Report   entity.CommunityReport `json:"report"`
	Username string                 `json:"username"`
}

type ListReportsByReturnPointIDRequest struct {
	ReturnPointID string `json:"report_id" validate:"required,uuid"`
}

func CommunityReportSuccessResponse(report entity.CommunityReport, user entity.User) CreateCommunityReportResponse {
	return CreateCommunityReportResponse{
		Report:   report,
		Username: user.Username,
	}
}

func ListCommunityReportsSuccessReponse(reports map[string]entity.CommunityReport) []CreateCommunityReportResponse {
	reportsRes := []CreateCommunityReportResponse{}

	for username, report := range reports {
		reportsRes = append(reportsRes, CreateCommunityReportResponse{
			Report:   report,
			Username: username,
		})
	}

	return reportsRes
}
