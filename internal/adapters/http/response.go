package http

import (
	"encoding/json"
	"errors"
	"go-starter/internal/adapters/validator"
	"go-starter/internal/domain/entities"
	"net/http"
)

// response represents a response body format
type response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// newResponse is a helper function to create a response body
func newResponse(success bool, message string, data any) response {
	return response{
		Success: success,
		Message: message,
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
	rsp := newResponse(true, "Success", data)
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
	ID       int    `json:"id" example:"1"`
	Username string `json:"username" example:"john"`
}

// newUserResponse is a helper function to create a response body for handling user data
func newUserResponse(user *entities.User) userResponse {
	return userResponse{
		ID:       user.ID,
		Username: user.Username,
	}
}
