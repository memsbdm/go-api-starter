package redis

import (
	"context"
	"go-starter/internal/domain"
	"go-starter/internal/domain/ports"
	"strings"
	"sync"
	"time"
)

// RedisMock implements ports.CacheRepository interface and provides access to in-memory data
type RedisMock struct {
	data map[string][]byte
	mu   sync.Mutex
}

// NewMock creates a new mock instance of Redis
func NewMock() ports.CacheRepository {
	return &RedisMock{
		data: make(map[string][]byte),
		mu:   sync.Mutex{},
	}
}

// Set stores the value in the redis database
func (rm *RedisMock) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.data[key] = value
	return nil
}

// Get retrieves the value from the redis database
func (rm *RedisMock) Get(ctx context.Context, key string) ([]byte, error) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if value, ok := rm.data[key]; ok {
		return value, nil
	}
	return nil, domain.ErrCacheNotFound
}

// Delete removes the value from the redis database
func (rm *RedisMock) Delete(ctx context.Context, key string) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	delete(rm.data, key)
	return nil
}

// DeleteByPrefix removes the value from the redis database with the given prefix
func (rm *RedisMock) DeleteByPrefix(ctx context.Context, prefix string) error {
	for key := range rm.data {
		if strings.HasPrefix(key, prefix) {
			delete(rm.data, key)
		}
	}
	return nil
}

// Close closes the connection to the redis database
func (rm *RedisMock) Close() error {
	return nil
}
