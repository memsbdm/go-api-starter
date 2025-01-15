package services

import (
	"context"
	"errors"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"go-starter/internal/domain/ports"
	"go-starter/internal/domain/utils"
)

// AuthService implements ports.AuthService interface.
type AuthService struct {
	userSvc  ports.UserService
	tokenSvc ports.TokenService
}

// NewAuthService creates a new instance of AuthService.
func NewAuthService(userSvc ports.UserService, tokenSvc ports.TokenService) *AuthService {
	return &AuthService{
		userSvc:  userSvc,
		tokenSvc: tokenSvc,
	}
}

// Login authenticates a user.
// Returns an access token and a refresh token upon successful authentication,
// or an error if the login fails (e.g., due to incorrect credentials).
func (as *AuthService) Login(ctx context.Context, username, password string) (string, string, error) {
	user, err := as.userSvc.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", "", domain.ErrInvalidCredentials
		}
		return "", "", domain.ErrInternal
	}

	err = utils.ComparePassword(password, user.Password)
	if err != nil {
		return "", "", domain.ErrInvalidCredentials
	}

	accessToken, err := as.tokenSvc.GenerateAccessToken(user)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := as.tokenSvc.GenerateRefreshToken(ctx, user.ID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// Register registers a new user in the system.
// Returns the created user entity and an error if the registration fails
// (e.g., due to username already existing or validation issues).
func (as *AuthService) Register(ctx context.Context, user *entities.User) (*entities.User, error) {
	return as.userSvc.Register(ctx, user)
}

// Refresh generates new access and refresh tokens using the previous refresh token.
// Returns the new access token, new refresh token, and an error if the refresh fails
// (e.g., if the previous refresh token is invalid or expired).
func (as *AuthService) Refresh(ctx context.Context, previousRefreshToken string) (string, string, error) {
	claims, err := as.tokenSvc.ValidateAndParseRefreshToken(ctx, previousRefreshToken)
	if err != nil {
		return "", "", err
	}

	user, err := as.userSvc.GetByID(ctx, claims.Subject)
	if err != nil {
		return "", "", err
	}

	accessToken, err := as.tokenSvc.GenerateAccessToken(user)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := as.tokenSvc.GenerateRefreshToken(ctx, claims.Subject)
	if err != nil {
		return "", "", err
	}

	err = as.tokenSvc.RevokeRefreshToken(ctx, previousRefreshToken)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// Logout invalidates the specified refresh token, effectively logging the user out.
// Returns an error if the logout operation fails (e.g., if the refresh token is not found).
func (as *AuthService) Logout(ctx context.Context, refreshToken string) error {
	return as.tokenSvc.RevokeRefreshToken(ctx, refreshToken)
}
