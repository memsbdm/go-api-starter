package ports

import (
	"context"
	"time"
)

// CacheService is an interface for interacting with cache-related business logic.
type CacheService interface {
	// Set stores the value in the cache with a specified key and time-to-live (TTL).
	// Returns an error if the operation fails (e.g., if the cache is unreachable).
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error

	// Get retrieves the value associated with the specified key from the cache.
	// Returns the value as a byte slice and an error if the key is not found (domain.ErrCacheNotFound)
	// or if there are issues accessing the cache.
	Get(ctx context.Context, key string) ([]byte, error)

	// Delete removes the value associated with the specified key from the cache.
	// Returns an error if the operation fails (e.g., if there are issues accessing the cache).
	Delete(ctx context.Context, key string) error

	// DeleteByPrefix removes all values from the cache that match the given prefix.
	// Returns an error if the operation fails (e.g., if there are issues accessing the cache).
	DeleteByPrefix(ctx context.Context, prefix string) error

	// Close closes the connection to the cache server, ensuring that all resources are freed.
	// Returns an error if the operation fails (e.g., if there are issues closing the connection).
	Close() error
}

// CacheRepository is an interface for interacting with cache-related data.
type CacheRepository interface {
	// Set stores the value in the cache with a specified key and time-to-live (TTL).
	// Returns an error if the operation fails (e.g., if the cache is unreachable).
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error

	// Get retrieves the value associated with the specified key from the cache.
	// Returns the value as a byte slice and an error if the key is not found
	// or if there are issues accessing the cache.
	Get(ctx context.Context, key string) ([]byte, error)

	// Delete removes the value associated with the specified key from the cache.
	// Returns an error if the operation fails (e.g., if there are issues accessing the cache).
	Delete(ctx context.Context, key string) error

	// DeleteByPrefix removes all values from the cache that match the given prefix.
	// Returns an error if the operation fails (e.g., if there are issues accessing the cache).
	DeleteByPrefix(ctx context.Context, prefix string) error

	// Close closes the connection to the cache server, ensuring that all resources are freed.
	// Returns an error if the operation fails (e.g., if there are issues closing the connection).
	Close() error

	// Eval executes a Lua script on the cache server.
	// Returns the result of the script and an error if the operation fails (e.g., if there are issues executing the script).
	Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error)
}
