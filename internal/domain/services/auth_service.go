package services

import (
	"context"
	"errors"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"go-starter/internal/domain/ports"
	"go-starter/internal/domain/utils"
)

// AuthService implements ports.AuthService interface
type AuthService struct {
	userSvc  ports.UserService
	tokenSvc ports.TokenService
}

// NewAuthService creates a token service instance
func NewAuthService(userSvc ports.UserService, tokenSvc ports.TokenService) *AuthService {
	return &AuthService{
		userSvc:  userSvc,
		tokenSvc: tokenSvc,
	}
}

// Login authenticates a user. Returns a token string pointer (could be JWT, session ID, etc.) or an error if login fails
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
		return "", "", domain.ErrInternal
	}

	refreshToken, err := as.tokenSvc.GenerateRefreshToken(user)
	if err != nil {
		return "", "", domain.ErrInternal
	}

	return accessToken, refreshToken, nil
}

// RefreshToken generates a new access token and a new refresh token. It returns an error if the token is invalid or expired
func (as *AuthService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	claims, err := as.tokenSvc.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", domain.ErrInvalidToken
	}

	user, err := as.userSvc.GetByID(ctx, claims.UserID)
	if err != nil {
		return "", "", domain.ErrInternal
	}

	accessToken, err := as.tokenSvc.GenerateAccessToken(user)
	if err != nil {
		return "", "", domain.ErrInternal
	}

	newRefreshToken, err := as.tokenSvc.GenerateRefreshToken(user)
	if err != nil {
		return "", "", domain.ErrInternal
	}

	return accessToken, newRefreshToken, nil
}

// Register registers a new user
func (as *AuthService) Register(ctx context.Context, user *entities.User) (*entities.User, error) {
	return as.userSvc.Register(ctx, user)
}
