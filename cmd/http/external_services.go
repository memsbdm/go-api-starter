package main

import (
	"context"
	"database/sql"
	"fmt"
	"go-starter/config"
	"go-starter/internal/adapters/errtracker"
	"go-starter/internal/adapters/mailer"
	"go-starter/internal/adapters/storage/cache"
	"go-starter/internal/adapters/storage/fileupload"
	"go-starter/internal/adapters/storage/postgres"
	"go-starter/internal/domain/ports"
	"log/slog"
)

// externalServices holds connections to external services like database and cache.
// It encapsulates all external dependencies required by the application.
type externalServices struct {
	db         *sql.DB
	cache      ports.CacheRepository
	errTracker ports.ErrTrackerAdapter
	mailer     ports.MailerAdapter
	fileUpload ports.FileUploadAdapter
}

// initializeExternalServices sets up connections to all external services .
// It returns the initialized services and a cleanup function to properly close all connections.
// The cleanup function should be deferred by the caller.
// If any service fails to initialize, it ensures proper cleanup of already initialized services.
func initializeExternalServices(ctx context.Context, cfg *config.Container) (*externalServices, func(), error) {
	errTracker := initializeErrTracker(cfg)

	db, err := initializeDatabase(ctx, cfg, errTracker)
	if err != nil {
		return nil, nil, err
	}

	cache, err := initializeCache(ctx, cfg, errTracker)
	if err != nil {
		_ = db.Close()
		return nil, nil, err
	}

	ses, err := initializeMailer(cfg, errTracker)
	if err != nil {
		_ = db.Close()
		_ = cache.Close()
		return nil, nil, err
	}

	s3, err := initializeFileUpload(cfg.FileUpload, errTracker)
	if err != nil {
		_ = db.Close()
		_ = cache.Close()
		return nil, nil, err
	}

	cleanup := createCleanupFunction(db, cache, errTracker)

	return &externalServices{
		db:         db,
		cache:      cache,
		errTracker: errTracker,
		mailer:     ses,
		fileUpload: s3,
	}, cleanup, nil
}

func initializeErrTracker(cfg *config.Container) ports.ErrTrackerAdapter {
	var errTracker ports.ErrTrackerAdapter
	errTracker = errtracker.NewErrTrackerAdapterMock()
	if cfg.Application.Env != config.EnvDevelopment {
		errTracker = errtracker.NewSentryAdapter(cfg)
	}
	return errTracker
}

func initializeDatabase(ctx context.Context, cfg *config.Container, errTracker ports.ErrTrackerAdapter) (*sql.DB, error) {
	db, err := postgres.New(ctx, cfg.DB, errTracker)
	if err != nil {
		return nil, err
	}
	slog.Info("Successfully connected to the database")
	return db, nil
}

func initializeCache(ctx context.Context, cfg *config.Container, errTracker ports.ErrTrackerAdapter) (ports.CacheRepository, error) {
	cache, err := cache.New(ctx, cfg.Redis, errTracker)
	if err != nil {
		return nil, err
	}
	slog.Info("Successfully connected to the cache service")
	return cache, nil
}

func initializeMailer(cfg *config.Container, errTracker ports.ErrTrackerAdapter) (ports.MailerAdapter, error) {
	adapter, err := mailer.NewSESAdapter(cfg.Mailer, errTracker)
	if err != nil {
		return nil, err
	}
	return adapter, nil
}

func initializeFileUpload(fileUploadCfg *config.FileUpload, errTracker ports.ErrTrackerAdapter) (ports.FileUploadAdapter, error) {
	adapter, err := fileupload.NewS3Adapter(fileUploadCfg, errTracker)
	if err != nil {
		return nil, err
	}
	return adapter, nil
}

func createCleanupFunction(db *sql.DB, cache ports.CacheRepository, errTracker ports.ErrTrackerAdapter) func() {
	return func() {
		if err := db.Close(); err != nil {
			err = fmt.Errorf("failed to close database connection: %w", err)
			errTracker.CaptureException(err)
			slog.Error(err.Error())
		}
		if err := cache.Close(); err != nil {
			slog.Error(err.Error())
		}
	}
}
