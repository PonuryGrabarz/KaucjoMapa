package presenter

type CreateCommunityReportVoteRequest struct {
	Vote string `json:"vote" validate:"required,oneof=confirm reject"`
}
