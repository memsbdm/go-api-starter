package services

import (
	"context"
	"errors"
	"go-starter/internal/adapters/http/helpers"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
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
	repo     ports.UserRepository
	cacheSvc ports.CacheService
	tokenSvc ports.TokenService
}

// NewUserService creates a new instance of UserService.
func NewUserService(repo ports.UserRepository, cacheSvc ports.CacheService, tokenSvc ports.TokenService) *UserService {
	return &UserService{
		repo:     repo,
		cacheSvc: cacheSvc,
		tokenSvc: tokenSvc,
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

// Register creates a new user account in the system.
// Returns the created user or an error if the registration fails (e.g., due to validation issues).
func (us *UserService) Register(ctx context.Context, user *entities.User) (*entities.User, error) {
	err := validateRegisterRequest(user)
	if err != nil {
		return nil, err
	}

	if err := us.repo.CheckEmailAvailability(ctx, user.Email); err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyTaken) {
			return nil, domain.ErrEmailAlreadyTaken
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
		if errors.Is(err, domain.ErrUsernameAlreadyTaken) {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return created, nil
}

func (us *UserService) VerifyEmail(ctx context.Context, token string) error {
	userID, err := us.tokenSvc.VerifyAndInvalidateSecureToken(ctx, entities.EmailVerificationToken, token)
	if err != nil {
		return err
	}

	user, err := us.repo.VerifyEmail(ctx, userID)
	if err != nil {
		return domain.ErrInternal
	}

	err = us.cacheUser(ctx, user)
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

// validateAndFormatUsername checks if the provided username meets the required criteria.
// Returns an error if any validation fails.
func validateAndFormatUsername(user *entities.User) error {
	user.Username = strings.TrimSpace(user.Username)
	if user.Username == "" {
		return domain.ErrUsernameRequired
	}
	if len(user.Username) < domain.UsernameMinLength {
		return domain.ErrUsernameTooShort
	}
	if len(user.Username) > domain.UsernameMaxLength {
		return domain.ErrUsernameTooLong
	}
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !usernameRegex.MatchString(user.Username) {
		return domain.ErrUsernameInvalid
	}
	return nil
}

// validatePassword checks if the provided password meets the required criteria.
// Returns an error if any validation fails.
func validatePassword(password string) error {
	if password == "" {
		return domain.ErrPasswordRequired
	}
	if len(password) < domain.PasswordMinLength {
		return domain.ErrPasswordTooShort
	}
	return nil
}

// validateAndFormatName checks if the provided name meets the required criteria and formats it.
// Returns an error if any validation fails.
func validateAndFormatName(user *entities.User) error {
	user.Name = strings.TrimSpace(user.Name)
	if user.Name == "" {
		return domain.ErrNameRequired
	}
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

// validateAndFormatEmail validates and formats a user email.
// Returns an error if any validation fails.
func validateAndFormatEmail(user *entities.User) error {
	user.Email = strings.TrimSpace(user.Email)
	if user.Email == "" {
		return domain.ErrEmailRequired
	}

	ok := helpers.IsValidEmail(user.Email)
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
	if err := validateAndFormatUsername(user); err != nil {
		return err
	}
	if err := validatePassword(user.Password); err != nil {
		return err
	}
	if err := validateAndFormatEmail(user); err != nil {
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
