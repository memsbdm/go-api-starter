package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"go-starter/internal/domain/ports"

	"github.com/lib/pq"
)

// UserRepository implements the ports.UserRepository interface and provides access to the database.
type UserRepository struct {
	executor   QueryExecutor
	errTracker ports.ErrTrackerAdapter
}

// NewUserRepository creates and returns a new UserRepository instance.
func NewUserRepository(db *sql.DB, errTracker ports.ErrTrackerAdapter) *UserRepository {
	return &UserRepository{
		executor:   db,
		errTracker: errTracker,
	}
}

// NewUserRepositoryWithExecutor creates and returns a new UserRepository instance with a custom executor.
func NewUserRepositoryWithExecutor(executor QueryExecutor, errTracker ports.ErrTrackerAdapter) *UserRepository {
	return &UserRepository{
		executor:   executor,
		errTracker: errTracker,
	}
}

// UserRepository queries
const (
	getByIDQuery                = `SELECT id, created_at, updated_at, name, username, email, is_email_verified, role_id FROM users WHERE id = $1`
	getByUsernameQuery          = `SELECT id, created_at, updated_at, name, username, password, email, is_email_verified, role_id FROM users WHERE username = $1`
	getIDByVerifiedEmailQuery   = `SELECT id FROM users WHERE email = $1 AND is_email_verified = true`
	checkEmailAvailabilityQuery = `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND is_email_verified = true)`
	createUserQuery             = `INSERT INTO users (name, username, password, email) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at, name, username, email, is_email_verified, role_id`
	updatePasswordQuery         = `UPDATE users SET password = $1 WHERE id = $2 `
	verifyEmailQuery            = `UPDATE users SET is_email_verified = true WHERE id = $1 `
)

// GetByID selects a user by their unique identifier from the database.
// Returns the user entity if found or an error if not found or any other issue occurs.
func (ur *UserRepository) GetByID(ctx context.Context, id entities.UserID) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	var (
		uuidStr string
		user    entities.User
	)
	err := ur.executor.QueryRowContext(ctx, getByIDQuery, id.String()).Scan(&uuidStr, &user.CreatedAt, &user.UpdatedAt, &user.Name, &user.Username, &user.Email, &user.IsEmailVerified, &user.RoleID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, fmt.Errorf("%w: id=%s", domain.ErrUserNotFound, id)
		default:
			err = fmt.Errorf("failed to get user %s: %w", id.String(), err)
			ur.errTracker.CaptureException(err)
			return nil, err
		}
	}

	userID, err := entities.ParseUserID(uuidStr)
	if err != nil {
		err = fmt.Errorf("failed to parse user id %s: %w", uuidStr, err)
		ur.errTracker.CaptureException(err)
		return nil, err
	}
	user.ID = userID

	return &user, nil
}

// GetByUsername selects a user by their username from the database.
// Returns the user entity if found or an error if not found or any other issue occurs.
func (ur *UserRepository) GetByUsername(ctx context.Context, username string) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	user := &entities.User{}
	var uuidStr string
	err := ur.executor.QueryRowContext(ctx, getByUsernameQuery, username).Scan(&uuidStr, &user.CreatedAt, &user.UpdatedAt, &user.Name, &user.Username, &user.Password, &user.Email, &user.IsEmailVerified, &user.RoleID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, domain.ErrUserNotFound
		default:
			err = fmt.Errorf("failed to get user %s: %w", username, err)
			ur.errTracker.CaptureException(err)
			return nil, err
		}
	}

	parsedID, err := entities.ParseUserID(uuidStr)
	if err != nil {
		err = fmt.Errorf("failed to parse user id %s: %w", uuidStr, err)
		ur.errTracker.CaptureException(err)
		return nil, err
	}
	user.ID = parsedID

	return user, nil
}

