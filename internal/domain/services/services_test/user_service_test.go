//go:build !integration

package services_test

import (
	"bytes"
	"context"
	"errors"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"go-starter/internal/domain/services"
	"go-starter/internal/domain/utils"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func newValidUserToCreate() *entities.User {
	return &entities.User{
		Name:     "John Doe",
		Username: "example",
		Password: "secret123",
		Email:    "example@example.com",
	}
}

func TestUserService_Register(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()
	builder := NewTestBuilder().Build()
	userToCreate := newValidUserToCreate()
	createdUser, err := builder.UserService.Register(ctx, userToCreate)
	if err != nil {
		t.Fatalf("error while registering user: %v", err)
	}

	tests := map[string]struct {
		input       *entities.User
		expectedErr error
	}{
		"register valid user": {
			input: &entities.User{
				Name:     userToCreate.Name,
				Username: "success",
				Password: userToCreate.Password,
				Email:    userToCreate.Email,
			},
			expectedErr: nil,
		},
		"register user with conflicting username": {
			input: &entities.User{
				Name:     userToCreate.Name,
				Username: createdUser.Username,
				Password: userToCreate.Password,
				Email:    userToCreate.Email,
			},
			expectedErr: domain.ErrUsernameConflict,
		},
		"register user with conflicting not verified email": {
			input: &entities.User{
				Name:     userToCreate.Name,
				Username: "conflict",
				Password: userToCreate.Password,
				Email:    userToCreate.Email,
			},
			expectedErr: nil,
		},
		"register user with short password": {
			input: &entities.User{
				Name:     userToCreate.Name,
				Username: createdUser.Username,
				Password: "short",
				Email:    userToCreate.Email,
			},
			expectedErr: domain.ErrPasswordTooShort,
		},
		"register user with short username": {
			input: &entities.User{
				Name:     userToCreate.Name,
				Username: strings.Repeat("x", domain.UsernameMinLength-1),
				Password: userToCreate.Password,
				Email:    userToCreate.Email,
			},
			expectedErr: domain.ErrUsernameTooShort,
		},
		"register user with long username": {
			input: &entities.User{
				Name:     userToCreate.Name,
				Username: strings.Repeat("x", domain.UsernameMaxLength+1),
				Password: userToCreate.Password,
				Email:    userToCreate.Email,
			},
			expectedErr: domain.ErrUsernameTooLong,
		},
		"register user with invalid username": {
			input: &entities.User{
				Name:     userToCreate.Name,
				Username: "invalid%@",
				Password: userToCreate.Password,
				Email:    userToCreate.Email,
			},
			expectedErr: domain.ErrUsernameInvalid,
		},
		"register user with long name": {
			input: &entities.User{
				Name:     strings.Repeat("x", domain.NameMaxLength+1),
				Username: userToCreate.Username,
				Password: userToCreate.Password,
				Email:    userToCreate.Email,
			},
			expectedErr: domain.ErrNameTooLong,
		},
		"register user with long name containing emojis": {
			input: &entities.User{
				Name:     strings.Repeat("🥵", domain.NameMaxLength/4+1),
				Username: userToCreate.Username,
				Password: userToCreate.Password,
				Email:    userToCreate.Email,
			},
			expectedErr: domain.ErrNameTooLong,
		},
	}

	// Act & Assert
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			result, err := builder.UserService.Register(ctx, tt.input)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if err == nil && result.Username != tt.input.Username {
				t.Errorf("expected username %s, got %s", tt.input.Username, result.Username)
			}
		})
	}
}

func TestUserService_Register_EmailVerification(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()
	builder := NewTestBuilder().Build()

	userToCreate := newValidUserToCreate()
	createdUser, err := builder.UserService.Register(ctx, userToCreate)
	if err != nil {
		t.Fatalf("error while registering user: %v", err)
	}

	_, err = builder.UserRepo.VerifyEmail(ctx, createdUser.ID)
	if err != nil {
		t.Fatalf("error while verifying email: %v", err)
	}

	testedUser := &entities.User{
		Username: "valid",
		Password: userToCreate.Password,
		Name:     userToCreate.Name,
		Email:    userToCreate.Email,
	}

	// Act & Assert
	newUser, err := builder.UserService.Register(ctx, testedUser)
	if !errors.Is(err, domain.ErrEmailConflict) {
		t.Errorf("expected error %v, got %v", domain.ErrEmailConflict, err)
	}
	if newUser != nil {
		t.Errorf("expected user to be nil, got %v", newUser)
	}
}

