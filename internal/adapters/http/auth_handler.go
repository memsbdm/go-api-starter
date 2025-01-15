package http

import (
	"go-starter/internal/adapters/validator"
	"go-starter/internal/domain/entities"
	"go-starter/internal/domain/ports"
	"net/http"
)

// AuthHandler is responsible for handling HTTP requests related to authentication operations.
// It acts as a bridge between the HTTP layer and the authentication service.
type AuthHandler struct {
	svc ports.AuthService
}

// NewAuthHandler initializes and returns a new AuthHandler instance.
// It accepts an implementation of the AuthService interface to handle authentication logic.
func NewAuthHandler(svc ports.AuthService) *AuthHandler {
	return &AuthHandler{
		svc: svc,
	}
}

// loginRequest represents the structure of the request body used for logging in a user.
type loginRequest struct {
	Username string `json:"username" validate:"required" example:"john"`
	Password string `json:"password" validate:"required" example:"secret123"`
}

// Login godoc
//
//	@Summary		Login a user
//	@Description	Authenticate a user account
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			loginRequest	body loginRequest true "Login request"
//	@Success		200	{object}	response[loginResponse]	"Access and refresh tokens"
//	@Failure		401	{object}	errorResponse	"Unauthorized / credentials error"
//	@Failure		403	{object}	errorResponse	"Forbidden error"
//	@Failure		422	{object}	errorResponse	"Validation error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/v1/auth/login [post]
func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payload loginRequest

	if err := validator.ValidateRequest(w, r, &payload); err != nil {
		handleValidationError(w, err)
		return
	}
	accessToken, refreshToken, err := ah.svc.Login(ctx, payload.Username, payload.Password)
	if err != nil {
		handleError(w, err)
		return
	}

	response := newLoginResponse(accessToken, refreshToken)
	handleSuccess(w, http.StatusOK, response)
}

// registerUserRequest represents the structure of the request body used for registering a new user.
type registerUserRequest struct {
	Username string `json:"username" validate:"required" example:"john"`
	Password string `json:"password" validate:"required,min=8" example:"secret123"`
}

// Register godoc
//
//	@Summary		Register a new user
//	@Description	Create a new user account
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			registerUserRequest	body registerUserRequest true "Register request"
//	@Success		200	{object}	response[userResponse]	"Created user"
//	@Failure		403	{object}	errorResponse	"Forbidden error"
//	@Failure		409	{object}	errorResponse	"Duplication error"
//	@Failure		422	{object}	errorResponse	"Validation error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/v1/auth/register [post]
func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payload registerUserRequest

	if err := validator.ValidateRequest(w, r, &payload); err != nil {
		handleValidationError(w, err)
		return
	}

	user := entities.User{
		Username: payload.Username,
		Password: payload.Password,
	}

	created, err := ah.svc.Register(ctx, &user)
	if err != nil {
		handleError(w, err)
		return
	}

	response := newUserResponse(created)
	handleSuccess(w, http.StatusCreated, response)
}

// refreshTokenRequest represents the structure of the request body used for refreshing token or revoke existing one.
type refreshTokenRequest struct {
	RefreshToken string `validate:"required,jwt" example:"eyJhbGci..."`
}

// Refresh godoc
//
//	@Summary		Generate a new access token and refresh token
//	@Description	Generate a new access token and refresh token
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			refreshTokenRequest	body refreshTokenRequest true "Refresh token request"
//	@Success		200	{object}	response[loginResponse]	"Access and refresh tokens"
//	@Failure		401	{object}	errorResponse	"Unauthorized error"
//	@Failure		422	{object}	errorResponse	"Validation error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/v1/auth/refresh [post]
//	@Security		BearerAuth
func (ah *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payload refreshTokenRequest

	if err := validator.ValidateRequest(w, r, &payload); err != nil {
		handleValidationError(w, err)
		return
	}

	accessToken, refreshToken, err := ah.svc.Refresh(ctx, payload.RefreshToken)
	if err != nil {
		handleError(w, err)
		return
	}

	response := newLoginResponse(accessToken, refreshToken)
	handleSuccess(w, http.StatusOK, response)
}

// Logout godoc
//
//	@Summary		Logout an authenticated user
//	@Description	Logout an authenticated user
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			refreshTokenRequest	body refreshTokenRequest true "Refresh token request"
//	@Success		200	{object}	emptyResponse	"Success"
//	@Failure		401	{object}	errorResponse	"Unauthorized error"
//	@Failure		422	{object}	errorResponse	"Validation error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/v1/auth/logout [delete]
//	@Security		BearerAuth
func (ah *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var payload refreshTokenRequest

	if err := validator.ValidateRequest(w, r, &payload); err != nil {
		handleValidationError(w, err)
		return
	}

	err := ah.svc.Logout(ctx, payload.RefreshToken)
	if err != nil {
		handleError(w, err)
		return
	}

	handleSuccess(w, http.StatusOK, nil)
}
