package main

import (
	"context"
	"database/sql"
	"fmt"
	"go-starter/config"
	"go-starter/internal/adapters/errortracker"
	"go-starter/internal/adapters/mailer"
	"go-starter/internal/adapters/mocks"
	"go-starter/internal/adapters/storage/postgres"
	"go-starter/internal/adapters/storage/redis"
	"go-starter/internal/domain/ports"
	"log/slog"
)

// externalServices holds connections to external services like database and cache.
// It encapsulates all external dependencies required by the application.
type externalServices struct {
	db         *sql.DB
	cache      ports.CacheRepository
	errTracker ports.ErrorTracker
	mailer     ports.MailerRepository
}

// initializeExternalServices sets up connections to all external services .
// It returns the initialized services and a cleanup function to properly close all connections.
// The cleanup function should be deferred by the caller.
// If any service fails to initialize, it ensures proper cleanup of already initialized services.
func initializeExternalServices(ctx context.Context, cfg *config.Container) (*externalServices, func(), error) {
	// Init error tracker
	var errTracker ports.ErrorTracker
	errTracker = mocks.NewErrorTrackerMock(cfg.ErrTracker)
	if cfg.Application.Env == config.EnvProduction {
		errTracker = errortracker.NewSentryErrorTracker(cfg.ErrTracker)
	}

	// Init database
	db, err := postgres.New(ctx, cfg.DB)
	if err != nil {
		err = fmt.Errorf("failed to connect to database: %w", err)
		errTracker.CaptureException(err)
		return nil, nil, err
	}
	slog.Info("Successfully connected to the database")

	// Init cache service
	cache, err := redis.New(ctx, cfg.Redis)
	if err != nil {
		errTracker.CaptureException(err)
		return nil, nil, fmt.Errorf("failed to connect to cache service: %w", err)
	}
	slog.Info("Successfully connected to the cache service")

	// Init mailer
	smtp, err := mailer.New(cfg.Mailer)
	if err != nil {
		errTracker.CaptureException(err)
		return nil, nil, fmt.Errorf("failed to initialize mailer: %w", err)
	}

	cleanup := func() {
		if err := db.Close(); err != nil {
			err = fmt.Errorf("failed to close database connection: %w", err)
			errTracker.CaptureException(err)
			slog.Error(err.Error())
		}
		if err := cache.Close(); err != nil {
			err = fmt.Errorf("failed to close cache connection: %w", err)
			errTracker.CaptureException(err)
			slog.Error(err.Error())
		}
		if err := smtp.Close(); err != nil {
			err = fmt.Errorf("failed to close mailer connection: %w", err)
			errTracker.CaptureException(err)
			slog.Error(err.Error())
		}
	}

	return &externalServices{
		db:         db,
		cache:      cache,
		errTracker: errTracker,
		mailer:     smtp,
	}, cleanup, nil
}
