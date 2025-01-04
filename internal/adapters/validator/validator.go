package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

var ErrInvalidJSON = errors.New("invalid json")

var validationMessages = map[string]error{
	"registerUserRequest.Username.required": errors.New("username is required"),
}

// ValidateRequest takes a payload from an HTTP request and verifies it. It can return an error for an invalid JSON or
// for custom payload validation errors.
func ValidateRequest(w http.ResponseWriter, r *http.Request, payload interface{}) []error {
	maxBytes := 1_048_576 // 1mb
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	defer func() {
		err := r.Body.Close()
		if err != nil {
			slog.Error("error closing body: %v", err)
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
