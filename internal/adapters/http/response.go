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
	Success  bool     `json:"success"`
	Messages []string `json:"messages"`
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

// userResponse represents a user response body
type userResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

// newUserResponse is a helper function to create a response body for handling user data
func newUserResponse(user *entities.User) userResponse {
	return userResponse{
		ID:       user.ID,
		Username: user.Username,
	}
}
