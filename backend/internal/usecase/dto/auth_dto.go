package dto

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
