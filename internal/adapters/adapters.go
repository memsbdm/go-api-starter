package adapters

import (
	"database/sql"
	"go-starter/internal/adapters/storage/postgres/repositories"
	"go-starter/internal/adapters/token"
	"go-starter/internal/domain/ports"
)

// Adapters holds all repository implementations for the application.
type Adapters struct {
	UserRepository  ports.UserRepository
	TokenRepository ports.TokenRepository
	CacheRepository ports.CacheRepository
	ErrTracker      ports.ErrorTracker
}

// New creates and initializes a new Adapters instance with the provided dependencies.
func New(db *sql.DB, timeGenerator ports.TimeGenerator, cache ports.CacheRepository, errTracker ports.ErrorTracker) *Adapters {
	return &Adapters{
		UserRepository:  repositories.NewUserRepository(db, errTracker),
		TokenRepository: token.NewTokenRepository(timeGenerator, errTracker),
		CacheRepository: cache,
		ErrTracker:      errTracker,
	}
}
