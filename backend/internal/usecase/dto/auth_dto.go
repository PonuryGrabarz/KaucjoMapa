package dto

import "github.com/google/uuid"

type LoginUserInput struct {
	Username string
	Password string
}

type RefreshTokenInput struct {
	RefreshToken string
	AccessToken  string
}

type RevokeTokenInput struct {
	RefreshToken string
	AccessToken  string
}

type ChangePasswordInput struct {
	UserID            uuid.UUID
	OldPassword       string
	NewPassword       string
	RepeatNewPassword string
}
