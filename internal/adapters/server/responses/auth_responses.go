package responses

import (
	"go-starter/internal/domain/entities"
)

// LoginResponse represents the structure of a response body for a successful authentication.
type LoginResponse struct {
	AccessToken string       `json:"access_token"`
	User        UserResponse `json:"user"`
}

// NewLoginResponse is a helper function that creates a LoginResponse.
func NewLoginResponse(token string, user *entities.User) LoginResponse {
	return LoginResponse{
		AccessToken: token,
		User:        NewUserResponse(user),
	}
}
