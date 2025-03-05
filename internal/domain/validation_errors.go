package domain

import (
	"errors"
	"fmt"
)

// Validation constants
const (
	NameMaxLength     = 50
	UsernameMinLength = 4
	UsernameMaxLength = 15
	PasswordMinLength = 8
	EmailMaxLength    = 254
)

// Required validation errors
var (
	// ErrUsernameRequired represents an error when the username is required but not provided.
	ErrUsernameRequired = errors.New("username is required")
	// ErrPasswordRequired represents an error when the password is required but not provided.
	ErrPasswordRequired = errors.New("password is required")
	// ErrPasswordConfirmationRequired represents an error when the password confirmation is required but not provided.
	ErrPasswordConfirmationRequired = errors.New("password confirmation required")
	// ErrNameRequired represents an error when name is required but not provided.
	ErrNameRequired = errors.New("name is required")
	// ErrEmailRequired represents an error when email is required but not provided.
	ErrEmailRequired = errors.New("email is required")
)

// Other validation errors
var (
	// ErrPasswordsNotMatch represents an error when the provided passwords do not match.
	ErrPasswordsNotMatch = errors.New("passwords does not match")
	// ErrPasswordTooShort represents an error when the password is too short, less than the minimum required length.
	ErrPasswordTooShort = fmt.Errorf("password is too short, it should be at least %d characters", PasswordMinLength)
	// ErrUsernameTooShort represents an error when the username is too short, less than the minimum required length.
	ErrUsernameTooShort = fmt.Errorf("username is too short, it should be at least %d characters", UsernameMinLength)
	// ErrUsernameConflict represents a conflict error when trying to create a user with an existing username.
	ErrUsernameConflict = errors.New("username already taken")
	// ErrUsernameTooLong represents an error when the username is too long, greater than the minimum required length.
	ErrUsernameTooLong = fmt.Errorf("username is too long, it should be at most %d characters", UsernameMaxLength)
	// ErrUsernameInvalid represents an error when the username is invalid, not respecting the regex pattern.
	ErrUsernameInvalid = errors.New("username can only contain alphanumeric characters and underscore")
	// ErrNameTooLong represents an error when the name is too long, greater than the minimum required length.
	ErrNameTooLong = fmt.Errorf("name is too long, it should be at most %d characters", NameMaxLength)
	// ErrEmailInvalid represents an error when the email is invalid.
	ErrEmailInvalid = errors.New("email is invalid")
	// ErrEmailConflict represents a conflict error when trying to create a user with an existing email.
	ErrEmailConflict = errors.New("email already taken")
)
