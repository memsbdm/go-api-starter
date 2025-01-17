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
// Returns auth tokens upon successful authentication,
// or an error if the login fails (e.g., due to incorrect credentials).
func (as *AuthService) Login(ctx context.Context, username, password string) (*entities.User, *entities.AuthTokens, error) {
	user, err := as.userSvc.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, nil, domain.ErrInvalidCredentials
		}
		return nil, nil, domain.ErrInternal
	}

	err = utils.ComparePassword(password, user.Password)
	if err != nil {
		return nil, nil, domain.ErrInvalidCredentials
	}

	accessToken, err := as.tokenSvc.GenerateAccessToken(user)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := as.tokenSvc.GenerateRefreshToken(ctx, user.ID)
	if err != nil {
		return nil, nil, err
	}

	authTokens := &entities.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return user, authTokens, nil
}

// Register registers a new user in the system.
// Returns the created user entity and an error if the registration fails
// (e.g., due to username already existing or validation issues).
func (as *AuthService) Register(ctx context.Context, user *entities.User) (*entities.User, error) {
	return as.userSvc.Register(ctx, user)
}

// Refresh generates new access and refresh tokens using the previous refresh token.
// Returns the new auth tokens, and an error if the refresh fails
// (e.g., if the previous refresh token is invalid or expired).
func (as *AuthService) Refresh(ctx context.Context, previousRefreshToken string) (*entities.AuthTokens, error) {
	if previousRefreshToken == "" {
		return nil, domain.ErrRefreshTokenRequired
	}
	claims, err := as.tokenSvc.ValidateAndParseRefreshToken(ctx, previousRefreshToken)
	if err != nil {
		return nil, err
	}

	user, err := as.userSvc.GetByID(ctx, claims.Subject)
	if err != nil {
		return nil, err
	}

	accessToken, err := as.tokenSvc.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := as.tokenSvc.GenerateRefreshToken(ctx, claims.Subject)
	if err != nil {
		return nil, err
	}

	err = as.tokenSvc.RevokeRefreshToken(ctx, previousRefreshToken)
	if err != nil {
		return nil, err
	}

	authTokens := &entities.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return authTokens, nil
}

// Logout invalidates the specified refresh token, effectively logging the user out.
// Returns an error if the logout operation fails (e.g., if the refresh token is not found).
func (as *AuthService) Logout(ctx context.Context, refreshToken string) error {
	return as.tokenSvc.RevokeRefreshToken(ctx, refreshToken)
}
