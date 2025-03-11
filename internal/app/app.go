package app

import (
	"context"
	"go-starter/config"
	"go-starter/internal/adapters"
	"go-starter/internal/adapters/errtracker"
	"go-starter/internal/adapters/http/handlers"
	"go-starter/internal/domain/ports"
	"go-starter/internal/domain/services"
	"log/slog"
)

// TODO
type Application struct {
	ErrTracker ports.ErrTrackerAdapter
	Handlers   *handlers.Handlers
	Services   *services.Services
}

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

// TODO
func createCleanupFunction(apiAdapters *adapters.Adapters) func() {
	return func() {
		slog.Info("cleaning app")
		apiAdapters.DB.Close()
		apiAdapters.CacheRepository.Close()
		// err
	}
}
