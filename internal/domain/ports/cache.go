package ports

import (
	"context"
	"time"
)

// CacheService is an interface for interacting with cache-related business logic
type CacheService interface {
	// GenerateCacheKey generates a cache key based on the input parameters
	GenerateCacheKey(prefix string, params any) string
	// GenerateCacheKeyParams generates a cache params based on the input parameters
	GenerateCacheKeyParams(params ...any) string
	// Serialize marshals the input data into an array of bytes
	Serialize(data any) ([]byte, error)
	// Deserialize unmarshals the input data into the output interface
	Deserialize(data []byte, output any) error
	// Set stores the value in the cache
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	// Get retrieves the value from the cache
	Get(ctx context.Context, key string) ([]byte, error)
	// Delete removes the value from the cache
	Delete(ctx context.Context, key string) error
	// DeleteByPrefix removes the value from the cache with the given prefix
	DeleteByPrefix(ctx context.Context, prefix string) error
	// Close closes the connection to the cache server
	Close() error
}

// CacheRepository is an interface for interacting with cache-related data
type CacheRepository interface {
	// Set stores the value in the cache
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	// Get retrieves the value from the cache
	Get(ctx context.Context, key string) ([]byte, error)
	// Delete removes the value from the cache
	Delete(ctx context.Context, key string) error
	// DeleteByPrefix removes the value from the cache with the given prefix
	DeleteByPrefix(ctx context.Context, prefix string) error
	// Close closes the connection to the cache server
	Close() error
}
