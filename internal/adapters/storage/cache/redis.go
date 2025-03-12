package cache

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
	client     *redis.Client
	errTracker ports.ErrTrackerAdapter
}

// New creates a new instance of Redis.
func New(ctx context.Context, redisCfg *config.Redis, errTracker ports.ErrTrackerAdapter) (ports.CacheRepository, error) {
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
		errTracker.CaptureException(fmt.Errorf("failed to ping redis: %w", err))
		return nil, err
	}

	return &Redis{client: client, errTracker: errTracker}, nil
}

// Set stores the value in the cache with a specified key and time-to-live (TTL).
// Returns an error if the operation fails (e.g., if the cache is unreachable).
func (r *Redis) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	err := r.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		r.errTracker.CaptureException(fmt.Errorf("failed to set value in redis: %w", err))
		return err
	}
	return nil
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
		r.errTracker.CaptureException(fmt.Errorf("failed to get value from redis: %w", err))
		return nil, err
	}
	return []byte(res), nil
}

// Delete removes the value associated with the specified key from the cache.
// Returns an error if the operation fails (e.g., if there are issues accessing the cache).
func (r *Redis) Delete(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		r.errTracker.CaptureException(fmt.Errorf("failed to delete value from redis: %w", err))
		return err
	}
	return nil
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
			r.errTracker.CaptureException(fmt.Errorf("failed to scan values from redis: %w", err))
			return err
		}

		for _, key := range keys {
			err := r.client.Del(ctx, key).Err()
			if err != nil {
				r.errTracker.CaptureException(fmt.Errorf("failed to delete value from redis: %w", err))
				return err
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
	err := r.client.Close()
	if err != nil {
		r.errTracker.CaptureException(fmt.Errorf("failed to close redis connection: %w", err))
		return err
	}
	return nil
}

// Eval executes a Lua script in Redis.
// The script is executed atomically and can access keys and arguments passed to it.
func (r *Redis) Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error) {
	result, err := r.client.Eval(ctx, script, keys, args...).Result()
	if err != nil {
		r.errTracker.CaptureException(fmt.Errorf("failed to execute Lua script in redis: %w", err))
		return nil, err
	}
	return result, nil
}
