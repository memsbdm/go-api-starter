package http

import (
	"encoding/json"
	"errors"
	"go-starter/internal/adapters/validator"
	"go-starter/internal/domain/entities"
	"go-starter/pkg/env"
	"net/http"
	"time"
)

// response represents a generic API success response structure.
type response[T any] struct {
	Success bool `json:"success"`
	Data    T    `json:"data,omitempty"`
}

// newSuccessResponse is a helper function that creates a new success response.
func newSuccessResponse[T any](data T) response[T] {
	return response[T]{
		Success: true,
		Data:    data,
	}
}

// handleError sends an error response to the client.
// It determines the appropriate HTTP status code based on the provided error
// and returns a standardized error response format.
func handleError(w http.ResponseWriter, err error) {
	status, ok := domainHttpErrMap[err]
	if !ok {
		status = http.StatusInternalServerError
	}

	if status == http.StatusUnprocessableEntity {
		handleValidationError(w, []error{err})
		return
	}

	errResp := newErrorResponse([]error{err})
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = encoder.Encode(errResp)
}

// emptyResponse represents a success response without any data.
// It indicates that the request was successful, but no additional information is provided.
type emptyResponse struct {
	Success bool `json:"success" example:"true"`
}

// errorResponse represents the format of an error response body.
type errorResponse struct {
	Success  bool     `json:"success" example:"false"`
	Messages []string `json:"messages" example:"Error message 1, Error message 2"`
}

// newErrorResponse is a helper function that creates an error response body from a slice of error messages.
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

// handleSuccess sends a success response with the specified status code and optional data.
// It encodes the response in JSON format and sets the appropriate HTTP headers.
func handleSuccess(w http.ResponseWriter, statusCode int, data any) {
	rsp := newSuccessResponse(data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	encoder := json.NewEncoder(w)
	_ = encoder.Encode(rsp)
}

// handleValidationError sends an error response specifically for request validation errors.
// It sets the appropriate HTTP status code based on the type of validation error.
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

// healthResponse represents the structure of the health check response.
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

// userResponse represents the structure of a response body containing user information.
type userResponse struct {
	ID              string    `json:"id" example:"6b947a32-8919-4974-9ef3-048a556b0b75"`
	CreatedAt       time.Time `json:"created_at" example:"2024-08-15T16:23:33.455225Z"`
	UpdatedAt       time.Time `json:"updated_at" example:"2025-01-15T14:29:33.455225Z"`
	Name            string    `json:"name" example:"John Doe"`
	Username        string    `json:"username" example:"john"`
	IsEmailVerified bool      `json:"is_email_verified" example:"true"`
}

// newUserResponse is a helper function that creates a userResponse from a user entity.
func newUserResponse(user *entities.User) userResponse {
	return userResponse{
		ID:              user.ID.String(),
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
		Name:            user.Name,
		Username:        user.Username,
		IsEmailVerified: user.IsEmailVerified,
	}
}

// getUserByIDResponse represents the structure of a response body containing user information.
type getUserByIDResponse struct {
	ID       string `json:"id" example:"6b947a32-8919-4974-9ef3-048a556b0b75"`
	Name     string `json:"name" example:"John Doe"`
	Username string `json:"username" example:"john"`
}

// newGetUserByIDResponse is a helper function that creates a userResponse from a user entity.
func newGetUserByIDResponse(user *entities.User) getUserByIDResponse {
	return getUserByIDResponse{
		ID:       user.ID.String(),
		Name:     user.Name,
		Username: user.Username,
	}
}

// loginResponse represents the structure of a response body for successful authentication.
type loginResponse struct {
	Tokens authTokensResponse `json:"tokens"`
	User   userResponse       `json:"user"`
}

// authTokensResponse represents the structure of a response body with an access token, a refresh token
// and their expiration time in ms.
type authTokensResponse struct {
	AccessToken             string `json:"access_token" example:"eyJhbGciOi..."`
	RefreshToken            string `json:"refresh_token" example:"eyJhbGciOi..."`
	AccessTokenExpiredInMs  int64  `json:"access_token_expired_in_ms" example:"50000"`
	RefreshTokenExpiredInMs int64  `json:"refresh_token_expired_in_ms" example:"890000"`
}

// newAuthTokensResponse is a helper function that creates a authTokensResponse from an auth tokens entity.
func newAuthTokensResponse(tokens *entities.AuthTokens) authTokensResponse {
	accessTokenDuration := env.GetDuration("ACCESS_TOKEN_DURATION")
	refreshTokenDuration := env.GetDuration("REFRESH_TOKEN_DURATION")
	return authTokensResponse{
		AccessToken:             string(tokens.AccessToken),
		RefreshToken:            string(tokens.RefreshToken),
		AccessTokenExpiredInMs:  accessTokenDuration.Milliseconds(),
		RefreshTokenExpiredInMs: refreshTokenDuration.Milliseconds(),
	}
}

// newLoginResponse is a helper function that creates a loginResponse with the provided auth tokens and user.
func newLoginResponse(tokens *entities.AuthTokens, user *entities.User) loginResponse {
	return loginResponse{
		Tokens: newAuthTokensResponse(tokens),
		User:   newUserResponse(user),
	}
}

// refreshTokenResponse represents the structure of a response body for successful token refresh.
type refreshTokenResponse struct {
	Tokens authTokensResponse `json:"tokens"`
}

// newRefreshTokenResponse is a helper function that creates a refreshTokenResponse with the provided auth tokens.
func newRefreshTokenResponse(tokens *entities.AuthTokens) refreshTokenResponse {
	return refreshTokenResponse{
		Tokens: newAuthTokensResponse(tokens),
	}
}
