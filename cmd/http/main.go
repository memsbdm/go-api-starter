package main

import (
	"context"
	"errors"
	"fmt"
	"go-starter/config"
	"go-starter/internal/adapters"
	"go-starter/internal/adapters/http"
	"go-starter/internal/adapters/http/handlers"
	"go-starter/internal/adapters/logger"
	"go-starter/internal/adapters/timegen"
	"go-starter/internal/domain/services"
	"log/slog"
	httpx "net/http"
	"os"
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
		if errors.Is(err, httpx.ErrServerClosed) {
			os.Exit(0)
		}
		slog.Error(err.Error())
		os.Exit(1)
	}
	os.Exit(0)
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

	timeGenerator := timegen.NewTimeGenerator()

	apiAdapters := adapters.New(extServices.db, timeGenerator, extServices.cache, extServices.errTracker, extServices.mailer)
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
