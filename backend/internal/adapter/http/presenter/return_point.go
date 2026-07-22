package presenter

import (
	"time"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
)

type CreateReturnPointRequest struct {
	Name          string  `json:"name"`
	Type          string  `json:"type" enums:"butelkomat,sklep,punkt_skupu,inny"`
	Address       string  `json:"address"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	Pet           bool    `json:"pet"`
	Can           bool    `json:"can"`
	Glass         bool    `json:"glass"`
	OpeningHours  string  `json:"opening_hours"`
	Description   string  `json:"description"`
	CurrentStatus string  `json:"current_status" enums:"dostepny,zapelniony,nieczynny,awaria,brak_informacji,do_weryfikacji"`
}

type UpdateKaucjomatRequest struct {
	Name         *string  `json:"name"`
	Type         *string  `json:"type" enums:"butelkomat,sklep,punkt_skupu,inny"`
	Address      *string  `json:"address"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
	Pet          *bool    `json:"pet"`
	Can          *bool    `json:"can"`
	Glass        *bool    `json:"glass"`
	OpeningHours *string  `json:"opening_hours"`
	Description  *string  `json:"description"`
}

type ReturnPointResponse struct {
	ID            string  `json:"id" format:"uuid"`
	ShortName     string  `json:"short_name"`
	Name          string  `json:"name"`
	BgColor       string  `json:"bg_color"`
	TextColor     string  `json:"text_color"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"Longitude"`
	Address       string  `json:"address"`
	Pet           bool    `json:"pet"`
	Can           bool    `json:"can"`
	Glass         bool    `json:"glass"`
	OpeningHours  string  `json:"opening_hours"`
	Description   string  `json:"description"`
	CurrentStatus string  `json:"current_status" enums:"dostepny,zapelniony,nieczynny,awaria,brak_informacji,do_weryfikacji"`
}

type ReturnPointPage struct {
	ReturnPoints []ReturnPointResponse `json:"return_points"`
	MaxPage      int64                 `json:"max_page"`
}

// ReturnPointDetailsResponse describes the full return point returned by an update.
type ReturnPointDetailsResponse struct {
	ID            string    `json:"id" format:"uuid"`
	UserID        *string   `json:"user_id" format:"uuid"`
	Name          string    `json:"name"`
	Type          string    `json:"type" enums:"butelkomat,sklep,punkt_skupu,inny"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	Address       string    `json:"address"`
	Pet           bool      `json:"pet"`
	Can           bool      `json:"can"`
	Glass         bool      `json:"glass"`
	OpeningHours  string    `json:"opening_hours"`
	Description   string    `json:"description"`
	CurrentStatus string    `json:"current_status" enums:"dostepny,zapelniony,nieczynny,awaria,brak_informacji,do_weryfikacji"`
	IsVerified    bool      `json:"is_verified"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func CreateReturnPointSuccessResponse(p entity.ReturnPoint) ReturnPointResponse {
	style := ResolvePointStyle(p.Name)

	return ReturnPointResponse{
		ID:            p.ID.String(),
		ShortName:     style.ShortName,
		Name:          p.Name,
		BgColor:       style.BgColor,
		TextColor:     style.TextColor,
		Latitude:      p.Latitude,
		Longitude:     p.Longitude,
		Address:       p.Address,
		Pet:           p.Pet,
		Can:           p.Can,
		Glass:         p.Glass,
		OpeningHours:  p.OpeningHours,
		Description:   p.Description,
		CurrentStatus: string(p.CurrentStatus),
	}
}

func GetReturnPointSuccessResponse(returnPoints []entity.ReturnPoint, maxPage int64) ReturnPointPage {
	mapPointsDetails := []ReturnPointResponse{}

	for _, point := range returnPoints {
		mapPointsDetails = append(mapPointsDetails, CreateReturnPointSuccessResponse(point))
	}

	return ReturnPointPage{
		ReturnPoints: mapPointsDetails,
		MaxPage:      maxPage,
	}
}
