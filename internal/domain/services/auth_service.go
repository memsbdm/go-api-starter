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
		return "", "", err
	}

	refreshToken, err := as.tokenSvc.GenerateRefreshToken(ctx, user.ID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// Register registers a new user
func (as *AuthService) Register(ctx context.Context, user *entities.User) (*entities.User, error) {
	return as.userSvc.Register(ctx, user)
}

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

func (as *AuthService) Logout(ctx context.Context, refreshToken string) error {
	return as.tokenSvc.RevokeRefreshToken(ctx, refreshToken)
}
