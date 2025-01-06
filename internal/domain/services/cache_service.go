package services

import (
	"context"
	"encoding/json"
	"fmt"
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

// GenerateCacheKey generates a cache key based on the input parameters
func (cs *CacheService) GenerateCacheKey(prefix string, params any) string {
	return fmt.Sprintf("%s:%v", prefix, params)
}

// GenerateCacheKeyParams generates a cache params based on the input parameters
func (cs *CacheService) GenerateCacheKeyParams(params ...any) string {
	var str string

	for i, param := range params {
		str += fmt.Sprintf("%v", param)

		last := len(params) - 1
		if i != last {
			str += "-"
		}
	}

	return str
}

// Serialize marshals the input data into an array of bytes
func (cs *CacheService) Serialize(data any) ([]byte, error) {
	return json.Marshal(data)
}

// Deserialize unmarshals the input data into the output interface
func (cs *CacheService) Deserialize(data []byte, output any) error {
	return json.Unmarshal(data, output)
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
