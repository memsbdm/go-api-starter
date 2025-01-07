package domain

import "errors"

var (
	ErrInternal                  = errors.New("internal error")
	ErrInvalidUserId             = errors.New("invalid user id")
	ErrUserNotFound              = errors.New("user not found")
	ErrUserUsernameAlreadyExists = errors.New("user username already exists")
	ErrCannotParseUUID           = errors.New("cannot parse uuid")
)
