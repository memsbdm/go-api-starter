package handlers

import (
	"go-starter/internal/adapters/http/helpers"
	"go-starter/internal/adapters/http/responses"
	"go-starter/internal/adapters/validator"
	"go-starter/internal/domain/entities"
	"go-starter/internal/domain/ports"
	"net/http"
	"strings"
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
//	@Success		200	{object}	responses.Response[responses.LoginResponse]	"Login response"
//	@Failure		400	{object}	responses.ErrorResponse	"Bad request error"
//	@Failure		401	{object}	responses.ErrorResponse	"Unauthorized / credentials error"
//	@Failure		403	{object}	responses.ErrorResponse	"Forbidden error"
//	@Failure		422	{object}	responses.ErrorResponse	"Validation error"
//	@Failure		500	{object}	responses.ErrorResponse	"Internal server error"
//	@Router			/v1/auth/login [post]
func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payload loginRequest
	if err := validator.ValidateRequest(w, r, &payload); err != nil {
		responses.HandleValidationError(w, err)
		return
	}

	payload.Username = strings.TrimSpace(payload.Username)
	user, accessToken, err := ah.svc.Login(ctx, payload.Username, payload.Password)
	if err != nil {
		responses.HandleError(w, err)
		return
	}

	response := responses.NewLoginResponse(accessToken, user)
	responses.HandleSuccess(w, http.StatusOK, response)
}

// registerRequest represents the structure of the request body used for registering a new user.
type registerRequest struct {
	Name     string `json:"name" validate:"required,min=1,max=50" example:"John Doe"`
	Username string `json:"username" validate:"required,min=4,max=15" example:"john"`
	Password string `json:"password" validate:"required,min=8" example:"secret123"`
	Email    string `json:"email" validate:"required,email" example:"john@example.com"`
}

// Register godoc
//
//	@Summary		Register a new user
//	@Description	Create a new user account
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			registerRequest	body registerRequest true "Register request"
//	@Success		200	{object}	responses.Response[responses.LoginResponse]	"Created user"
//	@Failure		400	{object}	responses.ErrorResponse	"Bad request error"
//	@Failure		403	{object}	responses.ErrorResponse	"Forbidden error"
//	@Failure		409	{object}	responses.ErrorResponse	"Duplication error"
//	@Failure		422	{object}	responses.ErrorResponse	"Validation error"
//	@Failure		500	{object}	responses.ErrorResponse	"Internal server error"
//	@Router			/v1/auth/register [post]
func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payload registerRequest

	if err := validator.ValidateRequest(w, r, &payload); err != nil {
		responses.HandleValidationError(w, err)
		return
	}

	user := &entities.User{
		Name:     payload.Name,
		Username: payload.Username,
		Password: payload.Password,
		Email:    payload.Email,
	}

	created, err := ah.svc.Register(ctx, user)
	if err != nil {
		responses.HandleError(w, err)
		return
	}

	authenticatedUser, authTokens, err := ah.svc.Login(ctx, created.Username, payload.Password)
	if err != nil {
		responses.HandleError(w, err)
		return
	}

	response := responses.NewLoginResponse(authTokens, authenticatedUser)
	responses.HandleSuccess(w, http.StatusCreated, response)
}

// Logout godoc
//
//	@Summary		Logout an authenticated user
//	@Description	Logout an authenticated user
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	responses.EmptyResponse	"Success"
//	@Failure		401	{object}	responses.ErrorResponse	"Unauthorized error"
//	@Failure		403	{object}	responses.ErrorResponse	"Forbidden error"
//	@Failure		500	{object}	responses.ErrorResponse	"Internal server error"
//	@Router			/v1/auth/logout [delete]
//	@Security		BearerAuth
func (ah *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	accessToken, err := helpers.ExtractTokenFromHeader(r)
	if err != nil {
		responses.HandleError(w, err)
		return
	}

	err = ah.svc.Logout(ctx, accessToken)
	if err != nil {
		responses.HandleError(w, err)
		return
	}

	responses.HandleSuccess(w, http.StatusOK, nil)
}
