package main

import (
	"context"
	"database/sql"
	"fmt"
	"go-starter/config"
	"go-starter/internal/adapters/storage/postgres"
	"go-starter/internal/adapters/storage/postgres/seed"
	"log/slog"
	"os"
)

func main() {
	if err := run(); err != nil {
		slog.Error("application error", "error", err)
		os.Exit(1)
	}
	os.Exit(0)
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
	if err := seed.SeedUsers(ctx, cfg, db); err != nil {
		return fmt.Errorf("seeding users: %w", err)
	}

	slog.Info("seeding completed successfully")
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
	slog.Info("successfully closed database connection")
}
