package http

import (
	"github.com/google/uuid"
	"go-starter/internal/adapters/validator"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"go-starter/internal/domain/ports"
	"net/http"
)

// UserHandler represents the HTTP handler for user-related requests.
type UserHandler struct {
	svc ports.UserService
}

// NewUserHandler creates and returns a new UserHandler instance.
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

	claims, err := extractAccessTokenClaims(ctx)
	if err != nil {
		handleError(w, domain.ErrInternal)
		return
	}

	user, err := uh.svc.GetByID(ctx, claims.Subject)
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
//	@Success		200	{object}	response[getUserByIDResponse]	"User displayed"
//	@Failure		400	{object}	errorResponse	"Incorrect User ID"
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

	user, err := uh.svc.GetByID(ctx, entities.UserID(userUUID))
	if err != nil {
		handleError(w, err)
		return
	}

	response := newGetUserByIDResponse(user)
	handleSuccess(w, http.StatusOK, response)
}

// updatePasswordRequest represents the structure of the request body used for updating a user password.
type updatePasswordRequest struct {
	Password             string `json:"password" validate:"required,min=8,eqfield=PasswordConfirmation" example:"secret123"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required" example:"secret123"`
}

// UpdatePassword godoc
//
//	@Summary		Update user password
//	@Description	Update user password
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			updatePasswordRequest	body updatePasswordRequest true "Update user password request"
//	@Success		200	{object}	emptyResponse	"Success"
//	@Failure		401	{object}	errorResponse	"Unauthorized error"
//	@Failure		422	{object}	errorResponse	"Validation error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/v1/users/password [patch]
//	@Security		BearerAuth
func (uh *UserHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var payload updatePasswordRequest
	if err := validator.ValidateRequest(w, r, &payload); err != nil {
		handleValidationError(w, err)
		return
	}

	updateUserParams := entities.UpdateUserParams{
		Password:             &payload.Password,
		PasswordConfirmation: &payload.PasswordConfirmation,
	}

	claims, err := extractAccessTokenClaims(ctx)
	if err != nil {
		handleError(w, domain.ErrInternal)
		return
	}

	err = uh.svc.UpdatePassword(ctx, claims.Subject, updateUserParams)
	if err != nil {
		handleError(w, err)
		return
	}

	handleSuccess(w, http.StatusOK, nil)
}
