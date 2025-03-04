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
	RoleID          int       `json:"role_id" example:"1"`
	AvatarURL       string    `json:"avatar_url" example:"https://example.com/avatar.jpg"`
}

// NewUserResponse is a helper function that creates a UserResponse from a user entity.
func NewUserResponse(user *entities.User) UserResponse {
	var avatarURL string
	if user.AvatarURL != nil {
		avatarURL = *user.AvatarURL
	}

	return UserResponse{
		ID:              user.ID.String(),
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
		Name:            user.Name,
		Username:        user.Username,
		Email:           user.Email,
		IsEmailVerified: user.IsEmailVerified,
		RoleID:          user.RoleID.Int(),
		AvatarURL:       avatarURL,
	}
}

// GetUserByIDResponse represents the structure of a response body containing user information.
type GetUserByIDResponse struct {
	ID        string `json:"id" example:"6b947a32-8919-4974-9ef3-048a556b0b75"`
	Name      string `json:"name" example:"John Doe"`
	Username  string `json:"username" example:"john"`
	AvatarURL string `json:"avatar_url" example:"https://example.com/avatar.jpg"`
}

// NewGetUserByIDResponse is a helper function that creates a UserResponse from a user entity.
func NewGetUserByIDResponse(user *entities.User) GetUserByIDResponse {
	var avatarURL string
	if user.AvatarURL != nil {
		avatarURL = *user.AvatarURL
	}

	return GetUserByIDResponse{
		ID:        user.ID.String(),
		Name:      user.Name,
		Username:  user.Username,
		AvatarURL: avatarURL,
	}
}

// UploadAvatarResponse represents the structure of a response body containing the avatar URL.
type UploadAvatarResponse struct {
	AvatarURL string `json:"avatar_url" example:"https://example.com/avatar.jpg"`
}

// NewUploadAvatarResponse is a helper function that creates a UploadAvatarResponse from an avatar URL.
func NewUploadAvatarResponse(avatarURL string) UploadAvatarResponse {
	return UploadAvatarResponse{AvatarURL: avatarURL}
}
