package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"go-starter/internal/domain/ports"
	"time"
)

const (
	// QueryTimeoutDuration specifies the timeout duration for database queries.
	QueryTimeoutDuration = time.Second * 5
)

// QueryExecutor represents a SQL query executor (DB or Transaction)
type QueryExecutor interface {
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

// withTx is a helper function that starts a transaction, executes a function, and commits or rolls back the transaction based on the function's return value.
// It also captures and logs any errors that occur during the transaction.
func withTx(db *sql.DB, ctx context.Context, errTracker ports.ErrTrackerAdapter, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		errTracker.CaptureException(err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			errTracker.CaptureException(rbErr)
			return fmt.Errorf("failed to rollback transaction: %v (original error: %w)", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		errTracker.CaptureException(err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
