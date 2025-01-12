package ports

import (
	"context"
)

// AuthService is an interface for interacting with authentication
type AuthService interface {
	// Login authenticates a user.
	// Returns an access token, a refresh token or an error if login fails
	Login(ctx context.Context, username, password string) (string, string, error)
}
