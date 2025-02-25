package services

import (
	"context"
	"errors"
	"go-starter/config"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"go-starter/internal/domain/mailtemplates"
	"go-starter/internal/domain/ports"
	"go-starter/internal/domain/utils"
	"strings"
)

// AuthService implements ports.AuthService interface.
type AuthService struct {
	appCfg     *config.App
	userSvc    ports.UserService
	tokenSvc   ports.TokenService
	errTracker ports.ErrTrackerAdapter
	mailerSvc  ports.MailerService
}

// NewAuthService creates a new instance of AuthService.
func NewAuthService(
	appCfg *config.App,
	userSvc ports.UserService,
	tokenSvc ports.TokenService,
	errTracker ports.ErrTrackerAdapter,
	mailerSvc ports.MailerService,
) *AuthService {

	return &AuthService{
		appCfg:     appCfg,
		userSvc:    userSvc,
		tokenSvc:   tokenSvc,
		errTracker: errTracker,
		mailerSvc:  mailerSvc,
	}
}

// Login authenticates a user.
// Returns auth tokens upon successful authentication,
// or an error if the login fails (e.g., due to incorrect credentials).
func (as *AuthService) Login(ctx context.Context, username, password string) (*entities.User, *entities.AuthTokens, error) {
	if strings.TrimSpace(username) == "" {
		return nil, nil, domain.ErrUsernameRequired
	}
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
	createdUser, err := as.userSvc.Register(ctx, user)
	if err != nil {
		return nil, err
	}

	token, err := as.tokenSvc.GenerateOneTimeToken(ctx, entities.EmailVerificationToken, createdUser.ID)
	if err != nil {
		return nil, err
	}

	err = as.mailerSvc.Send(&ports.EmailMessage{
		To:      []string{createdUser.Email},
		Subject: "Verify your email!",
		Body:    mailtemplates.VerifyEmail(as.appCfg.BaseURL, token),
	})
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

// Refresh generates new access and refresh tokens using the previous refresh token.
// Returns the new auth tokens, and an error if the refresh fails
// (e.g., if the previous refresh token is invalid or expired).
func (as *AuthService) Refresh(ctx context.Context, previousRefreshToken string) (*entities.AuthTokens, error) {
	if previousRefreshToken == "" {
		return nil, domain.ErrRefreshTokenRequired
	}
	claims, err := as.tokenSvc.VerifyAndParseRefreshToken(ctx, previousRefreshToken)
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

	refreshToken, err := as.tokenSvc.GenerateRefreshToken(ctx, user.ID)
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
