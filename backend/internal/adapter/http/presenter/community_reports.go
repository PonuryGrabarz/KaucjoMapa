package presenter

import "github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"

type CreateCommunityReportRequest struct {
	ReturnPointID string `json:"report_id" validate:"required,uuid"`
	Status        string `json:"status" validate:"required,oneof=dostepny nieczynny bark_informacji"`
	Description   string `json:"description" validate:"required"`
}

type CreateCommunityReportResponse struct {
	Report     entity.CommunityReport `json:"report"`
	Username   string                 `json:"username"`
	Reputation int32                  `json:"reputation"`
}

type ListReportsByReturnPointIDRequest struct {
	ReturnPointID string `json:"report_id" validate:"required,uuid"`
}

func CommunityReportSuccessResponse(report entity.CommunityReport, user entity.User) CreateCommunityReportResponse {
	return CreateCommunityReportResponse{
		Report:     report,
		Username:   user.Username,
		Reputation: user.Reputation,
	}
}

func ListCommunityReportsSuccessReponse(userReports []entity.UserCommunityReport) []CreateCommunityReportResponse {
	reportsRes := []CreateCommunityReportResponse{}

	for _, report := range userReports {
		reportsRes = append(reportsRes, CreateCommunityReportResponse{
			Report:     report.Report,
			Username:   report.Username,
			Reputation: report.Reputation,
		})
	}

	return reportsRes
}

func GetCommunityReportSuccessResponse(report entity.CommunityReport, user entity.User) CreateCommunityReportResponse {
	return CreateCommunityReportResponse{
		Report:   report,
		Username: user.Username,
	}
}
