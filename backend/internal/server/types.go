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

type MapPointDetails struct {
	ReturnPointID uuid.UUID `json:"return_point_id"`
	Name          string    `json:"name"`
	ShortName     string    `json:"short_name"`
	BgColor       string    `json:"bg_color"`
	TextColor     string    `json:"text_color"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	Address       string    `json:"address"`
	Pet           bool      `json:"pet"`
	Can           bool      `json:"can"`
	Glass         bool      `json:"glass"`
	OpeningHours  string    `json:"opening_hours"`
	Description   string    `json:"description"`
	CurrentStatus string    `json:"current_status"`
}

type User struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	Username   string    `json:"username"`
	Token      string    `json:"token"`
	Reputation int32     `json:"reputation"`
	IsAdmin    bool      `json:"is_admin"`
	IsBanned   bool      `json:"is_banned"`
	CreatedAt  string    `json:"created_at"`
	UpdatedAt  string    `json:"updated_at"`
}

type CreateKaucjomatRequest struct {
	Name          string  `json:"name"`
	Type          string  `json:"type"`
	Address       string  `json:"address"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	OpeningHours  string  `json:"opening_hours"`
	Description   string  `json:"description"`
	CurrentStatus string  `json:"current_status"`
}
