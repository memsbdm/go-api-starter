package main

import (
	"context"
	"database/sql"
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
	// Load environment variables
	cfg := config.New()

	ctx := context.Background()
	db, err := postgres.New(ctx, cfg.DB)
	if err != nil {
		slog.Error("failed to connect to database")
		os.Exit(1)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			slog.Error("failed to close database connection")
		} else {
			slog.Info("Successfully closed database connection")
		}
	}()
	slog.Info("Successfully connected to the database")

	seedUsers(db)
}

func seedUsers(db *sql.DB) {
	userRepo := repositories.NewUserRepository(db)
	timeGenerator := timegen.NewRealTimeGenerator()
	cacheService := redis.NewCacheMock(timeGenerator)
	userService := services.NewUserService(userRepo, cacheService)
	slog.Info("Seeding users...")
	seed.Users(userService)
}
