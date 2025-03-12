package main

import (
	"context"
	"errors"
	"fmt"
	"go-starter/config"
	"go-starter/internal/adapters/logger"
	"go-starter/internal/adapters/server"
	"go-starter/internal/app"
	"go-starter/internal/domain/ports"
	"log/slog"
	"net/http"
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
	if err := run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error(err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func run() error {
	// Load environment variables
	cfg := config.New()
	logger.New(cfg.Application)

	slog.Info("starting the application")

	ctx := context.Background()
	app, cleanup := app.New(ctx, cfg)
	defer cleanup()

	handler := server.SetupRoutes(app.Handlers, app.Services, app.Adapters)
	srv := server.New(cfg.HTTP, handler)

	done := make(chan bool)
	go gracefulShutdown(srv, done, app.ErrTracker)

	err := srv.Serve()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		err = fmt.Errorf("http server error: %s", err)
		app.ErrTracker.CaptureException(err)
		return err
	}
	app.ErrTracker.Flush(2 * time.Second)
	<-done
	slog.Info("graceful shutdown complete")
	return err
}

// gracefulShutdown manages the graceful shutdown process of the HTTP server.
func gracefulShutdown(server *server.Server, done chan bool, errTrackerAdapter ports.ErrTrackerAdapter) {
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
		errTrackerAdapter.CaptureException(err)
		slog.Error(err.Error())
	}

	slog.Info("server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}