func TestUserService_GetByID(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()
	builder := NewTestBuilder().Build()

	userToCreate := newValidUserToCreate()

	createdUser, err := builder.UserService.Register(ctx, userToCreate)
	if err != nil {
		t.Fatalf("error while registering user: %v", err)
	}

	tests := map[string]struct {
		input       entities.UserID
		expectedErr error
	}{
		"get valid user by id": {
			input:       createdUser.ID,
			expectedErr: nil,
		},
		"get non-existing user by id": {
			input:       entities.UserID(uuid.New()),
			expectedErr: domain.ErrUserNotFound,
		},
	}

	// Act & Assert
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			_, err := builder.UserService.GetByID(ctx, tt.input)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestUserService_GetByID_Cache(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()
	builder := NewTestBuilder().Build()

	userToCreate := newValidUserToCreate()
	createdUser, err := builder.UserService.Register(ctx, userToCreate)
	if err != nil {
		t.Fatalf("error while registering user: %v", err)
	}

	_, err = builder.UserService.GetByID(ctx, createdUser.ID)
	if err != nil {
		t.Fatalf("error while fetching user: %v", err)
	}

	// Act & Assert
	cachedUser, err := builder.CacheService.Get(ctx, utils.GenerateCacheKey(services.UserCachePrefix, createdUser.ID))
	if err != nil {
		t.Errorf("error while getting user from cache: %v", err)
	}

	var deserializedUser entities.User
	err = utils.Deserialize(cachedUser, &deserializedUser)
	if err != nil {
		t.Errorf("error while deserializing user: %v", err)
	}

	if deserializedUser.ID != createdUser.ID {
		t.Errorf("deserialized user does not match cache")
	}
}

func TestUserService_GetByUsername(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()
	builder := NewTestBuilder().Build()

	userToCreate := newValidUserToCreate()

	createdUser, err := builder.UserService.Register(ctx, userToCreate)
	if err != nil {
		t.Fatalf("error while registering user: %v", err)
	}

	tests := map[string]struct {
		input       string
		expectedErr error
	}{
		"get valid user by username": {
			input:       createdUser.Username,
			expectedErr: nil,
		},
		"get non-existing user by username": {
			input:       "non-existing",
			expectedErr: domain.ErrUserNotFound,
		},
	}

	// Act & Assert
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			result, err := builder.UserService.GetByUsername(ctx, tt.input)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if err == nil && result.Username != tt.input {
				t.Errorf("expected username %s, got %s", tt.input, result.Username)
			}
		})
	}
}

