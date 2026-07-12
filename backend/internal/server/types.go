package server

import "github.com/google/uuid"

type MapPoint struct {
	ReturnPointID uuid.UUID `json:"return_point_id"`
	ShortName     string    `json:"short_name"`
	BgColor       string    `json:"bg_color"`
	TextColor     string    `json:"text_color"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
}

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	Reputation   int32     `json:"reputation"`
	IsAdmin      bool      `json:"is_admin"`
	IsBanned     bool      `json:"is_banned"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
}
