package services

import (
	"context"
	"errors"
	"go-starter/internal/domain"
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

	accessToken, err := as.tokenSvc.GenerateToken(user)
	if err != nil {
		return "", "", domain.ErrInternal
	}

	refreshToken, err := as.tokenSvc.GenerateRefreshToken(user)
	if err != nil {
		return "", "", domain.ErrInternal
	}

	return accessToken, refreshToken, nil
}
