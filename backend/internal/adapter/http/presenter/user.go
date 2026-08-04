package presenter

type CreateUserRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type ChangePasswordRequest struct {
	OldPassword       string `json:"old_password" validate:"required"`
	NewPassword       string `json:"new_password" validate:"required,min=8"`
	RepeatNewPassword string `json:"repeat_new_password" validate:"required,min=8"`
}

type ChangeAvatarResponse struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	HasAvatar  bool   `json:"has_avatar"`
	Reputation int    `json:"reputation"`
	IsAdmin    bool   `json:"is_admin"`
}

func UserSuccessResponse(token string) TokenResponse {
	return TokenResponse{
		Token: token,
	}
}
