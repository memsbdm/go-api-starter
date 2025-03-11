package responses

import (
	"encoding/json"
	"errors"
	"go-starter/internal/adapters/server/apierrors"
	"go-starter/internal/adapters/validator"
	"net/http"
)

// Response represents a generic API success response structure.
type Response[T any] struct {
	Success bool `json:"success"`
	Data    T    `json:"data,omitempty"`
}

// NewSuccessResponse is a helper function that creates a new success response.
func NewSuccessResponse[T any](data T) Response[T] {
	return Response[T]{
		Success: true,
		Data:    data,
	}
}

// HandleError sends an error response to the client.
// It determines the appropriate HTTP status code based on the provided error
// and returns a standardized error response format.
func HandleError(w http.ResponseWriter, err error) {
	status, ok := apierrors.DomainHttpErrMap[err]
	if !ok {
		status = http.StatusInternalServerError
	}

	if status == http.StatusUnprocessableEntity {
		HandleValidationError(w, []error{err})
		return
	}

	errResp := NewErrorResponse([]error{err})
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = encoder.Encode(errResp)
}

// EmptyResponse represents a success response without any data.
// It indicates that the request was successful, but no additional information is provided.
type EmptyResponse struct {
	Success bool `json:"success" example:"true"`
}

// ErrorResponse represents the format of an error response body.
type ErrorResponse struct {
	Success  bool     `json:"success" example:"false"`
	Messages []string `json:"messages" example:"Error message 1,Error message 2"`
}

// NewErrorResponse is a helper function that creates an error response body from a slice of error messages.
func NewErrorResponse(errs []error) ErrorResponse {
	errsStr := make([]string, len(errs))
	for i, err := range errs {
		errsStr[i] = err.Error()
	}
	return ErrorResponse{
		Success:  false,
		Messages: errsStr,
	}
}

// HandleSuccess sends a success response with the specified status code and optional data.
// It encodes the response in JSON format and sets the appropriate HTTP headers.
func HandleSuccess(w http.ResponseWriter, statusCode int, data any) {
	rsp := NewSuccessResponse(data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	encoder := json.NewEncoder(w)
	_ = encoder.Encode(rsp)
}

// HandleValidationError sends an error response specifically for request validation errors.
// It sets the appropriate HTTP status code based on the type of validation error.
func HandleValidationError(w http.ResponseWriter, errs []error) {
	w.Header().Set("Content-Type", "application/json")
	errRsp := NewErrorResponse(errs)

	if errors.Is(errs[0], validator.ErrInvalidJSON) {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}

	_ = json.NewEncoder(w).Encode(errRsp)
}
