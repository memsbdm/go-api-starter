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
)

// AuthService implements ports.AuthService interface.
type AuthService struct {
	cfg       *config.Container
	userSvc   ports.UserService
	tokenSvc  ports.TokenService
	mailerSvc ports.MailerService
}

// NewAuthService creates a new instance of AuthService.
func NewAuthService(
	cfg *config.Container,
	userSvc ports.UserService,
	tokenSvc ports.TokenService,
	mailerSvc ports.MailerService,
) *AuthService {

	return &AuthService{
		cfg:       cfg,
		userSvc:   userSvc,
		tokenSvc:  tokenSvc,
		mailerSvc: mailerSvc,
	}
}

// Login authenticates a user.
// Returns auth tokens upon successful authentication,
// or an error if the login fails (e.g., due to incorrect credentials).
func (as *AuthService) Login(ctx context.Context, username, password string) (*entities.User, string, error) {
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
		Body:    mailtemplates.VerifyEmail(as.cfg.Application.BaseURL, token, as.cfg.Token.EmailVerificationTokenDuration),
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

// SendPasswordResetEmail sends a password reset email to the user.
// Returns an error if the email fails to send.
func (as *AuthService) SendPasswordResetEmail(ctx context.Context, email string) error {
	userID, err := as.userSvc.GetIDByVerifiedEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil
		}
		return err
	}

	token, err := as.tokenSvc.GenerateOneTimeToken(ctx, entities.PasswordResetToken, userID)
	if err != nil {
		return err
	}

	err = as.mailerSvc.Send(&ports.EmailMessage{
		To:      []string{email},
		Subject: "Reset your password!",
		Body:    mailtemplates.ResetPassword(as.cfg.Application.BaseURL, token, as.cfg.Token.PasswordResetTokenDuration),
	})
	if err != nil {
		return err
	}

	return nil
}

// VerifyPasswordResetToken verifies a password reset token.
// Returns an error if the token is invalid.
func (as *AuthService) VerifyPasswordResetToken(ctx context.Context, token string) error {
	_, err := as.tokenSvc.VerifyOneTimeToken(ctx, entities.PasswordResetToken, token)
	return err
}

// ResetPassword resets a user's password.
// Returns an error if the password reset fails.
func (as *AuthService) ResetPassword(ctx context.Context, token, password, passwordConfirmation string) error {
	userID, err := as.tokenSvc.VerifyOneTimeToken(ctx, entities.PasswordResetToken, token)
	if err != nil {
		return err
	}

	err = as.userSvc.UpdatePassword(ctx, userID, entities.UpdateUserParams{
		Password:             &password,
		PasswordConfirmation: &passwordConfirmation,
	})
	if err != nil {
		return err
	}

	return as.tokenSvc.ConsumeOneTimeToken(ctx, entities.PasswordResetToken, token)
}
