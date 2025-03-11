package app

import (
	"context"
	"go-starter/config"
	"go-starter/internal/adapters"
	"go-starter/internal/adapters/errtracker"
	"go-starter/internal/adapters/server/handlers"
	"go-starter/internal/domain/ports"
	"go-starter/internal/domain/services"
	"log/slog"
)

// Application is the main application struct.
type Application struct {
	ErrTracker ports.ErrTrackerAdapter
	Handlers   *handlers.Handlers
	Services   *services.Services
}

// New creates a new Application instance.
func New(ctx context.Context, cfg *config.Container) (*Application, func()) {
	errTracker := errtracker.New(cfg)

	apiAdapters := adapters.New(ctx, cfg, errTracker)
	apiServices := services.New(cfg, apiAdapters)
	apiHandlers := handlers.New(apiServices, errTracker)

	cleanup := createCleanupFunction(apiAdapters)

	app := &Application{
		ErrTracker: errTracker,
		Handlers:   apiHandlers,
		Services:   apiServices,
	}

	return app, cleanup
}

// createCleanupFunction creates a cleanup function for the application.
func createCleanupFunction(apiAdapters *adapters.Adapters) func() {
	return func() {
		slog.Info("cleaning app")
		err := apiAdapters.DB.Close()
		if err != nil {
			slog.Error("failed to close database", "error", err)
		}
		err = apiAdapters.CacheRepository.Close()
		if err != nil {
			slog.Error("failed to close cache repository", "error", err)
		}
	}
}
