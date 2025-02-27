package services

import (
	"context"
	"errors"
	"go-starter/internal/domain"
	"go-starter/internal/domain/ports"
	"time"
)

// CacheService implements ports.CacheService interface and provides access to the cache repository
type CacheService struct {
	repo ports.CacheRepository
}

// NewCacheService creates a new cache service instance
func NewCacheService(repo ports.CacheRepository) *CacheService {
	return &CacheService{
		repo: repo,
	}
}

// Set stores the value in the cache with a specified key and time-to-live (TTL).
// Returns an error if the operation fails (e.g., if the cache is unreachable).
func (cs *CacheService) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	err := cs.repo.Set(ctx, key, value, ttl)
	if err != nil {
		return domain.ErrInternal
	}
	return nil
}

// Get retrieves the value associated with the specified key from the cache.
// Returns the value as a byte slice and an error if the key is not found
// or if there are issues accessing the cache.
func (cs *CacheService) Get(ctx context.Context, key string) ([]byte, error) {
	value, err := cs.repo.Get(ctx, key)
	if err != nil {
		if errors.Is(err, domain.ErrCacheNotFound) {
			return nil, domain.ErrCacheNotFound
		}
		return nil, domain.ErrInternal
	}
	return value, nil
}

// Delete removes the value associated with the specified key from the cache.
// Returns an error if the operation fails (e.g., if there are issues accessing the cache).
func (cs *CacheService) Delete(ctx context.Context, key string) error {
	err := cs.repo.Delete(ctx, key)
	if err != nil {
		return domain.ErrInternal
	}
	return nil
}

// DeleteByPrefix removes all values from the cache that match the given prefix.
// Returns an error if the operation fails (e.g., if there are issues accessing the cache).
func (cs *CacheService) DeleteByPrefix(ctx context.Context, prefix string) error {
	err := cs.repo.DeleteByPrefix(ctx, prefix)
	if err != nil {
		return domain.ErrInternal
	}
	return nil
}

// Close closes the connection to the cache server, ensuring that all resources are freed.
// Returns an error if the operation fails (e.g., if there are issues closing the connection).
func (cs *CacheService) Close() error {
	err := cs.repo.Close()
	if err != nil {
		return domain.ErrInternal
	}
	return nil
}
