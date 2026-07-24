package entity

import "errors"

// ErrorKind klasyfikuje błąd domenowy. Warstwa HTTP mapuje Kind na kod stanu.
type ErrorKind int

const (
	KindInternal ErrorKind = iota
	KindNotFound
	KindBadRequest
	KindUnauthorized
	KindForbidden
	KindConflict
	KindInvalidCredentials
)

// FieldError opisuje pojedynczy błąd walidacji pola. Jest bezpieczny do
// pokazania klientowi, bo dotyczy wyłącznie danych, które sam przysłał.
type FieldError struct {
	Field   string
	Rule    string
	Param   string
	Message string
}

// AppError to jedyny typ błędu, jaki warstwy usecase/repo/handler propagują w
// górę. Message jest kuratorowany i trafia do klienta; Err przechowuje surowy
// błąd wewnętrzny wyłącznie do logów i NIGDY nie jest eksponowany.
type AppError struct {
	Kind    ErrorKind
	Message string
	Fields  []FieldError
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error { return e.Err }

func NewNotFound(msg string) *AppError { return &AppError{Kind: KindNotFound, Message: msg} }

func NewForbidden(msg string) *AppError { return &AppError{Kind: KindForbidden, Message: msg} }

func NewUnauthorized(msg string) *AppError { return &AppError{Kind: KindUnauthorized, Message: msg} }

func NewConflict(msg string) *AppError { return &AppError{Kind: KindConflict, Message: msg} }

func NewInvalidCredentials() *AppError {
	return &AppError{Kind: KindInvalidCredentials, Message: "invalid credentials"}
}

func NewBadRequest(msg string, fields ...FieldError) *AppError {
	return &AppError{Kind: KindBadRequest, Message: msg, Fields: fields}
}

// NewInternal chowa surowy błąd: klient dostaje generyk, log dostaje Err.
func NewInternal(err error) *AppError {
	return &AppError{Kind: KindInternal, Message: "internal server error", Err: err}
}

// AsAppError wyciąga *AppError z łańcucha błędów (obsługuje wrapping).
func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

// IsKind zwraca true, gdy err jest AppError danego rodzaju.
func IsKind(err error, kind ErrorKind) bool {
	appErr, ok := AsAppError(err)
	return ok && appErr.Kind == kind
}
