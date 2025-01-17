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
}

// New creates and initializes a new Adapters instance with the provided dependencies.
func New(db *sql.DB, timeGenerator ports.TimeGenerator, cache ports.CacheRepository) *Adapters {
	return &Adapters{
		UserRepository:  repositories.NewUserRepository(db),
		TokenRepository: token.NewTokenRepository(timeGenerator),
		CacheRepository: cache,
	}
}
