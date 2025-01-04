package main

import (
	"context"
	"go-starter/config"
	"go-starter/internal/adapters/http"
	"go-starter/internal/adapters/logger"
	"go-starter/internal/adapters/storage/postgres"
	"go-starter/internal/adapters/storage/postgres/repositories"
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
		slog.Error("Failed to connect to database")
		os.Exit(1)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			slog.Error("Failed to close database connection")
		}
	}()
	slog.Info("Successfully connected to the database")

	// Dependency injection
	// User
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := http.NewUserHandler(userService)

	// Init and start server
	srv := http.New(cfg.HTTP, *userHandler)
	srv.Serve()
}
