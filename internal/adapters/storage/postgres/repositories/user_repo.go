package repositories

import (
	"context"
	"database/sql"
	"errors"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
)

// UserRepository implements ports.UserRepository interface and provides access to the database
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// GetByID gets a user by ID from the database
func (ur *UserRepository) GetByID(ctx context.Context, id int) (*entities.User, error) {
	query := `SELECT id, username FROM users WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	user := &entities.User{}
	err := ur.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, domain.ErrUserNotFound
		default:
			return nil, err
		}
	}
	return user, nil
}

// Create creates a new user in the database
func (ur *UserRepository) Create(ctx context.Context, user *entities.User) (*entities.User, error) {
	query := `INSERT INTO users (username) VALUES ($1) RETURNING id, username`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := ur.db.QueryRowContext(ctx, query, user.Username).Scan(&user.ID, &user.Username)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return nil, domain.ErrUserUsernameAlreadyExists
		default:
			return nil, err
		}
	}
	return user, nil
}
