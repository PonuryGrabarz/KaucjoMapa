package postgres

import (
	"database/sql"
	"errors"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
	"github.com/lib/pq"
)

// mapError tłumaczy błędy z warstwy bazy danych (database/sql + sterownik
// lib/pq) na błędy domenowe entity.AppError. Dzięki temu szczegóły PostgreSQL
// (kody SQLSTATE, nazwy kolumn, komunikaty constraintów) nigdy nie wyciekają
// do warstwy HTTP ani do klienta — surowy błąd zostaje tylko w logach.
func mapError(err error) error {
	if err == nil {
		return nil
	}

	// Brak wiersza to najczęstszy przypadek — nie znaleziono zasobu.
	if errors.Is(err, sql.ErrNoRows) {
		return entity.NewNotFound("resource not found")
	}

	// Błędy specyficzne dla PostgreSQL niosą kod SQLSTATE.
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23505": // unique_violation — naruszenie unikalności
			return entity.NewConflict("resource already exists")
		case "23503", // foreign_key_violation — wskazano nieistniejący zasób
			"23502", // not_null_violation
			"23514", // check_violation
			"22P02", // invalid_text_representation (np. zły uuid / enum)
			"22007", // invalid_datetime_format
			"22008", // datetime_field_overflow
			"22003": // numeric_value_out_of_range
			return entity.NewBadRequest("invalid request")
		}
	}

	// Wszystko inne (padnięte połączenie, timeout, nieznany kod) chowamy za
	// generyczny błąd serwera, zachowując surowy błąd do logów.
	return entity.NewInternal(err)
}
