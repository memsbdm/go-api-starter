package adapters

import (
	"database/sql"
	"go-starter/internal/adapters/storage/database/repositories"
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
	FileUploadAdapter ports.FileUploadAdapter
}

// New creates and initializes a new Adapters instance with the provided dependencies.
func New(db *sql.DB, timeGenerator ports.TimeGenerator, cache ports.CacheRepository, errTracker ports.ErrTrackerAdapter, mailer ports.MailerAdapter, fileUpload ports.FileUploadAdapter) *Adapters {
	return &Adapters{
		UserRepository:    repositories.NewUserRepository(db, errTracker),
		TokenRepository:   token.NewTokenProvider(timeGenerator, errTracker),
		CacheRepository:   cache,
		ErrTrackerAdapter: errTracker,
		MailerAdapter:     mailer,
		FileUploadAdapter: fileUpload,
	}
}
