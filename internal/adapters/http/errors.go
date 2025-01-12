package http

import (
	"go-starter/internal/domain"
	"net/http"
)

// domainHttpErrMap is a map of defined entities error messages and their corresponding http status codes
var domainHttpErrMap = map[error]int{
	domain.ErrInternal:     http.StatusInternalServerError,
	domain.ErrForbidden:    http.StatusForbidden,
	domain.ErrUnauthorized: http.StatusUnauthorized,

	domain.ErrInvalidToken: http.StatusUnauthorized,

	domain.ErrInvalidCredentials: http.StatusUnauthorized,

	domain.ErrInvalidUserId:             http.StatusBadRequest,
	domain.ErrUserNotFound:              http.StatusNotFound,
	domain.ErrUserUsernameAlreadyExists: http.StatusConflict,
}
