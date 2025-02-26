package redis

import (
	"context"
	"errors"
	"fmt"
	"go-starter/config"
	"go-starter/internal/domain"
	"go-starter/internal/domain/ports"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis implements the ports.CacheRepository interface and provides access to the Redis library.
type Redis struct {
	client *redis.Client
}

// New creates a new instance of Redis.
func New(ctx context.Context, redisCfg *config.Redis) (ports.CacheRepository, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
		// TLSConfig: &tls.Config{
		// 	MinVersion: tls.VersionTLS12, // TODO: Uncomment this when we have a valid certificate
		// },
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &Redis{client}, nil
}

// Set stores the value in the cache with a specified key and time-to-live (TTL).
// Returns an error if the operation fails (e.g., if the cache is unreachable).
func (r *Redis) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

// Get retrieves the value associated with the specified key from the cache.
// Returns the value as a byte slice and an error if the key is not found
// or if there are issues accessing the cache.
func (r *Redis) Get(ctx context.Context, key string) ([]byte, error) {
	res, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, domain.ErrCacheNotFound
		}
		return nil, err
	}
	return []byte(res), nil
}

// Delete removes the value associated with the specified key from the cache.
// Returns an error if the operation fails (e.g., if there are issues accessing the cache).
func (r *Redis) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// DeleteByPrefix removes all values from the cache that match the given prefix.
// Returns an error if the operation fails (e.g., if there are issues accessing the cache).
func (r *Redis) DeleteByPrefix(ctx context.Context, prefix string) error {
	var cursor uint64
	prefix = prefix + "*"
	for {
		var err error
		keys, cursor, err := r.client.Scan(ctx, cursor, prefix, 100).Result()
		if err != nil {
			return fmt.Errorf("scan error: %w", err)
		}

		for _, key := range keys {
			err := r.client.Del(ctx, key).Err()
			if err != nil {
				return fmt.Errorf("delete error: %w", err)
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}

// Close closes the connection to the cache server, ensuring that all resources are freed.
// Returns an error if the operation fails (e.g., if there are issues closing the connection).
func (r *Redis) Close() error {
	return r.client.Close()
}
