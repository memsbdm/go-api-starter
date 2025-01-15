package domain

import (
	"errors"
	"fmt"
)

const PasswordMinLength = 8

// Returned validation errors
var (
	// ErrPasswordRequired represents an error when the password is required but not provided.
	ErrPasswordRequired = errors.New("password is required")
	// ErrPasswordConfirmationRequired represents an error when the password confirmation is required but not provided.
	ErrPasswordConfirmationRequired = errors.New("password confirmation required")
	// ErrPasswordsNotMatch represents an error when the provided passwords do not match.
	ErrPasswordsNotMatch = errors.New("passwords does not match")
	// ErrPasswordTooShort represents an error when the password is too short, less than the minimum required length.
	ErrPasswordTooShort = errors.New(fmt.Sprintf("password is too short, it should be at least %d characters", PasswordMinLength))
)
