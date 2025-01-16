package http

import (
	"go-starter/internal/domain"
	"net/http"
)

// domainHttpErrMap maps domain-specific error types to their corresponding HTTP status codes.
var domainHttpErrMap = map[error]int{
	domain.ErrInternal:     http.StatusInternalServerError,
	domain.ErrForbidden:    http.StatusForbidden,
	domain.ErrUnauthorized: http.StatusUnauthorized,
	domain.ErrBadRequest:   http.StatusBadRequest,

	domain.ErrInvalidToken:       http.StatusUnauthorized,
	domain.ErrInvalidCredentials: http.StatusUnauthorized,

	domain.ErrInvalidUserId:        http.StatusBadRequest,
	domain.ErrUserNotFound:         http.StatusNotFound,
	domain.ErrUsernameAlreadyTaken: http.StatusConflict,

	// Validation errors

	// Users
	domain.ErrUsernameRequired:             http.StatusUnprocessableEntity,
	domain.ErrUsernameTooShort:             http.StatusUnprocessableEntity,
	domain.ErrUsernameTooLong:              http.StatusUnprocessableEntity,
	domain.ErrUsernameInvalid:              http.StatusUnprocessableEntity,
	domain.ErrPasswordRequired:             http.StatusUnprocessableEntity,
	domain.ErrPasswordsNotMatch:            http.StatusUnprocessableEntity,
	domain.ErrPasswordTooShort:             http.StatusUnprocessableEntity,
	domain.ErrPasswordConfirmationRequired: http.StatusUnprocessableEntity,
}
