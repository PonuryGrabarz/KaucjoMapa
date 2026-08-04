package entity

import (
	"time"

	"github.com/google/uuid"
)

type Category string
type ErrorReportStatus string

const (
	CategoryAddressLocationError Category = "adres_lokalizacja"
	CategoryDescriptionError     Category = "opis_punktu"
	CategoryRefundMethodError    Category = "formy_rozliczenia"
	CategoryOpeningHoursError    Category = "godziny_otwarcia"
	CategoryAcceptedBottlesError Category = "przyjmowane_opakowania"
	CategoryOtherError           Category = "inne_dane"
)

const (
	ErrorReportStatusPending  ErrorReportStatus = "oczekujace"
	ErrorReportStatusReviewed ErrorReportStatus = "rozpatrzone"
	ErrorReportStatusRejected ErrorReportStatus = "odrzucone"
)

type ErrorReport struct {
	ID                  uuid.UUID
	UserID              *uuid.UUID
	ReturnPointID       *uuid.UUID
	Categories          []Category
	Description         string
	SuggestedCorrection string
	Status              ErrorReportStatus
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
