package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"go-starter/internal/domain/ports"

	"github.com/google/uuid"

	"github.com/lib/pq"
)

// UserRepository implements the ports.UserRepository interface and provides access to the database.
type UserRepository struct {
	db         *sql.DB
	errTracker ports.ErrTrackerAdapter
}

// NewUserRepository creates and returns a new UserRepository instance.
func NewUserRepository(db *sql.DB, errTracker ports.ErrTrackerAdapter) *UserRepository {
	return &UserRepository{
		db:         db,
		errTracker: errTracker,
	}
}

// GetByID selects a user by their unique identifier from the database.
// Returns the user entity if found or an error if not found or any other issue occurs.
func (ur *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	const query = `SELECT id, created_at, updated_at, name, username, email, is_email_verified FROM users WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	var (
		uuidStr string
		user    entities.User
	)
	err := ur.db.QueryRowContext(ctx, query, id.String()).Scan(&uuidStr, &user.CreatedAt, &user.UpdatedAt, &user.Name, &user.Username, &user.Email, &user.IsEmailVerified)
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
	query := `SELECT id, created_at, updated_at, name, username, password, email, is_email_verified FROM users WHERE username = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	user := &entities.User{}
	var uuidStr string
	err := ur.db.QueryRowContext(ctx, query, username).Scan(&uuidStr, &user.CreatedAt, &user.UpdatedAt, &user.Name, &user.Username, &user.Password, &user.Email, &user.IsEmailVerified)

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

// CheckEmailAvailability checks if an email is available for registration.
// Returns an error if the email is already taken.
func (ur *UserRepository) CheckEmailAvailability(ctx context.Context, email string) error {
	const query = `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND is_email_verified = true)`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	var exists bool
	if err := ur.db.QueryRowContext(ctx, query, email).Scan(&exists); err != nil {
		err = fmt.Errorf("failed to check email verification status: %w", err)
		ur.errTracker.CaptureException(err)
		return err
	}

	if exists {
		return domain.ErrEmailAlreadyTaken
	}

	return nil
}

// Create inserts a new user into the database.
// Returns the created user or an error if the operation fails (e.g., due to a database constraint violation).
func (ur *UserRepository) Create(ctx context.Context, user *entities.User) (*entities.User, error) {
	const query = `INSERT INTO users (name, username, password, email) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at, name, username, email, is_email_verified`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var (
		uuidStr     string
		createdUser = *user
	)

	err := ur.db.QueryRowContext(
		ctx,
		query,
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
	)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23505": // Code unique_violation
				if pqErr.Constraint == "users_username_key" {
					return nil, domain.ErrUsernameAlreadyTaken
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
func (ur *UserRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	const query = `UPDATE users SET password = $1 WHERE id = $2 `
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := ur.db.ExecContext(ctx, query, newPassword, userID.String())
	if err != nil {
		err = fmt.Errorf("failed to update user password for user %s: %w", userID.String(), err)
		ur.errTracker.CaptureException(err)
		return err
	}

	return nil
}

// VerifyEmail updates the email verification status of a user.
func (ur *UserRepository) VerifyEmail(ctx context.Context, userID uuid.UUID) (*entities.User, error) {
	const query = `UPDATE users SET is_email_verified = true WHERE id = $1 `
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := ur.db.ExecContext(ctx, query, userID.String())
	if err != nil {
		err = fmt.Errorf("failed to update users email verification status for user %s: %w", userID.String(), err)
		ur.errTracker.CaptureException(err)
		return nil, err
	}

	return ur.GetByID(ctx, userID)
}
