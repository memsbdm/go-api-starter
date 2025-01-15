package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"go-starter/internal/domain"
	"log/slog"
	"net/http"
)

var Validate *validator.Validate

// init initializes the validator with required struct validation enabled.
func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

// ErrInvalidJSON is returned when the JSON payload is invalid.
var ErrInvalidJSON = errors.New("invalid json")

// validationMessages holds custom error messages for specific validation failures.
var validationMessages = map[string]error{
	// Auth
	"loginRequest.Username.required":            domain.ErrUsernameRequired,
	"loginRequest.Password.required":            domain.ErrPasswordRequired,
	"refreshTokenRequest.RefreshToken.required": errors.New("refresh_token is required"),
	"refreshTokenRequest.RefreshToken.jwt":      domain.ErrInvalidToken,

	// Users
	"registerUserRequest.Username.required":               domain.ErrUsernameRequired,
	"registerUserRequest.Password.required":               domain.ErrPasswordRequired,
	"registerUserRequest.Password.min":                    domain.ErrPasswordTooShort,
	"updatePasswordRequest.Password.required":             domain.ErrPasswordRequired,
	"updatePasswordRequest.Password.min":                  domain.ErrPasswordTooShort,
	"updatePasswordRequest.Password.eqfield":              domain.ErrPasswordsNotMatch,
	"updatePasswordRequest.PasswordConfirmation.required": domain.ErrPasswordConfirmationRequired,
}

// ValidateRequest takes a payload from an HTTP request and verifies it.
func ValidateRequest(w http.ResponseWriter, r *http.Request, payload interface{}) []error {
	maxBytes := 1_048_576 // 1mb
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	defer func() {
		err := r.Body.Close()
		if err != nil {
			errMsg := fmt.Sprintf("failed to close body: %v", err)
			slog.Error(errMsg)
		}
	}()
	decoder := json.NewDecoder(r.Body)

	// Validate JSON format
	if err := decoder.Decode(&payload); err != nil {
		return []error{ErrInvalidJSON}
	}

	// Validate payload
	var errs []error
	if err := Validate.Struct(payload); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			field := fmt.Sprintf("%s.%s", err.StructNamespace(), err.Tag())
			message, ok := validationMessages[field]
			if !ok {
				message = errors.New(fmt.Sprintf("Validation failed on field '%s' for condition '%s'", err.Field(), err.Tag()))
			}
			errs = append(errs, message)
		}
	}
	return errs
}