// GetByVerifiedEmail returns the user ID for a verified email.
// Returns an error if the user is not found or any other issue occurs.
func (ur *UserRepository) GetIDByVerifiedEmail(ctx context.Context, email string) (entities.UserID, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	var uuidStr string
	err := ur.executor.QueryRowContext(ctx, getIDByVerifiedEmailQuery, email).Scan(&uuidStr)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return entities.NilUserID, domain.ErrUserNotFound
		default:
			err = fmt.Errorf("failed to get user %s: %w", email, err)
			ur.errTracker.CaptureException(err)
			return entities.NilUserID, err
		}
	}

	parsedID, err := entities.ParseUserID(uuidStr)
	if err != nil {
		err = fmt.Errorf("failed to parse user id %s: %w", uuidStr, err)
		ur.errTracker.CaptureException(err)
		return entities.NilUserID, err
	}

	return parsedID, nil
}

// CheckEmailAvailability checks if an email is available for registration.
// Returns an error if the email is already taken.
func (ur *UserRepository) CheckEmailAvailability(ctx context.Context, email string) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	var exists bool
	if err := ur.executor.QueryRowContext(ctx, checkEmailAvailabilityQuery, email).Scan(&exists); err != nil {
		err = fmt.Errorf("failed to check email verification status: %w", err)
		ur.errTracker.CaptureException(err)
		return err
	}

	if exists {
		return domain.ErrEmailConflict
	}

	return nil
}

// Create inserts a new user into the database.
// Returns the created user or an error if the operation fails (e.g., due to a database constraint violation).
func (ur *UserRepository) Create(ctx context.Context, user *entities.User) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var (
		uuidStr     string
		createdUser = *user
	)

	err := ur.executor.QueryRowContext(
		ctx,
		createUserQuery,
		user.Name,
		user.Username,
		user.Password,
		user.Email,
	).Scan(
		&uuidStr,
		&createdUser.CreatedAt,
		&createdUser.UpdatedAt,
		&createdUser.Name,
		&createdUser.Username,
		&createdUser.Email,
		&createdUser.IsEmailVerified,
		&createdUser.RoleID,
	)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23505": // Code unique_violation
				if pqErr.Constraint == "users_username_key" {
					return nil, domain.ErrUsernameConflict
				}
			}
		}
		err = fmt.Errorf("failed to insert user %s: %w", user.Username, err)
		ur.errTracker.CaptureException(err)
		return nil, err
	}

	parsedID, err := entities.ParseUserID(uuidStr)
	if err != nil {
		err = fmt.Errorf("failed to parse user id %s: %w", uuidStr, err)
		ur.errTracker.CaptureException(err)
		return nil, err
	}
	user.ID = parsedID

	return user, nil
}

// UpdatePassword updates a user password.
// Returns an error if the update fails.
func (ur *UserRepository) UpdatePassword(ctx context.Context, userID entities.UserID, newPassword string) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := ur.executor.ExecContext(ctx, updatePasswordQuery, newPassword, userID.String())
	if err != nil {
		err = fmt.Errorf("failed to update user password for user %s: %w", userID.String(), err)
		ur.errTracker.CaptureException(err)
		return err
	}

	return nil
}

// VerifyEmail updates the email verification status of a user.
func (ur *UserRepository) VerifyEmail(ctx context.Context, userID entities.UserID) (*entities.User, error) {
	var returnedUser *entities.User
	return returnedUser, withTx(ur.executor.(*sql.DB), ctx, ur.errTracker, func(tx *sql.Tx) error {
		txRepo := NewUserRepositoryWithExecutor(tx, ur.errTracker)

		user, err := txRepo.GetByID(ctx, userID)
		if err != nil {
			return err
		}

		err = txRepo.CheckEmailAvailability(ctx, user.Email)
		if err != nil {
			return domain.ErrEmailAlreadyVerified
		}

		_, err = tx.ExecContext(ctx, verifyEmailQuery, userID.String())
		if err != nil {
			err = fmt.Errorf("failed to update email verification status: %w", err)
			ur.errTracker.CaptureException(err)
			return err
		}

		user.IsEmailVerified = true
		returnedUser = user
		return nil
	})
}
