package http

import (
	"go-starter/internal/adapters/validator"
	"go-starter/internal/domain/ports"
	"net/http"
)

// AuthHandler represents the HTTP handler for token-related requests
type AuthHandler struct {
	svc ports.AuthService
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(svc ports.AuthService) *AuthHandler {
	return &AuthHandler{
		svc: svc,
	}
}

// loginRequest represents the request body for login a user
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
//	@Param			LoginRequest	body loginRequest true "Login request"
//	@Success		200	{object}	loginResponse	"User logged in"
//	@Failure		401	{object}	errorResponse	"Unauthorized / credentials error"
//	@Failure		403	{object}	errorResponse	"Forbidden error"
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
