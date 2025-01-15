package repositories

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
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
	query := `SELECT id, username FROM users WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	user := &entities.User{}
	var uuidStr string
	err := ur.db.QueryRowContext(ctx, query, id.String()).Scan(&uuidStr, &user.Username)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, domain.ErrUserNotFound
		default:
			return nil, err
		}
	}

	parsedID, err := uuid.Parse(uuidStr)
	if err != nil {
		return nil, err
	}
	user.ID = entities.UserID(parsedID)

	return user, nil
}

// GetByUsername selects a user by their username from the database.
// Returns the user entity if found or an error if not found or any other issue occurs.
func (ur *UserRepository) GetByUsername(ctx context.Context, username string) (*entities.User, error) {
	query := `SELECT id, username, password FROM users WHERE username = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	user := &entities.User{}
	var uuidStr string
	err := ur.db.QueryRowContext(ctx, query, username).Scan(&uuidStr, &user.Username, &user.Password)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, domain.ErrUserNotFound
		default:
			return nil, err
		}
	}

	parsedID, err := uuid.Parse(uuidStr)
	if err != nil {
		return nil, err
	}
	user.ID = entities.UserID(parsedID)

	return user, nil
}

// Create inserts a new user into the database.
// Returns the created user or an error if the operation fails (e.g., due to a database constraint violation).
func (ur *UserRepository) Create(ctx context.Context, user *entities.User) (*entities.User, error) {
	query := `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id, username, created_at, updated_at`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var uuidStr string
	err := ur.db.QueryRowContext(ctx, query, user.Username, user.Password).Scan(&uuidStr, &user.Username, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return nil, domain.ErrUserUsernameAlreadyExists
		default:
			return nil, err
		}
	}
	parsedID, err := uuid.Parse(uuidStr)
	if err != nil {
		return nil, err
	}
	user.ID = entities.UserID(parsedID)

	return user, nil
}
