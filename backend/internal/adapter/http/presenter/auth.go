package presenter

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

func LoginSuccessResponse(token string) TokenResponse {
	return TokenResponse{
		Token: token,
	}
}

func RefreshTokenSuccessResponse(token string) TokenResponse {
	return TokenResponse{
		Token: token,
	}
}
