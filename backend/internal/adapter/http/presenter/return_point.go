package presenter

import (
	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/gofiber/fiber/v3"
)

type CreateReturnPointRequest struct {
	Name          string  `json:"name"`
	Type          string  `json:"type"`
	Address       string  `json:"address"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	OpeningHours  string  `json:"opening_hours"`
	Description   string  `json:"description"`
	CurrentStatus string  `json:"current_status"`
}

type UpdateKaucjomatRequest struct {
	Name         *string  `json:"name"`
	Type         *string  `json:"type"`
	Address      *string  `json:"address"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
	Pet          *bool    `json:"pet"`
	Can          *bool    `json:"can"`
	Glass        *bool    `json:"glass"`
	OpeningHours *string  `json:"opening_hours"`
	Description  *string  `json:"description"`
}

func CreateReturnPointSuccessResponse(p entity.ReturnPoint) *fiber.Map {
	style := ResolvePointStyle(p.Name)

	return &fiber.Map{
		"ID":            p.ID,
		"ShortName":     style.ShortName,
		"Name":          p.Name,
		"BgColor":       style.BgColor,
		"TextColor":     style.TextColor,
		"Latitude":      p.Latitude,
		"Longitude":     p.Longitude,
		"Address":       p.Address,
		"Pet":           p.Pet,
		"Can":           p.Can,
		"Glass":         p.Glass,
		"OpeningHours":  p.OpeningHours,
		"Description":   p.Description,
		"CurrentStatus": p.CurrentStatus,
	}
}
