package mocks

import (
	"context"
	"go-starter/internal/domain"
	"go-starter/internal/domain/ports"
	"strings"
	"sync"
	"time"
)

// CacheMock implements the ports.CacheRepository interface and provides an in-memory cache simulation.
// It allows for testing caching functionalities without the need for a real caching system.
type CacheMock struct {
	data          map[string][]byte
	timer         map[string]time.Time
	mu            sync.Mutex
	timeGenerator ports.TimeGenerator
}

// NewCacheMock creates a new mock instance of the cache.
func NewCacheMock(timeGenerator ports.TimeGenerator) ports.CacheRepository {
	return &CacheMock{
		data:          make(map[string][]byte),
		timer:         make(map[string]time.Time),
		mu:            sync.Mutex{},
		timeGenerator: timeGenerator,
	}
}

// Set stores the value in the cache with a specified key and time-to-live (TTL).
// Returns an error if the operation fails (e.g., if the cache is unreachable).
func (cm *CacheMock) Set(_ context.Context, key string, value []byte, ttl time.Duration) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.data[key] = value
	cm.timer[key] = cm.timeGenerator.Now().Add(ttl)
	return nil
}

// Get retrieves the value associated with the specified key from the cache.
// Returns the value as a byte slice and an error if the key is not found
// or if there are issues accessing the cache.
func (cm *CacheMock) Get(_ context.Context, key string) ([]byte, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if value, ok := cm.timer[key]; ok {
		if value.After(cm.timeGenerator.Now()) {
			return cm.data[key], nil
		}
	}
	return nil, domain.ErrCacheNotFound
}

// Delete removes the value associated with the specified key from the cache.
// Returns an error if the operation fails (e.g., if there are issues accessing the cache).
func (cm *CacheMock) Delete(_ context.Context, key string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.data, key)
	delete(cm.timer, key)
	return nil
}

// DeleteByPrefix removes all values from the cache that match the given prefix.
// Returns an error if the operation fails (e.g., if there are issues accessing the cache).
func (cm *CacheMock) DeleteByPrefix(_ context.Context, prefix string) error {
	for key := range cm.data {
		if strings.HasPrefix(key, prefix) {
			delete(cm.data, key)
		}
	}
	for key := range cm.timer {
		if strings.HasPrefix(key, prefix) {
			delete(cm.timer, key)
		}
	}
	return nil
}

// Close closes the connection to the cache server, ensuring that all resources are freed.
// Returns an error if the operation fails (e.g., if there are issues closing the connection).
func (cm *CacheMock) Close() error {
	return nil
}
