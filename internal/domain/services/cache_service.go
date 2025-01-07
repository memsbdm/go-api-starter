package services

import (
	"context"
	"go-starter/internal/domain/ports"
	"time"
)

// CacheService implements ports.CacheService interface and provides access to the cache repository
type CacheService struct {
	repo ports.CacheRepository
}

// NewCacheService creates a new user services instance
func NewCacheService(repo ports.CacheRepository) *CacheService {
	return &CacheService{
		repo: repo,
	}
}

// Set stores the value in the cache
func (cs *CacheService) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return cs.repo.Set(ctx, key, value, ttl)
}

// Get retrieves the value from the cache
func (cs *CacheService) Get(ctx context.Context, key string) ([]byte, error) {
	return cs.repo.Get(ctx, key)
}

// Delete removes the value from the cache
func (cs *CacheService) Delete(ctx context.Context, key string) error {
	return cs.repo.Delete(ctx, key)
}

// DeleteByPrefix removes the value from the cache with the given prefix
func (cs *CacheService) DeleteByPrefix(ctx context.Context, prefix string) error {
	return cs.repo.DeleteByPrefix(ctx, prefix)
}

// Close closes the connection to the cache server
func (cs *CacheService) Close() error {
	return cs.repo.Close()
}
