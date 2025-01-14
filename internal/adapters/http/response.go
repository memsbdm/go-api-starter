package http

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"go-starter/internal/adapters/validator"
	"go-starter/internal/domain/entities"
	"net/http"
)

// response represents a generic API success response
type response[T any] struct {
	Success bool `json:"success"`
	Data    T    `json:"data,omitempty"`
}

// newSuccessResponse is a helper function to create a success response
func newSuccessResponse[T any](data T) response[T] {
	return response[T]{
		Success: true,
		Data:    data,
	}
}

// handleError sends an error response
func handleError(w http.ResponseWriter, err error) {
	status, ok := domainHttpErrMap[err]
	if !ok {
		status = http.StatusInternalServerError
	}

	errResp := newErrorResponse([]error{err})
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = encoder.Encode(errResp)
}

// emptyResponse represents a success response without data
type emptyResponse struct {
	Success bool `json:"success" example:"true"`
}

// errorResponse represents an error response body format
type errorResponse struct {
	Success  bool     `json:"success" example:"false"`
	Messages []string `json:"messages" example:"Error message 1, Error message 2"`
}

// newErrorResponse is a helper function to create an error response body
func newErrorResponse(errs []error) errorResponse {
	errsStr := make([]string, len(errs))
	for i, err := range errs {
		errsStr[i] = err.Error()
	}
	return errorResponse{
		Success:  false,
		Messages: errsStr,
	}
}

// handleSuccess sends a success response with the specified status code and optional data
func handleSuccess(w http.ResponseWriter, statusCode int, data any) {
	rsp := newSuccessResponse(data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	encoder := json.NewEncoder(w)
	_ = encoder.Encode(rsp)
}

// handleValidationError sends an error response for some specific request validation error
func handleValidationError(w http.ResponseWriter, errs []error) {
	w.Header().Set("Content-Type", "application/json")
	errRsp := newErrorResponse(errs)

	if errors.Is(errs[0], validator.ErrInvalidJSON) {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}

	_ = json.NewEncoder(w).Encode(errRsp)
}

// Custom responses

type healthResponse struct {
	Idle              string `json:"idle" example:"1"`
	InUse             string `json:"in_use" example:"0"`
	MaxIdleClosed     string `json:"max_idle_closed" example:"0"`
	MaxLifetimeClosed string `json:"max_lifetime_closed" example:"0"`
	Message           string `json:"message" example:"It's healthy'"`
	OpenConnections   string `json:"open_connections" example:"1"`
	Status            string `json:"status" example:"up"`
	WaitCount         string `json:"wait_count" example:"0"`
	WaitDuration      string `json:"wait_duration" example:"0s"`
}

// userResponse represents a user response body
type userResponse struct {
	ID              uuid.UUID `json:"id" example:"1"`
	Username        string    `json:"username" example:"john"`
	IsEmailVerified bool      `json:"is_email_verified" example:"true"`
}

// newUserResponse is a helper function to create a response body for handling user data
func newUserResponse(user *entities.User) userResponse {
	return userResponse{
		ID:              user.ID.UUID(),
		Username:        user.Username,
		IsEmailVerified: user.IsEmailVerified,
	}
}

// loginResponse represents a successful authentication body
type loginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// newLoginResponse is a helper function to create a response body for handling successful token
func newLoginResponse(accessToken, refreshToken string) loginResponse {
	return loginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

}
