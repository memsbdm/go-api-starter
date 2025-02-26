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
func (as *AuthService) Login(ctx context.Context, username, password string) (*entities.User, string, error) {
	if strings.TrimSpace(username) == "" {
		return nil, "", domain.ErrUsernameRequired
	}
	user, err := as.userSvc.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, "", domain.ErrInvalidCredentials
		}
		return nil, "", domain.ErrInternal
	}

	err = utils.ComparePassword(password, user.Password)
	if err != nil {
		return nil, "", domain.ErrInvalidCredentials
	}

	accessToken, err := as.tokenSvc.GenerateAuthToken(ctx, user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, accessToken, nil
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

// Logout logs out a user from the system.
// Returns an error if the logout fails.
func (as *AuthService) Logout(ctx context.Context, accessToken string) error {
	return as.tokenSvc.RevokeAuthToken(ctx, accessToken)
}
