package postgres

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"go-starter/config"
	"time"
)

// New creates a postgres database instance
func New(c context.Context, config *config.DB) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.Addr)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)

	duration, err := time.ParseDuration(config.MaxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
