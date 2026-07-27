package presenter

import (
	"time"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/google/uuid"
)

type OpeningHours struct {
	MonSun string `json:"mon-sun"`
	MonSat string `json:"mon-sat"`
	MonFri string `json:"mon-fri"`
	Sat    string `json:"sat"`
	Sun    any    `json:"sun"`
}

type CreateReturnPointRequest struct {
	Name          string       `json:"name" validate:"required,max=255"`
	Type          string       `json:"type" validate:"required,oneof=butelkomat sklep punkt_skupu inny" enums:"butelkomat,sklep,punkt_skupu,inny"`
	Address       string       `json:"address"`
	Latitude      float64      `json:"latitude" validate:"latitude"`
	Longitude     float64      `json:"longitude" validate:"longitude"`
	Pet           bool         `json:"pet"`
	Can           bool         `json:"can"`
	Glass         bool         `json:"glass"`
	RefundMethod  []string     `json:"refund_method" validate:"required,dive,oneof=gotowka bon karta" enums:"gotowka,bon,karta"`
	OpeningHours  OpeningHours `json:"opening_hours"`
	Description   string       `json:"description"`
	CurrentStatus string       `json:"current_status" validate:"required,oneof=dostepny awaria brak_informacji" enums:"dostepny,awaria,brak_informacji"`
}

type UpdateKaucjomatRequest struct {
	Name         *string       `json:"name" validate:"omitempty,required,max=255"`
	Type         *string       `json:"type" validate:"omitempty,oneof=butelkomat sklep punkt_skupu inny" enums:"butelkomat,sklep,punkt_skupu,inny"`
	Address      *string       `json:"address"`
	Latitude     *float64      `json:"latitude" validate:"omitempty,latitude"`
	Longitude    *float64      `json:"longitude" validate:"omitempty,longitude"`
	Pet          *bool         `json:"pet"`
	Can          *bool         `json:"can"`
	Glass        *bool         `json:"glass"`
	RefundMethod *[]string     `json:"refund_method" validate:"omitempty,dive,oneof=gotowka bon karta" enums:"gotowka,bon,karta"`
	OpeningHours *OpeningHours `json:"opening_hours"`
	Description  *string       `json:"description"`
}

type ReturnPointDetailsResponse struct {
	ID            string       `json:"id" format:"uuid"`
	UserID        *uuid.UUID   `json:"user_id" format:"uuid"`
	Name          string       `json:"name"`
	ShortName     string       `json:"short_name"`
	BgColor       string       `json:"bg_color"`
	TextColor     string       `json:"text_color"`
	Type          string       `json:"type" enums:"butelkomat,sklep,punkt_skupu,inny"`
	Latitude      float64      `json:"latitude"`
	Longitude     float64      `json:"longitude"`
	Address       string       `json:"address"`
	Pet           bool         `json:"pet"`
	Can           bool         `json:"can"`
	Glass         bool         `json:"glass"`
	OpeningHours  OpeningHours `json:"opening_hours"`
	RefundMethod  []string     `json:"refund_method" enums:"gotowka,bon,karta"`
	Description   string       `json:"description"`
	CurrentStatus string       `json:"current_status" validate:"required,oneof=dostepny awaria brak_informacji" enums:"dostepny,awaria,brak_informacji"`
	IsVerified    bool         `json:"is_verified"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

type Pagination struct {
	Page       int64 `json:"page"`
	PageSize   int64 `json:"page_size"`
	TotalPages int64 `json:"total_pages"`
	TotalCount int64 `json:"total_count"`
}

type ReturnPointPage struct {
	ReturnPoints []ReturnPointDetailsResponse `json:"return_points"`
	Pagination   Pagination                   `json:"pagination"`
}

func CreateReturnPointSuccessResponse(p entity.ReturnPoint) ReturnPointDetailsResponse {
	style := ResolvePointStyle(p.Name)
	refundMethod := []string{}

	for _, method := range p.RefundMethod {
		refundMethod = append(refundMethod, string(method))
	}

	return ReturnPointDetailsResponse{
		ID:           p.ID.String(),
		UserID:       p.UserID,
		Type:         string(p.Type),
		ShortName:    style.ShortName,
		Name:         p.Name,
		BgColor:      style.BgColor,
		TextColor:    style.TextColor,
		Latitude:     p.Latitude,
		Longitude:    p.Longitude,
		Address:      p.Address,
		Pet:          p.Pet,
		Can:          p.Can,
		Glass:        p.Glass,
		RefundMethod: refundMethod,
		OpeningHours: OpeningHours{
			MonFri: p.OpeningHours.MonFri,
			MonSat: p.OpeningHours.MonSat,
			MonSun: p.OpeningHours.MonSun,
			Sat:    p.OpeningHours.Sat,
			Sun:    p.OpeningHours.Sun,
		},
		Description:   p.Description,
		CurrentStatus: string(p.CurrentStatus),
		IsVerified:    p.IsVerified,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

func GetReturnPointSuccessResponse(returnPoints []entity.ReturnPoint, maxPage int64, page int64, pageSize int64, totalCount int64) ReturnPointPage {
	mapPointsDetails := []ReturnPointDetailsResponse{}

	for _, point := range returnPoints {
		mapPointsDetails = append(mapPointsDetails, CreateReturnPointSuccessResponse(point))
	}

	return ReturnPointPage{
		ReturnPoints: mapPointsDetails,
		Pagination: Pagination{
			Page:       page,
			PageSize:   pageSize,
			TotalPages: maxPage,
			TotalCount: totalCount,
		},
	}
}
