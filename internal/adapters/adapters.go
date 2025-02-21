package adapters

import (
	"database/sql"
	"go-starter/internal/adapters/storage/postgres/repositories"
	"go-starter/internal/adapters/token"
	"go-starter/internal/domain/ports"
)

// Adapters holds all repository implementations for the application.
type Adapters struct {
	UserRepository    ports.UserRepository
	TokenRepository   ports.TokenProvider
	CacheRepository   ports.CacheRepository
	ErrTrackerAdapter ports.ErrTrackerAdapter
	MailerAdapter     ports.MailerAdapter
}

// New creates and initializes a new Adapters instance with the provided dependencies.
func New(db *sql.DB, timeGenerator ports.TimeGenerator, cacheRepo ports.CacheRepository, errTrackerAdapter ports.ErrTrackerAdapter, mailerAdapter ports.MailerAdapter) *Adapters {
	return &Adapters{
		UserRepository:    repositories.NewUserRepository(db, errTrackerAdapter),
		TokenRepository:   token.NewJWTProvider(timeGenerator, errTrackerAdapter),
		CacheRepository:   cacheRepo,
		ErrTrackerAdapter: errTrackerAdapter,
		MailerAdapter:     mailerAdapter,
	}
}
