package entity

import (
	"time"

	"github.com/google/uuid"
)

type (
	PointType    string
	PointStatus  string
	RefundMethod string
)

const (
	PointTypeButelkomat       PointType    = "butelkomat"
	PointTypeSklep            PointType    = "sklep"
	PointTypePunktSkupu       PointType    = "punkt_skupu"
	PointTypeInny             PointType    = "inny"
	PointStatusDostepny       PointStatus  = "dostepny"
	PointStatusNieczynny      PointStatus  = "nieczynny"
	PointStatusBrakInformacji PointStatus  = "brak_informacji"
	RefundMethodCash          RefundMethod = "gotowka"
	RefundMethodVoucher       RefundMethod = "bon"
	RefundMethodCard          RefundMethod = "karta"
)

type OpeningHours struct {
	MonSun string
	MonSat string
	MonFri string
	Sat    string
	Sun    any
}

type ReturnPoint struct {
	ID            uuid.UUID
	UserID        *uuid.UUID
	Name          string
	ShortName     string
	BgColor       string
	TextColor     string
	Type          PointType
	Latitude      float64
	Longitude     float64
	DistanceKm    *float64
	Address       string
	Pet           bool
	Can           bool
	Glass         bool
	RefundMethod  []RefundMethod
	OpeningHours  OpeningHours
	Description   string
	CurrentStatus PointStatus
	IsVerified    bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
