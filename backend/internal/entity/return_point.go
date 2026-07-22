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
	ID            uuid.UUID
	UserID        *uuid.UUID
	Name          string
	Type          PointType
	Latitude      float64
	Longitude     float64
	Address       string
	Pet           bool
	Can           bool
	Glass         bool
	OpeningHours  string
	Description   string
	CurrentStatus PointStatus
	IsVerified    bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
