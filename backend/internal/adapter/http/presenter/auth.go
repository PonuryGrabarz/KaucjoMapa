package presenter

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
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
