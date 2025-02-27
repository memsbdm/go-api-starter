package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-starter/internal/domain"
	"log/slog"
	"net/http"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

const (
	// MaxRequestSize defines the maximum allowed size for request bodies (1MB)
	MaxRequestSize = 1 << 20
)

var (
	Validate *validator.Validate
	once     sync.Once
)

// init initializes the validator with required struct validation enabled.
func init() {
	once.Do(func() {
		Validate = validator.New(validator.WithRequiredStructEnabled())
		if err := Validate.RegisterValidation("notblank", notBlank); err != nil {
			slog.Error("failed to register notblank validation", "error", err)
		}
	})
}

// notBlank validates that the string length is greater than 0 after trimming whitespace.
func notBlank(fl validator.FieldLevel) bool {
	return len(strings.TrimSpace(fl.Field().String())) > 0
}

// ErrInvalidJSON is returned when the JSON payload is invalid.
var ErrInvalidJSON = errors.New("invalid json")

// validationMessages holds custom error messages for specific validation failures.
var validationMessages = map[string]error{
	// Auth
	"loginRequest.Username.notblank":    domain.ErrUsernameRequired,
	"loginRequest.Password.required":    domain.ErrPasswordRequired,
	"registerRequest.Name.notblank":     domain.ErrNameRequired,
	"registerRequest.Name.max":          domain.ErrNameTooLong,
	"registerRequest.Username.notblank": domain.ErrUsernameRequired,
	"registerRequest.Username.min":      domain.ErrUsernameTooShort,
	"registerRequest.Username.max":      domain.ErrUsernameTooLong,
	"registerRequest.Password.required": domain.ErrPasswordRequired,
	"registerRequest.Password.min":      domain.ErrPasswordTooShort,
	"registerRequest.Email.required":    domain.ErrEmailRequired,
	"registerRequest.Email.email":       domain.ErrEmailInvalid,

	// Users
	"updatePasswordRequest.Password.required":             domain.ErrPasswordRequired,
	"updatePasswordRequest.Password.min":                  domain.ErrPasswordTooShort,
	"updatePasswordRequest.Password.eqfield":              domain.ErrPasswordsNotMatch,
	"updatePasswordRequest.PasswordConfirmation.required": domain.ErrPasswordConfirmationRequired,
}

// ValidateRequest takes a payload from an HTTP request and verifies it.
func ValidateRequest(w http.ResponseWriter, r *http.Request, payload interface{}) []error {
	r.Body = http.MaxBytesReader(w, r.Body, MaxRequestSize)
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
				message = fmt.Errorf("validation failed on field '%s' for condition '%s'", err.Field(), err.Tag())
			}
			errs = append(errs, message)
		}
	}
	return errs
}
