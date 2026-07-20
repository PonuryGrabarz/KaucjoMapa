package entity

import (
	"time"

	"github.com/google/uuid"
)

type (
	PointType   string
	PointStatus string
)

const (
	PointTypeButelkomat       PointType   = "butelkomat"
	PointTypeSklep            PointType   = "sklep"
	PointTypePunktSkupu       PointType   = "punkt_skupu"
	PointTypeInny             PointType   = "inny"
	PointStatusDostepny       PointStatus = "dostepny"
	PointStatusZapelniony     PointStatus = "zapelniony"
	PointStatusNieczynny      PointStatus = "nieczynny"
	PointStatusAwaria         PointStatus = "awaria"
	PointStatusBrakInformacji PointStatus = "bark_informacji"
	PointStatusDoWeryfikacji  PointStatus = "do_weryfikacji"
)

type ReturnPoint struct {
	ID            uuid.UUID   `json:"id"`
	UserID        *uuid.UUID  `json:"user_id"`
	Name          string      `json:"name"`
	Type          PointType   `json:"type"`
	Latitude      float64     `json:"latitude"`
	Longitude     float64     `json:"longitude"`
	Address       string      `json:"address"`
	Pet           bool        `json:"pet"`
	Can           bool        `json:"can"`
	Glass         bool        `json:"glass"`
	OpeningHours  string      `json:"opening_hours"`
	Description   string      `json:"description"`
	CurrentStatus PointStatus `json:"current_status"`
	IsVerified    bool        `json:"is_verified"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}
