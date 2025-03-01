package handlers

import (
	"go-starter/internal/adapters/http/helpers"
	"go-starter/internal/adapters/http/responses"
	"go-starter/internal/adapters/validator"
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
//	@Success		200	{object}	responses.Response[responses.UserResponse]	"User displayed"
//	@Failure		401	{object}	responses.ErrorResponse	"Unauthorized error"
//	@Failure		500	{object}	responses.ErrorResponse	"Internal server error"
//	@Router			/v1/users/me [get]
//	@Security		BearerAuth
func (uh *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := helpers.GetUserIDFromContext(ctx)
	if err != nil {
		responses.HandleError(w, err)
		return
	}

	user, err := uh.svc.GetByID(ctx, userID)
	if err != nil {
		responses.HandleError(w, err)
		return
	}

	response := responses.NewUserResponse(user)
	responses.HandleSuccess(w, http.StatusOK, response)
}

// GetByID godoc
//
//	@Summary		Get a user
//	@Description	Get a user by id
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			uuid	path		string		true	"User ID" format(uuid)
//	@Success		200	{object}	responses.Response[responses.GetUserByIDResponse]	"User displayed"
//	@Failure		400	{object}	responses.ErrorResponse	"Incorrect User ID"
//	@Failure		404	{object}	responses.ErrorResponse	"Data not found error"
//	@Failure		500	{object}	responses.ErrorResponse	"Internal server error"
//	@Router			/v1/users/{uuid} [get]
//	@Security		BearerAuth
func (uh *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.PathValue("uuid")

	userID, err := entities.ParseUserID(id)
	if err != nil {
		responses.HandleError(w, err)
		return
	}

	user, err := uh.svc.GetByID(ctx, userID)
	if err != nil {
		responses.HandleError(w, err)
		return
	}

	response := responses.NewGetUserByIDResponse(user)
	responses.HandleSuccess(w, http.StatusOK, response)
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
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			updatePasswordRequest	body updatePasswordRequest true "Update user password request"
//	@Success		200	{object}	responses.EmptyResponse	"Success"
//	@Failure		400	{object}	responses.ErrorResponse	"Bad request error"
//	@Failure		401	{object}	responses.ErrorResponse	"Unauthorized error"
//	@Failure		422	{object}	responses.ErrorResponse	"Validation error"
//	@Failure		500	{object}	responses.ErrorResponse	"Internal server error"
//	@Router			/v1/users/me/password [patch]
//	@Security		BearerAuth
func (uh *UserHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var payload updatePasswordRequest
	if err := validator.ValidateRequest(w, r, &payload); err != nil {
		responses.HandleValidationError(w, err)
		return
	}

	updateUserParams := entities.UpdateUserParams{
		Password:             &payload.Password,
		PasswordConfirmation: &payload.PasswordConfirmation,
	}

	userID, err := helpers.GetUserIDFromContext(ctx)
	if err != nil {
		responses.HandleError(w, err)
		return
	}

	err = uh.svc.UpdatePassword(ctx, userID, updateUserParams)
	if err != nil {
		responses.HandleError(w, err)
		return
	}

	responses.HandleSuccess(w, http.StatusOK, nil)
}

// VerifyEmail godoc
//
//	@Summary		Verify user email
//	@Description	Verify user email
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			token	path		string		true	"Verification token"
//	@Success		200	{object}	responses.EmptyResponse	"User displayed"
//	@Failure		401	{object}	responses.ErrorResponse	"Unauthorized error / invalid token"
//	@Failure		409	{object}	responses.ErrorResponse	"Conflict error / already verified by another user"
//	@Failure		500	{object}	responses.ErrorResponse	"Internal server error"
//	@Router			/v1/users/me/email/verify/{token} [get]
func (uh *UserHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	token := r.PathValue("token")
	err := uh.svc.VerifyEmail(ctx, token)
	if err != nil {
		responses.HandleError(w, err)
		return
	}

	responses.HandleSuccess(w, http.StatusOK, nil)
}

// ResendEmailVerification godoc
//
//	@Summary		Resend user email verification
//	@Description	Resend user email verification
//	@Tags			Users
//	@Produce		json
//	@Success		200	{object}	responses.EmptyResponse	"Success"
//	@Failure		401	{object}	responses.ErrorResponse	"Unauthorized error"
//	@Failure		409	{object}	responses.ErrorResponse	"Conflict error / already verified"
//	@Failure		500	{object}	responses.ErrorResponse	"Internal server error"
//	@Router			/v1/users/me/email/verify/resend [post]
//	@Security		BearerAuth
func (uh *UserHandler) ResendEmailVerification(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := helpers.GetUserIDFromContext(ctx)
	if err != nil {
		responses.HandleError(w, err)
		return
	}

	err = uh.svc.ResendEmailVerification(ctx, userID)
	if err != nil {
		responses.HandleError(w, err)
		return
	}

	responses.HandleSuccess(w, http.StatusOK, nil)
}
