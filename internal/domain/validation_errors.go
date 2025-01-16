package domain

import (
	"errors"
	"fmt"
)

const (
	UsernameMinLength = 4
	UsernameMaxLength = 15
	PasswordMinLength = 8
)

// Returned validation errors
var (
	// Required

	// ErrUsernameRequired represents an error when the username is required but not provided.
	ErrUsernameRequired = errors.New("username is required")
	// ErrPasswordRequired represents an error when the password is required but not provided.
	ErrPasswordRequired = errors.New("password is required")
	// ErrPasswordConfirmationRequired represents an error when the password confirmation is required but not provided.
	ErrPasswordConfirmationRequired = errors.New("password confirmation required")

	// Other validation errors

	// ErrPasswordsNotMatch represents an error when the provided passwords do not match.
	ErrPasswordsNotMatch = errors.New("passwords does not match")
	// ErrPasswordTooShort represents an error when the password is too short, less than the minimum required length.
	ErrPasswordTooShort = errors.New(fmt.Sprintf("password is too short, it should be at least %d characters", PasswordMinLength))
	// ErrUsernameTooShort represents an error when the username is too short, less than the minimum required length.
	ErrUsernameTooShort = errors.New(fmt.Sprintf("username is too short, it should be at least %d characters", UsernameMinLength))
	// ErrUsernameTooLong represents an error when the username is too long, greater than the minimum required length.
	ErrUsernameTooLong = errors.New(fmt.Sprintf("username is too long, it should be at most %d characters", UsernameMaxLength))
	// ErrUsernameInvalid represents an error when the username is invalid, not respecting the regex pattern.
	ErrUsernameInvalid = errors.New("username can only contain alphanumeric characters and underscore")
)
