package domain

import "errors"

// Errors returned in responses.

// Generic errors.
var (
	// ErrInternal represents an internal error.
	ErrInternal = errors.New("internal error")
	// ErrForbidden represents a forbidden error.
	ErrForbidden = errors.New("forbidden")
	// ErrUnauthorized represents an unauthorized error.
	ErrUnauthorized = errors.New("unauthorized")
	// ErrBadRequest represents a bad request error.
	ErrBadRequest = errors.New("bad request")
)

// Mailer errors.
var (
	// ErrMailer represents a mailer error.
	ErrMailer = errors.New("mailer error")
)

// File upload errors.
var (
	// ErrFileUpload represents a file upload error.
	ErrFileUpload = errors.New("file upload error")
	// ErrFileTooLarge represents a file too large error.
	ErrFileTooLarge = errors.New("file too large")
	// ErrMissingBoundary represents a missing boundary error.
	ErrMissingBoundary = errors.New("missing boundary")
	// ErrInvalidMultipartForm represents an invalid multipart form error.
	ErrInvalidMultipartForm = errors.New("invalid multipart form")
	// ErrInvalidFileType represents an invalid file type error.
	ErrInvalidFileType = errors.New("invalid file type")
)

// Auth errors.
var (
	// ErrInvalidToken represents an error for an invalid token.
	ErrInvalidToken = errors.New("invalid token")
	// ErrInvalidCredentials represents an error for invalid login credentials.
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// User errors.
var (
	// ErrInvalidUserId represents an error for an invalid user ID format.
	ErrInvalidUserId = errors.New("invalid user id")
	// ErrUserNotFound represents an error when a user is not found.
	ErrUserNotFound = errors.New("user not found")
	// ErrEmailAlreadyVerified represents an error when a user's email is already verified.
	ErrEmailAlreadyVerified = errors.New("email already verified")
)

// Errors not returned in responses.
var (
	// ErrCacheNotFound represents an error for an empty cache value for a given key.
	ErrCacheNotFound = errors.New("cache not found")
)
