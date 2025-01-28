package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-starter/config"
	"go-starter/internal/adapters"
	"go-starter/internal/adapters/errortracker"
	"go-starter/internal/adapters/http"
	"go-starter/internal/adapters/http/handlers"
	"go-starter/internal/adapters/logger"
	"go-starter/internal/adapters/storage/postgres"
	"go-starter/internal/adapters/storage/postgres/repositories/mocks"
	"go-starter/internal/adapters/storage/redis"
	"go-starter/internal/adapters/timegen"
	"go-starter/internal/domain/ports"
	"go-starter/internal/domain/services"
	"log/slog"
	httpx "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title					Go Starter API
// @version					1.0
// @description				This is a simple starter API written in Go using net/http, Postgres database, and Redis cache.
//
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description				Type "Bearer" followed by a space and the access token.
func main() {
	if err := run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func run() error {
	// Load environment variables
	cfg := config.New()

	// Set logger
	logger.New(cfg.Application)

	slog.Info("Starting the application")

	// Init database
	ctx := context.Background()
	extServices, cleanup, err := initializeExternalServices(ctx, cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize external services: %w", err)
	}
	defer cleanup()

	// Dependency injection

	timeGenerator := timegen.NewRealTimeGenerator()

	apiAdapters := adapters.New(extServices.db, timeGenerator, extServices.cache, extServices.errTracker)
	apiServices := services.New(cfg, apiAdapters)
	apiHandlers := handlers.New(apiServices)

	// Init and start server
	srv := http.New(cfg.HTTP, apiHandlers, apiServices.TokenService, extServices.errTracker)

	done := make(chan bool, 1)
	go gracefulShutdown(srv, done, extServices.errTracker)

	err = srv.Serve()
	if err != nil && !errors.Is(err, httpx.ErrServerClosed) {
		err = fmt.Errorf("http server error: %s", err)
		extServices.errTracker.CaptureException(err)
		return err
	}
	extServices.errTracker.Flush(2 * time.Second)
	<-done
	slog.Info("Graceful shutdown complete.")
	return err
}

// externalServices holds connections to external services like database and cache.
// It encapsulates all external dependencies required by the application.
type externalServices struct {
	db         *sql.DB
	cache      ports.CacheRepository
	errTracker ports.ErrorTracker
}

// initializeExternalServices sets up connections to all external services .
// It returns the initialized services and a cleanup function to properly close all connections.
// The cleanup function should be deferred by the caller.
// If any service fails to initialize, it ensures proper cleanup of already initialized services.
func initializeExternalServices(ctx context.Context, cfg *config.Container) (*externalServices, func(), error) {
	// Init error tracker
	var errTracker ports.ErrorTracker
	errTracker = mocks.NewErrorTrackerMock(cfg.ErrTracker)
	if cfg.Application.Env == "production" {
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
		err := db.Close() // Clean up database connection if cache fails
		if err != nil {
			return nil, nil, fmt.Errorf("failed to close redis connection: %w", err)
		}
		return nil, nil, fmt.Errorf("failed to connect to cache service: %w", err)
	}
	slog.Info("Successfully connected to the cache service")

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
	}

	return &externalServices{
		db:         db,
		cache:      cache,
		errTracker: errTracker,
	}, cleanup, nil
}

// gracefulShutdown manages the graceful shutdown process of the HTTP server.
func gracefulShutdown(server *http.Server, done chan bool, errTracker ports.ErrorTracker) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	slog.Info("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		err = fmt.Errorf("server forced to shutdown with error: %v", err)
		errTracker.CaptureException(err)
		slog.Info(err.Error())
	}

	slog.Info("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}
