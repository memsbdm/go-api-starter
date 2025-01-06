package main

import (
	"context"
	"go-starter/config"
	"go-starter/internal/adapters/http"
	"go-starter/internal/adapters/logger"
	"go-starter/internal/adapters/storage/postgres"
	"go-starter/internal/adapters/storage/postgres/repositories"
	"go-starter/internal/adapters/storage/redis"
	"go-starter/internal/domain/services"
	"log/slog"
	"os"
)

func main() {
	// Load environment variables
	cfg := config.New()

	// Set logger
	logger.New(cfg.Application)

	slog.Info("Starting the application")

	// Init database
	ctx := context.Background()
	db, err := postgres.New(ctx, cfg.DB)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			slog.Error(err.Error())
		}
	}()
	slog.Info("Successfully connected to the database")

	// Init cache service
	cache, err := redis.New(ctx, cfg.Redis)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	defer func() {
		err := cache.Close()
		if err != nil {
			slog.Error(err.Error())
		}
	}()
	slog.Info("Successfully connected to the cache service")

	// Dependency injection
	// Health
	healthHandler := http.NewHealthHandler()

	// Cache
	cacheService := services.NewCacheService(cache)

	// User
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo, cacheService)
	userHandler := http.NewUserHandler(userService)

	// Init and start server
	srv := http.New(cfg.HTTP, *healthHandler, *userHandler)
	srv.Serve()
}
