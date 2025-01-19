package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-starter/config"
	"go-starter/internal/adapters"
	"go-starter/internal/adapters/http"
	"go-starter/internal/adapters/http/handlers"
	"go-starter/internal/adapters/logger"
	"go-starter/internal/adapters/storage/postgres"
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
// @description				This is a simple starter API written in Go using net/http, PostgresSQL database, and Redis cache.
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

	apiAdapters := adapters.New(extServices.db, timeGenerator, extServices.cache)
	apiServices := services.New(cfg, apiAdapters)
	apiHandlers := handlers.New(apiServices)

	// Init and start server
	srv := http.New(cfg.HTTP, apiHandlers, apiServices.TokenService)

	done := make(chan bool, 1)
	go gracefulShutdown(srv, done)

	err = srv.Serve()
	if err != nil && !errors.Is(err, httpx.ErrServerClosed) {
		return fmt.Errorf("http server error: %s", err)
	}
	<-done
	slog.Info("Graceful shutdown complete.")
	return err
}

// externalServices holds connections to external services like database and cache.
// It encapsulates all external dependencies required by the application.
type externalServices struct {
	db    *sql.DB
	cache ports.CacheRepository
}

// initializeExternalServices sets up connections to all external services .
// It returns the initialized services and a cleanup function to properly close all connections.
// The cleanup function should be deferred by the caller.
// If any service fails to initialize, it ensures proper cleanup of already initialized services.
func initializeExternalServices(ctx context.Context, cfg *config.Container) (*externalServices, func(), error) {
	// Init database
	db, err := postgres.New(ctx, cfg.DB)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	slog.Info("Successfully connected to the database")

	// Init cache service
	cache, err := redis.New(ctx, cfg.Redis)
	if err != nil {
		err := db.Close() // Clean up database connection if cache fails
		if err != nil {
			return nil, nil, fmt.Errorf("failed to close redis connection: %w", err)
		}
		return nil, nil, fmt.Errorf("failed to connect to cache service: %w", err)
	}
	slog.Info("Successfully connected to the cache service")

	cleanup := func() {
		if err := db.Close(); err != nil {
			slog.Error("failed to close database connection: " + err.Error())
		}
		if err := cache.Close(); err != nil {
			slog.Error("failed to close cache connection: " + err.Error())
		}
	}

	return &externalServices{
		db:    db,
		cache: cache,
	}, cleanup, nil
}

// gracefulShutdown manages the graceful shutdown process of the HTTP server.
func gracefulShutdown(server *http.Server, done chan bool) {
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
		slog.Info(fmt.Sprintf("Server forced to shutdown with error: %v", err))
	}

	slog.Info("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}
