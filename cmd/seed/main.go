package main

import (
	"context"
	"database/sql"
	"fmt"
	"go-starter/config"
	"go-starter/internal/adapters/storage/postgres"
	"go-starter/internal/adapters/storage/postgres/repositories"
	"go-starter/internal/adapters/storage/postgres/seed"
	"go-starter/internal/adapters/storage/redis"
	"go-starter/internal/adapters/timegen"
	"go-starter/internal/domain/services"
	"log/slog"
	"os"
)

func main() {
	if err := run(); err != nil {
		slog.Error("application error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	// Load environment variables
	cfg := config.New()

	db, err := initDatabase(ctx, cfg)
	if err != nil {
		return fmt.Errorf("initializing database: %w", err)
	}
	defer closeDB(db)

	// Seeders
	if err := seedUsers(ctx, db); err != nil {
		return fmt.Errorf("seeding users: %w", err)
	}

	slog.Info("Seeding completed successfully")
	return nil
}

func initDatabase(ctx context.Context, cfg *config.Container) (*sql.DB, error) {
	db, err := postgres.New(ctx, cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}
	slog.Info("Successfully connected to the database")
	return db, nil
}

func closeDB(db *sql.DB) {
	if err := db.Close(); err != nil {
		slog.Error("failed to close database connection", "error", err)
		return
	}
	slog.Info("Successfully closed database connection")
}

func seedUsers(ctx context.Context, db *sql.DB) error {
	// Initialize dependencies
	userRepo := repositories.NewUserRepository(db)
	timeGenerator := timegen.NewRealTimeGenerator()
	cacheService := redis.NewCacheMock(timeGenerator)
	userService := services.NewUserService(userRepo, cacheService)

	// Configure and run user generator
	slog.Info("Starting user seeding process")
	userGenerator := seed.NewUserGenerator(userService)

	opts := seed.GenerateUsersOptions{
		Count:     100,
		BatchSize: 30,
	}

	if err := userGenerator.GenerateUsers(ctx, opts); err != nil {
		return fmt.Errorf("generating users: %w", err)
	}

	return nil
}
