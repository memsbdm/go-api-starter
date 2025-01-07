package http

import (
	"github.com/google/uuid"
	"go-starter/internal/adapters/validator"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"go-starter/internal/domain/ports"
	"net/http"
)

// UserHandler represents the HTTP handler for user-related requests
type UserHandler struct {
	svc ports.UserService
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(svc ports.UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}

// GetByID godoc
//
//	@Summary		Get a user
//	@Description	Get a user by id
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int		true	"User ID"
//	@Success		200	{object}	userResponse	"User displayed"
//	@Failure		400	{object}	errorResponse	"Incorrect User ID"
//	@Failure		404	{object}	errorResponse	"Data not found error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/v1/users/{id} [get]
func (uh *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	userID, err := uuid.Parse(id)

	if err != nil {
		handleError(w, domain.ErrInvalidUserId)
		return
	}

	user, err := uh.svc.GetByID(ctx, entities.UserID(userID))
	if err != nil {
		handleError(w, err)
		return
	}

	response := newUserResponse(user)
	handleSuccess(w, http.StatusOK, response)
}

// registerUserRequest represents the request body for creating a user
type registerUserRequest struct {
	Username string `json:"username" validate:"required" example:"john"`
	Password string `json:"password" validate:"required" example:"secret"`
}

// Register godoc
//
//	@Summary		Register a new user
//	@Description	Create a new user account
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			registerUserRequest	body registerUserRequest true "Register request"
//	@Success		200	{object}	userResponse	"User created"
//	@Failure		409	{object}	errorResponse	"Duplication error"
//	@Failure		422	{object}	errorResponse	"Validation error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/v1/users [post]
func (uh *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payload registerUserRequest

	if err := validator.ValidateRequest(w, r, &payload); err != nil {
		handleValidationError(w, err)
		return
	}

	user := entities.User{
		Username: payload.Username,
	}

	created, err := uh.svc.Register(ctx, &user)
	if err != nil {
		handleError(w, err)
		return
	}

	response := newUserResponse(created)
	handleSuccess(w, http.StatusCreated, response)
}
