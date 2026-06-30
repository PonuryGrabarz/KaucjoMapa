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
