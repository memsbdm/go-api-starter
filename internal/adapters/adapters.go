package adapters

import (
	"context"
	"database/sql"
	"go-starter/config"
	"go-starter/internal/adapters/mailer"
	"go-starter/internal/adapters/storage/cache"
	"go-starter/internal/adapters/storage/database"
	"go-starter/internal/adapters/storage/database/migrations"
	"go-starter/internal/adapters/storage/database/repositories"
	"go-starter/internal/adapters/storage/fileupload"
	"go-starter/internal/adapters/timegen"
	"go-starter/internal/adapters/token"
	"go-starter/internal/domain/ports"
)

// Adapters holds all repository implementations for the application.
type Adapters struct {
	TimeGenerator     ports.TimeGenerator
	DB                *sql.DB
	UserRepository    ports.UserRepository
	TokenRepository   ports.TokenProvider
	CacheRepository   ports.CacheRepository
	ErrTrackerAdapter ports.ErrTrackerAdapter
	MailerAdapter     ports.MailerAdapter
	FileUploadAdapter ports.FileUploadAdapter
}

// New creates and initializes a new Adapters instance with the provided dependencies.
func New(ctx context.Context, cfg *config.Container, errTracker ports.ErrTrackerAdapter) *Adapters {
	timeGenerator := timegen.NewTimeGenerator()
	db := initializeDatabaseAndMigrate(ctx, cfg.DB, errTracker)

	return &Adapters{
		TimeGenerator:     timeGenerator,
		DB:                db,
		UserRepository:    repositories.NewUserRepository(db, errTracker),
		TokenRepository:   token.NewTokenProvider(timeGenerator, errTracker),
		CacheRepository:   initializeCache(ctx, cfg.Redis, errTracker),
		ErrTrackerAdapter: errTracker,
		MailerAdapter:     initializeMailer(cfg.Mailer, errTracker),
		FileUploadAdapter: initializeFileUpload(cfg.FileUpload, errTracker),
	}
}

// TODO
func initializeDatabaseAndMigrate(ctx context.Context, dbCfg *config.DB, errTracker ports.ErrTrackerAdapter) *sql.DB {
	db, err := database.New(ctx, dbCfg, errTracker)
	if err != nil {
		errTracker.CaptureException(err)
		panic(err)
	}

	err = database.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		errTracker.CaptureException(err)
		panic(err)
	}

	return db
}

func initializeCache(ctx context.Context, cacheCfg *config.Redis, errTracker ports.ErrTrackerAdapter) ports.CacheRepository {
	cache, err := cache.New(ctx, cacheCfg, errTracker)
	if err != nil {
		errTracker.CaptureException(err)
		panic(err)
	}
	return cache
}

func initializeMailer(mailerCfg *config.Mailer, errTracker ports.ErrTrackerAdapter) ports.MailerAdapter {
	mailer, err := mailer.NewSESAdapter(mailerCfg, errTracker)
	if err != nil {
		errTracker.CaptureException(err)
		panic(err)
	}
	return mailer
}

func initializeFileUpload(fileUploadCfg *config.FileUpload, errTracker ports.ErrTrackerAdapter) ports.FileUploadAdapter {
	fileUpload, err := fileupload.NewS3Adapter(fileUploadCfg, errTracker)
	if err != nil {
		errTracker.CaptureException(err)
		panic(err)
	}
	return fileUpload
}
