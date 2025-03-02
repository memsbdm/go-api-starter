package services

import (
	"context"
	"errors"
	"go-starter/config"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"go-starter/internal/domain/helpers"
	"go-starter/internal/domain/mailtemplates"
	"go-starter/internal/domain/ports"
	"go-starter/internal/domain/utils"
	"regexp"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// UserService implements ports.UserService interface and provides access to the user repository.
type UserService struct {
	repo      ports.UserRepository
	cacheSvc  ports.CacheService
	tokenSvc  ports.TokenService
	mailerSvc ports.MailerService
	appCfg    *config.App
}

// NewUserService creates a new instance of UserService.
func NewUserService(appCfg *config.App, repo ports.UserRepository, cacheSvc ports.CacheService, tokenSvc ports.TokenService, mailerSvc ports.MailerService) *UserService {
	return &UserService{
		repo:      repo,
		cacheSvc:  cacheSvc,
		tokenSvc:  tokenSvc,
		mailerSvc: mailerSvc,
		appCfg:    appCfg,
	}
}

// UserCachePrefix is the prefix for caching users.
const UserCachePrefix = "user"

// GetByID retrieves a user by their unique identifier.
// Returns the user entity if found or an error if not found or any other issue occurs.
func (us *UserService) GetByID(ctx context.Context, id entities.UserID) (*entities.User, error) {
	cacheUser, err := us.getUserFromCache(ctx, id)
	if err == nil {
		return cacheUser, nil
	}

	user, err := us.repo.GetByID(ctx, id.UUID())
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	err = us.cacheUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByUsername retrieves a user by their username.
// Returns the user entity if found or an error if not found or any other issue occurs.
func (us *UserService) GetByUsername(ctx context.Context, username string) (*entities.User, error) {
	user, err := us.repo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return user, nil
}

// GetIDByVerifiedEmail retrieves a user ID by their verified email.
// Returns the user ID if found or an error if not found or any other issue occurs.
func (us *UserService) GetIDByVerifiedEmail(ctx context.Context, email string) (entities.UserID, error) {
	userID, err := us.repo.GetIDByVerifiedEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return entities.UserID{}, err
		}
		return entities.UserID{}, domain.ErrInternal
	}

	return entities.UserID(userID), nil
}

// Register creates a new user account in the system.
// Returns the created user or an error if the registration fails (e.g., due to validation issues).
func (us *UserService) Register(ctx context.Context, user *entities.User) (*entities.User, error) {
	err := validateRegisterRequest(user)
	if err != nil {
		return nil, err
	}

	if err := us.repo.CheckEmailAvailability(ctx, user.Email); err != nil {
		if errors.Is(err, domain.ErrEmailConflict) {
			return nil, err
		}
		return nil, domain.ErrInternal
	}
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, domain.ErrInternal
	}

	userToCreate := &entities.User{
		Name:     user.Name,
		Username: user.Username,
		Password: hashedPassword,
		Email:    user.Email,
	}

	created, err := us.repo.Create(ctx, userToCreate)
	if err != nil {
		if errors.Is(err, domain.ErrUsernameConflict) {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return created, nil
}

// VerifyEmail verifies a user email.
// Returns an error if the verification of the token fails or the email is already verified by another user.
func (us *UserService) VerifyEmail(ctx context.Context, token string) error {
	userID, err := us.tokenSvc.VerifyAndConsumeOneTimeToken(ctx, entities.EmailVerificationToken, token)
	if err != nil {
		return err
	}

	user, err := us.repo.VerifyEmail(ctx, userID.UUID())
	if err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyVerified) {
			return err
		}
		return domain.ErrInternal
	}

	err = us.cacheUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

// ResendEmailVerification resends a user email verification email.
// Returns an error if the resend fails.
func (us *UserService) ResendEmailVerification(ctx context.Context, userID entities.UserID) error {
	user, err := us.GetByID(ctx, userID)
	if err != nil {
		return domain.ErrInternal
	}

	if user.IsEmailVerified {
		return domain.ErrEmailAlreadyVerified
	}

	token, err := us.tokenSvc.GenerateOneTimeToken(ctx, entities.EmailVerificationToken, userID)
	if err != nil {
		return err
	}

	err = us.mailerSvc.Send(&ports.EmailMessage{
		To:      []string{user.Email},
		Subject: "Verify your email!",
		Body:    mailtemplates.VerifyEmail(us.appCfg.BaseURL, token),
	})
	if err != nil {
		return err
	}

	return nil
}

