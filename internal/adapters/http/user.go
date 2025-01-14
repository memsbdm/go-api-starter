package http

import (
	"github.com/google/uuid"
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

// Me godoc
//
//	@Summary		Get authenticated user information
//	@Description	Get information of logged-in user
//	@Tags			Users
//	@Produce		json
//	@Success		200	{object}	response[userResponse]	"User displayed"
//	@Failure		401	{object}	errorResponse	"Unauthorized error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/v1/users/me [get]
//	@Security		BearerAuth
func (uh *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	authPayload, err := getAuthPayload(ctx, authorizationPayloadKey)
	if err != nil {
		handleError(w, domain.ErrInternal)
		return
	}

	user, err := uh.svc.GetByID(ctx, entities.UserID(authPayload.UserID))
	if err != nil {
		handleError(w, err)
		return
	}

	response := newUserResponse(user)
	handleSuccess(w, http.StatusOK, response)
}

// GetByID godoc
//
//	@Summary		Get a user
//	@Description	Get a user by id
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			uuid	path		string		true	"User ID" format(uuid)
//	@Success		200	{object}	response[userResponse]	"User displayed"
//	@Failure		400	{object}	errorResponse	"Incorrect User ID"
//	@Failure		401	{object}	errorResponse	"Unauthorized error"
//	@Failure		403	{object}	errorResponse	"Forbidden error"
//	@Failure		404	{object}	errorResponse	"Data not found error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/v1/users/{uuid} [get]
//	@Security		BearerAuth
func (uh *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("uuid")

	userUUID, err := uuid.Parse(id)
	if err != nil {
		handleError(w, domain.ErrInvalidUserId)
		return
	}

	authPayload, err := getAuthPayload(ctx, authorizationPayloadKey)
	if err != nil {
		handleError(w, domain.ErrInternal)
		return
	}

	userID := userUUID
	if authPayload.UserID != userID {
		handleError(w, domain.ErrForbidden)
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
