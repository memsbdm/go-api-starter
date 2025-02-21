package main

import (
	"context"
	"fmt"
	"go-starter/internal/adapters/http"
	"go-starter/internal/domain/ports"
	"log/slog"
	"os/signal"
	"syscall"
	"time"
)

// gracefulShutdown manages the graceful shutdown process of the HTTP server.
func gracefulShutdown(server *http.Server, done chan bool, errTrackerAdapter ports.ErrTrackerAdapter) {
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

	slog.Info("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}