// UpdatePassword updates a user password.
// Returns an error if the update fails (e.g., due to validation issues).
func (us *UserService) UpdatePassword(ctx context.Context, userID entities.UserID, params entities.UpdateUserParams) error {
	if params.Password == nil {
		return domain.ErrPasswordRequired
	}
	if params.PasswordConfirmation == nil {
		return domain.ErrPasswordConfirmationRequired
	}
	if *params.Password != *params.PasswordConfirmation {
		return domain.ErrPasswordsNotMatch
	}
	if len(*params.Password) < domain.PasswordMinLength {
		return domain.ErrPasswordTooShort
	}
	hashedPassword, err := utils.HashPassword(*params.Password)
	if err != nil {
		return domain.ErrInternal
	}
	err = us.repo.UpdatePassword(ctx, userID.UUID(), hashedPassword)
	if err != nil {
		return domain.ErrInternal
	}
	return nil
}

// validateUsername checks if the provided username meets the required criteria.
// Returns an error if any validation fails.
func validateUsername(username string) error {
	if len(username) < domain.UsernameMinLength {
		return domain.ErrUsernameTooShort
	}
	if len(username) > domain.UsernameMaxLength {
		return domain.ErrUsernameTooLong
	}
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !usernameRegex.MatchString(username) {
		return domain.ErrUsernameInvalid
	}
	return nil
}

// validatePassword checks if the provided password meets the required criteria.
// Returns an error if any validation fails.
func validatePassword(password string) error {
	if len(password) < domain.PasswordMinLength {
		return domain.ErrPasswordTooShort
	}
	return nil
}

// validateAndFormatName checks if the provided name meets the required criteria and formats it.
// Returns an error if any validation fails.
func validateAndFormatName(user *entities.User) error {
	if len(user.Name) > domain.NameMaxLength {
		return domain.ErrNameTooLong
	}

	fields := strings.Fields(user.Name)
	for i := range fields {
		titleCaser := cases.Title(language.English)
		fields[i] = titleCaser.String(strings.ToLower(fields[i]))
	}

	user.Name = strings.Join(fields, " ")

	return nil
}

// validateEmail validates an email address.
// Returns an error if the email is invalid.
func validateEmail(email string) error {
	ok := helpers.IsValidEmail(email)
	if !ok {
		return domain.ErrEmailInvalid
	}

	return nil
}

// validateRegisterRequest validates the registration details of a user.
// Returns an error if any validation fails.
func validateRegisterRequest(user *entities.User) error {
	if err := validateAndFormatName(user); err != nil {
		return err
	}
	if err := validateUsername(user.Username); err != nil {
		return err
	}
	if err := validatePassword(user.Password); err != nil {
		return err
	}
	if err := validateEmail(user.Email); err != nil {
		return err
	}
	return nil
}

// getUserFromCache retrieves a user from the cache.
// Returns an error if the retrieval fails.
func (us *UserService) getUserFromCache(ctx context.Context, userID entities.UserID) (*entities.User, error) {
	var user *entities.User
	cacheKey := utils.GenerateCacheKey(UserCachePrefix, userID.String())
	cachedUser, err := us.cacheSvc.Get(ctx, cacheKey)
	if err != nil {
		return nil, err
	}

	err = utils.Deserialize(cachedUser, &user)
	if err != nil {
		return nil, domain.ErrInternal
	}
	return user, nil
}

// cacheUser caches a user in the cache.
// Returns an error if the caching fails.
func (us *UserService) cacheUser(ctx context.Context, user *entities.User) error {
	userSerialized, err := utils.Serialize(user)
	if err != nil {
		return domain.ErrInternal
	}

	cacheKey := utils.GenerateCacheKey(UserCachePrefix, user.ID.String())

	err = us.cacheSvc.Set(ctx, cacheKey, userSerialized, time.Hour)
	if err != nil {
		return domain.ErrInternal
	}
	return nil
}
