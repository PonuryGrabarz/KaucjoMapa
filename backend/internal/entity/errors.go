package entity

import "errors"

var (
	ErrPointExistsNearby  = errors.New("point already exists nearby")
	ErrNotFound           = errors.New("resource not found")
	ErrForbidden          = errors.New("forbidden")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrServerError        = errors.New("internal server error")
	ErrBadRequest         = errors.New("bad request")
	ErrUnauthorized       = errors.New("unauthorized")
)
