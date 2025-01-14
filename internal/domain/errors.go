package domain

import "errors"

// Errors returned in responses
var (
	// ErrInternal represents an internal error
	ErrInternal = errors.New("internal error")
	// ErrForbidden represents a forbidden error
	ErrForbidden = errors.New("forbidden")
	// ErrUnauthorized represents an unauthorized error
	ErrUnauthorized = errors.New("unauthorized")
	// ErrBadRequest represents a bad request error
	ErrBadRequest = errors.New("bad request")

	// ErrInvalidToken represents an invalid token error
	ErrInvalidToken = errors.New("invalid token")

	// ErrInvalidCredentials represents a login error
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrInvalidUserId represents a wrong user id format error
	ErrInvalidUserId = errors.New("invalid user id")
	// ErrUserNotFound represents a user not found error
	ErrUserNotFound = errors.New("user not found")
	// ErrUserUsernameAlreadyExists represents a conflict during a user creation on the username field
	ErrUserUsernameAlreadyExists = errors.New("user username already exists")
)

// Errors not returned in responses
var (
	// ErrCacheNotFound represents an empty cache value for a given key
	ErrCacheNotFound = errors.New("cache not found")
	// ErrTokenPayloadNotFound represents an authorization token payload not found
	ErrTokenPayloadNotFound = errors.New("token payload not found, maybe missing token middleware")
)
