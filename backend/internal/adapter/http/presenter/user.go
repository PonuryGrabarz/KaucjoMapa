package presenter

type CreateUserRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func UserSuccessResponse(token string) TokenResponse {
	return TokenResponse{
		Token: token,
	}
}
