package postgres

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"go-starter/config"
	"log"
	"strconv"
	"sync"
	"time"
)

var (
	dbInstance *sql.DB
	mu         sync.RWMutex
)

// New creates a postgres database instance.
func New(ctx context.Context, config *config.DB) (*sql.DB, error) {
	mu.RLock()
	if dbInstance != nil {
		defer mu.RUnlock()
		return dbInstance, nil
	}
	mu.RUnlock()

	mu.Lock()
	defer mu.Unlock()

	// Double verification
	if dbInstance != nil {
		return dbInstance, nil
	}

	db, err := createConnection(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}

	dbInstance = db
	return dbInstance, nil
}

// createConnection establishes a new database connection with the given configuration.
func createConnection(c context.Context, config *config.DB) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.Addr)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxIdleTime(config.MaxIdleTime)

	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return db, nil
}

// Health checks the health status of the database.
func Health() map[string]string {
	mu.RLock()
	db := dbInstance
	mu.RUnlock()

	stats := make(map[string]string)

	if db == nil {
		stats["status"] = "down"
		stats["error"] = "database instance is nil"
		log.Fatal("database instance is nil")
		return stats
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Ping the database
	err := dbInstance.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	mu.RLock()
	dbStats := dbInstance.Stats()
	mu.RUnlock()

	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}
