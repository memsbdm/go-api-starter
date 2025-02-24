package responses

import (
	"go-starter/internal/domain/entities"
	"time"
)

// UserResponse represents the structure of a response body containing user information.
type UserResponse struct {
	ID              string    `json:"id" example:"6b947a32-8919-4974-9ef3-048a556b0b75"`
	CreatedAt       time.Time `json:"created_at" example:"2024-08-15T16:23:33.455225Z"`
	UpdatedAt       time.Time `json:"updated_at" example:"2025-01-15T14:29:33.455225Z"`
	Name            string    `json:"name" example:"John Doe"`
	Username        string    `json:"username" example:"john"`
	Email           string    `json:"email" example:"john@example.com"`
	IsEmailVerified bool      `json:"is_email_verified" example:"true"`
}

// NewUserResponse is a helper function that creates a UserResponse from a user entity.
func NewUserResponse(user *entities.User) UserResponse {
	return UserResponse{
		ID:              user.ID.String(),
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
		Name:            user.Name,
		Username:        user.Username,
		Email:           user.Email,
		IsEmailVerified: user.IsEmailVerified,
	}
}

// GetUserByIDResponse represents the structure of a response body containing user information.
type GetUserByIDResponse struct {
	ID       string `json:"id" example:"6b947a32-8919-4974-9ef3-048a556b0b75"`
	Name     string `json:"name" example:"John Doe"`
	Username string `json:"username" example:"john"`
}

// NewGetUserByIDResponse is a helper function that creates a UserResponse from a user entity.
func NewGetUserByIDResponse(user *entities.User) GetUserByIDResponse {
	return GetUserByIDResponse{
		ID:       user.ID.String(),
		Name:     user.Name,
		Username: user.Username,
	}
}
