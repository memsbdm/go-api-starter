package services

import (
	"context"
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

// Set stores the value in the cache
func (cs *CacheService) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	err := cs.repo.Set(ctx, key, value, ttl)
	if err != nil {
		return domain.ErrInternal
	}
	return nil
}

// Get retrieves the value from the cache
func (cs *CacheService) Get(ctx context.Context, key string) ([]byte, error) {
	value, err := cs.repo.Get(ctx, key)
	if err != nil {
		return nil, domain.ErrCacheNotFound
	}
	return value, nil
}

// Delete removes the value from the cache
func (cs *CacheService) Delete(ctx context.Context, key string) error {
	err := cs.repo.Delete(ctx, key)
	if err != nil {
		return domain.ErrInternal
	}
	return nil
}

// DeleteByPrefix removes the value from the cache with the given prefix
func (cs *CacheService) DeleteByPrefix(ctx context.Context, prefix string) error {
	err := cs.repo.DeleteByPrefix(ctx, prefix)
	if err != nil {
		return domain.ErrInternal
	}
	return nil
}

// Close closes the connection to the cache server
func (cs *CacheService) Close() error {
	return cs.repo.Close()
}
