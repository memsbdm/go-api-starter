package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
)

// UserRepository implements the ports.UserRepository interface and provides access to the database.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates and returns a new UserRepository instance.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// GetByID selects a user by their unique identifier from the database.
// Returns the user entity if found or an error if not found or any other issue occurs.
func (ur *UserRepository) GetByID(ctx context.Context, id entities.UserID) (*entities.User, error) {
	const query = `SELECT id, created_at, updated_at, name, username, is_email_verified FROM users WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	var (
		uuidStr string
		user    entities.User
	)
	err := ur.db.QueryRowContext(ctx, query, id.String()).Scan(&uuidStr, &user.CreatedAt, &user.UpdatedAt, &user.Name, &user.Username, &user.IsEmailVerified)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, fmt.Errorf("%w: id=%s", domain.ErrUserNotFound, id)
		default:
			return nil, fmt.Errorf("failed to get user: %w", err)
		}
	}

	userID, err := entities.ParseUserID(uuidStr)
	if err != nil {
		return nil, err
	}
	user.ID = userID

	return &user, nil
}

// GetByUsername selects a user by their username from the database.
// Returns the user entity if found or an error if not found or any other issue occurs.
func (ur *UserRepository) GetByUsername(ctx context.Context, username string) (*entities.User, error) {
	query := `SELECT id, created_at, updated_at, name, username, password FROM users WHERE username = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	user := &entities.User{}
	var uuidStr string
	err := ur.db.QueryRowContext(ctx, query, username).Scan(&uuidStr, &user.CreatedAt, &user.UpdatedAt, &user.Name, &user.Username, &user.Password)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, domain.ErrUserNotFound
		default:
			return nil, err
		}
	}

	parsedID, err := entities.ParseUserID(uuidStr)
	if err != nil {
		return nil, err
	}
	user.ID = parsedID

	return user, nil
}

// Create inserts a new user into the database.
// Returns the created user or an error if the operation fails (e.g., due to a database constraint violation).
func (ur *UserRepository) Create(ctx context.Context, user *entities.User) (*entities.User, error) {
	const query = `INSERT INTO users (name, username, password) VALUES ($1, $2, $3) RETURNING id, name, username, created_at, updated_at`
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
	).Scan(
		&uuidStr,
		&createdUser.Name,
		&createdUser.Username,
		&createdUser.CreatedAt,
		&createdUser.UpdatedAt,
	)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23505": // Code unique_violation
				if pqErr.Constraint == "users_username_key" {
					return nil, fmt.Errorf("%w: username=%s", domain.ErrUsernameAlreadyTaken, user.Username)
				}
			}
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	parsedID, err := entities.ParseUserID(uuidStr)
	if err != nil {
		return nil, err
	}
	user.ID = parsedID

	return user, nil
}

// UpdatePassword updates a user password.
// Returns an error if the update fails.
func (ur *UserRepository) UpdatePassword(ctx context.Context, userID entities.UserID, newPassword string) error {
	const query = `UPDATE users SET password = $1 WHERE id = $2 `
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := ur.db.ExecContext(ctx, query, newPassword, userID.String())
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}
