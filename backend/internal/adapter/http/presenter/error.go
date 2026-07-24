package presenter

// FieldErrorResponse opisuje jeden błąd walidacji pojedynczego pola.
type FieldErrorResponse struct {
	Field   string `json:"field" example:"current_status"`
	Rule    string `json:"rule" example:"oneof"`
	Param   string `json:"param,omitempty" example:"dostepny awaria brak_informacji"`
	Message string `json:"message" example:"current_status must be one of: dostepny awaria brak_informacji"`
}

type ErrorResponse struct {
	Error   string               `json:"error" example:"Invalid request body"`
	Details []FieldErrorResponse `json:"details,omitempty"`
}
