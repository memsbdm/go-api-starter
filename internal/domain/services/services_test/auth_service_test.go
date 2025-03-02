//go:build !integration

package services_test

import (
	"context"
	"errors"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"testing"
)

func TestAuthService_Login(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()
	builder := NewTestBuilder().Build()

	userToCreate := newValidUserToCreate()

	_, err := builder.AuthService.Register(ctx, userToCreate)
	if err != nil {
		t.Fatalf("failed to register user: %v", err)
	}

	type loginRequest struct {
		username string
		password string
	}

	tests := map[string]struct {
		input       *loginRequest
		expectedErr error
	}{
		"login success": {
			input: &loginRequest{
				username: userToCreate.Username,
				password: userToCreate.Password,
			},
			expectedErr: nil,
		},
		"login with bad credentials": {
			input: &loginRequest{
				username: "not-existing",
				password: "not-existing",
			},
			expectedErr: domain.ErrInvalidCredentials,
		},
	}

	// Act & Assert
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			_, _, err := builder.AuthService.Login(ctx, tt.input.username, tt.input.password)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestAuthService_Register_SendsEmail(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()
	builder := NewTestBuilder().Build()

	// Act & Assert
	userToCreate := newValidUserToCreate()
	_, err := builder.AuthService.Register(ctx, userToCreate)
	if err != nil {
		t.Fatalf("error while registering user: %v", err)
	}

	if v, ok := builder.MailerAdapter.(interface{ SentEmailsCount() int }); ok {
		if v.SentEmailsCount() != 1 {
			t.Errorf("expected 1 email to be sent, got %d", v.SentEmailsCount())
		}
	}
}

func TestAuthService_SendPasswordResetEmail(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()
	builder := NewTestBuilder().Build()

	tests := map[string]struct {
		input                  string
		prepare                func(builder *TestBuilder)
		expectedNbOfEmailsSent int
		expectedErr            error
	}{
		"send password reset email": {
			input: newValidUserToCreate().Email,
			prepare: func(builder *TestBuilder) {
				userToCreate := newValidUserToCreate()
				_, err := builder.AuthService.Register(ctx, userToCreate)
				if err != nil {
					t.Fatalf("error while registering user: %v", err)
				}
			},
			expectedNbOfEmailsSent: 1,
			expectedErr:            nil,
		},
		"send password reset email with non-existing user": {
			input:                  "non-existing@test.com",
			expectedNbOfEmailsSent: 0,
			expectedErr:            nil,
		},
	}

	// Act & Assert
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if tt.prepare != nil {
				tt.prepare(builder)
			}

			err := builder.AuthService.SendPasswordResetEmail(ctx, tt.input)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}

			if v, ok := builder.MailerAdapter.(interface{ SentEmailsCount() int }); ok {
				if v.SentEmailsCount() != tt.expectedNbOfEmailsSent {
					t.Errorf("expected %d emails to be sent, got %d", tt.expectedNbOfEmailsSent, v.SentEmailsCount())
				}
			}
		})
	}
}

func TestAuthService_ResetPassword(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		prepare     func(t *testing.T, builder *TestBuilder, ctx context.Context) (string, *entities.User)
		newPassword string
		confirm     string
		expectedErr error
	}{
		"successful password reset": {
			prepare: func(t *testing.T, builder *TestBuilder, ctx context.Context) (string, *entities.User) {
				return setupVerifiedUserWithResetToken(t, ctx, builder)

			},
			newPassword: "new-password",
			confirm:     "new-password",
			expectedErr: nil,
		},
		"password mismatch": {
			prepare: func(t *testing.T, builder *TestBuilder, ctx context.Context) (string, *entities.User) {
				return setupVerifiedUserWithResetToken(t, ctx, builder)
			},
			newPassword: "new-password",
			confirm:     "different-password",
			expectedErr: domain.ErrPasswordsNotMatch,
		},
		"invalid token": {
			prepare: func(t *testing.T, builder *TestBuilder, ctx context.Context) (string, *entities.User) {
				_, user := setupVerifiedUserWithResetToken(t, ctx, builder)
				return "invalid-token", user
			},
			newPassword: "new-password",
			confirm:     "new-password",
			expectedErr: domain.ErrInvalidToken,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			ctx := context.Background()
			builder := NewTestBuilder().Build()
			token, user := tt.prepare(t, builder, ctx)

			// Act & Assert
			err := builder.AuthService.ResetPassword(ctx, token, tt.newPassword, tt.confirm)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}

			if tt.expectedErr == nil {
				// Verify the password was actually changed by attempting to login
				_, _, err = builder.AuthService.Login(ctx, user.Username, tt.newPassword)
				if err != nil {
					t.Errorf("failed to login with new password: %v", err)
				}
			}
		})
	}
}

func setupVerifiedUserWithResetToken(t *testing.T, ctx context.Context, builder *TestBuilder) (string, *entities.User) {
	t.Helper()

	userToCreate := newValidUserToCreate()
	user, err := builder.AuthService.Register(ctx, userToCreate)
	if err != nil {
		t.Fatalf("error while registering user: %v", err)
	}

	_, err = builder.UserRepo.VerifyEmail(ctx, user.ID)
	if err != nil {
		t.Fatalf("error while verifying email: %v", err)
	}

	token, err := builder.TokenService.GenerateOneTimeToken(ctx, entities.PasswordResetToken, user.ID)
	if err != nil {
		t.Fatalf("error while generating one-time token: %v", err)
	}

	return token, user
}
