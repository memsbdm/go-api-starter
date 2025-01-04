package http

import (
	"go-starter/internal/adapters/validator"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"go-starter/internal/domain/ports"
	"net/http"
	"strconv"
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

// GetByID returns a user by id
func (uh *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		handleError(w, domain.ErrInvalidUserId)
		return
	}

	user, err := uh.svc.GetByID(ctx, userID)
	if err != nil {
		handleError(w, err)
		return
	}

	response := newUserResponse(user)
	handleSuccess(w, http.StatusOK, response)
}

// RegisterUserRequest represents the request body for creating a user
type RegisterUserRequest struct {
	Username string `json:"username" validate:"required"`
}

// Register registers a new user
func (uh *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payload RegisterUserRequest

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