func TestUserService_VerifyEmail(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()

	tests := map[string]struct {
		prepare     func() (*TestBuilder, string, entities.UserID)
		expectedErr error
	}{
		"verify email successfully": {
			prepare: func() (*TestBuilder, string, entities.UserID) {
				builder := NewTestBuilder().Build()
				user, err := builder.UserService.Register(ctx, newValidUserToCreate())
				if err != nil {
					t.Fatalf("error while registering user: %v", err)
				}
				token, err := builder.TokenService.GenerateOneTimeToken(ctx, entities.EmailVerificationToken, user.ID)
				if err != nil {
					t.Fatalf("error while generating token: %v", err)
				}
				return builder, token, user.ID
			},
			expectedErr: nil,
		},
		"verify email should fail for an already verified email": {
			prepare: func() (*TestBuilder, string, entities.UserID) {
				userToCreate := newValidUserToCreate()
				builder := NewTestBuilder().Build()
				user, err := builder.UserService.Register(ctx, userToCreate)
				if err != nil {
					t.Fatalf("error while registering user: %v", err)
				}

				user2, err := builder.UserService.Register(ctx, &entities.User{
					Username: "conflict",
					Email:    user.Email, // Same email as the first user
					Password: user.Password,
					Name:     user.Name,
				})
				if err != nil {
					t.Fatalf("error while registering user: %v", err)
				}

				token, err := builder.TokenService.GenerateOneTimeToken(ctx, entities.EmailVerificationToken, user.ID)
				if err != nil {
					t.Fatalf("error while generating token: %v", err)
				}

				token2, err := builder.TokenService.GenerateOneTimeToken(ctx, entities.EmailVerificationToken, user2.ID)
				if err != nil {
					t.Fatalf("error while generating token: %v", err)
				}

				// Verify the email of the first user, the second still have a verification email with a token
				err = builder.UserService.VerifyEmail(ctx, token)
				if err != nil {
					t.Fatalf("error while verifying email: %v", err)
				}
				return builder, token2, user2.ID
			},
			expectedErr: domain.ErrEmailAlreadyVerified,
		},
	}

	// Act & Assert
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			builder, token, userID := tt.prepare()
			err := builder.UserService.VerifyEmail(ctx, token)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if tt.expectedErr == nil {
				user, err := builder.CacheService.Get(ctx, utils.GenerateCacheKey(services.UserCachePrefix, userID))
				if err != nil {
					t.Errorf("error while getting user from cache: %v", err)
				}

				var deserializedUser entities.User
				err = utils.Deserialize(user, &deserializedUser)
				if err != nil {
					t.Errorf("error while deserializing user: %v", err)
				}

				if !deserializedUser.IsEmailVerified {
					t.Errorf("expected user to be verified in cache, got %v", deserializedUser.IsEmailVerified)
				}
			}
		})
	}
}
func TestUserService_ResendEmailVerification(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()

	tests := map[string]struct {
		prepare     func(t *testing.T) (entities.UserID, *TestBuilder)
		expectedErr error
	}{
		"resend email verification successfully": {
			prepare: func(t *testing.T) (entities.UserID, *TestBuilder) {
				builder := NewTestBuilder().Build()
				user, err := builder.UserService.Register(ctx, newValidUserToCreate())
				if err != nil {
					t.Fatalf("error while registering user: %v", err)
				}
				return user.ID, builder
			},
			expectedErr: nil,
		},
		"resend email verification should fail for an already verified user": {
			prepare: func(t *testing.T) (entities.UserID, *TestBuilder) {
				builder := NewTestBuilder().Build()
				user, err := builder.UserService.Register(ctx, newValidUserToCreate())
				if err != nil {
					t.Fatalf("error while registering user: %v", err)
				}

				_, err = builder.UserRepo.VerifyEmail(ctx, user.ID)
				if err != nil {
					t.Fatalf("error while verifying email: %v", err)
				}
				return user.ID, builder
			},
			expectedErr: domain.ErrEmailAlreadyVerified,
		},
	}

	// Act & Assert
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			userID, builder := tt.prepare(t)
			err := builder.UserService.ResendEmailVerification(ctx, userID)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestUserService_UpdatePassword(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()
	builder := NewTestBuilder().Build()
	userToCreate := newValidUserToCreate()
	user, err := builder.UserService.Register(ctx, userToCreate)
	if err != nil {
		t.Fatalf("error while registering user: %v", err)
	}

	validPassword := "secret123"
	shortPassword := "short"
	notMatchingPassword := "not-matching"

	tests := map[string]struct {
		input       entities.UpdateUserParams
		expectedErr error
	}{
		"update password successfully": {
			input: entities.UpdateUserParams{
				Password:             &validPassword,
				PasswordConfirmation: &validPassword,
			},
			expectedErr: nil,
		},
		"update with short password": {
			input: entities.UpdateUserParams{
				Password:             &shortPassword,
				PasswordConfirmation: &shortPassword,
			},
			expectedErr: domain.ErrPasswordTooShort,
		},
		"update with non-matching password": {
			input: entities.UpdateUserParams{
				Password:             &validPassword,
				PasswordConfirmation: &notMatchingPassword,
			},
			expectedErr: domain.ErrPasswordsNotMatch,
		},
		"update with missing password": {
			input: entities.UpdateUserParams{
				PasswordConfirmation: &validPassword,
			},
			expectedErr: domain.ErrPasswordRequired,
		},
		"update with missing password confirmation": {
			input: entities.UpdateUserParams{
				Password: &validPassword,
			},
			expectedErr: domain.ErrPasswordConfirmationRequired,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := builder.UserService.UpdatePassword(ctx, user.ID, tt.input)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestUserService_UpdateAvatar_Is_Caching_User(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()
	builder := NewTestBuilder().Build()
	user, err := builder.UserService.Register(ctx, newValidUserToCreate())
	if err != nil {
		t.Fatalf("error while registering user: %v", err)
	}

	tests := map[string]struct {
		input           string
		expectedURL     string
		expectCachedURL bool
	}{
		"update avatar successfully": {
			input:           "avatar.jpg",
			expectedURL:     "https://example.com/avatars/" + user.ID.String() + ".jpg",
			expectCachedURL: true,
		},
	}

	// Act & Assert
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			emptyFile := bytes.NewBuffer([]byte{})
			_, err := builder.UserService.UpdateAvatar(ctx, user.ID, tt.input, emptyFile)
			if err != nil {
				t.Fatalf("error while updating avatar: %v", err)
			}

			cachedUser, err := builder.CacheService.Get(ctx, utils.GenerateCacheKey(services.UserCachePrefix, user.ID))
			if err != nil {
				t.Errorf("error while getting user from cache: %v", err)
			}

			var deserializedUser entities.User
			err = utils.Deserialize(cachedUser, &deserializedUser)
			if err != nil {
				t.Fatalf("error while deserializing user: %v", err)
			}

			if tt.expectCachedURL && *deserializedUser.AvatarURL != tt.expectedURL {
				t.Errorf("expected cached URL to be %s, got %s", tt.expectedURL, *deserializedUser.AvatarURL)
			}
		})
	}
}

func TestUserService_DeleteAvatar_Is_Caching_User(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()
	builder := NewTestBuilder().Build()
	user, err := builder.UserService.Register(ctx, newValidUserToCreate())
	if err != nil {
		t.Fatalf("error while registering user: %v", err)
	}

	tests := map[string]struct {
		expectedErr error
	}{
		"delete avatar successfully": {
			expectedErr: nil,
		},
	}

	// Act & Assert
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := builder.UserService.DeleteAvatar(ctx, user.ID)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}

			cachedUser, err := builder.CacheService.Get(ctx, utils.GenerateCacheKey(services.UserCachePrefix, user.ID))
			if err != nil {
				t.Errorf("error while getting user from cache: %v", err)
			}

			var deserializedUser entities.User
			err = utils.Deserialize(cachedUser, &deserializedUser)
			if err != nil {
				t.Fatalf("error while deserializing user: %v", err)
			}

			if deserializedUser.AvatarURL != nil {
				t.Errorf("expected avatar URL to be empty, got %s", *deserializedUser.AvatarURL)
			}
		})
	}
}
