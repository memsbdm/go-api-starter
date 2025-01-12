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
	data          map[string][]byte
	timer         map[string]time.Time
	mu            sync.Mutex
	timeGenerator ports.TimeGenerator
}

// NewMock creates a new mock instance of Redis
func NewMock(timeGenerator ports.TimeGenerator) ports.CacheRepository {
	return &RedisMock{
		data:          make(map[string][]byte),
		timer:         make(map[string]time.Time),
		mu:            sync.Mutex{},
		timeGenerator: timeGenerator,
	}
}

// Set stores the value in the redis database
func (rm *RedisMock) Set(_ context.Context, key string, value []byte, ttl time.Duration) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.data[key] = value
	rm.timer[key] = rm.timeGenerator.Now().Add(ttl)
	return nil
}

// Get retrieves the value from the redis database
func (rm *RedisMock) Get(_ context.Context, key string) ([]byte, error) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if value, ok := rm.timer[key]; ok {
		if value.After(rm.timeGenerator.Now()) {
			return rm.data[key], nil
		}
	}
	return nil, domain.ErrCacheNotFound
}

// Delete removes the value from the redis database
func (rm *RedisMock) Delete(_ context.Context, key string) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	delete(rm.data, key)
	delete(rm.timer, key)
	return nil
}

// DeleteByPrefix removes the value from the redis database with the given prefix
func (rm *RedisMock) DeleteByPrefix(_ context.Context, prefix string) error {
	for key := range rm.data {
		if strings.HasPrefix(key, prefix) {
			delete(rm.data, key)
		}
	}
	for key := range rm.timer {
		if strings.HasPrefix(key, prefix) {
			delete(rm.timer, key)
		}
	}
	return nil
}

// Close closes the connection to the redis database
func (rm *RedisMock) Close() error {
	return nil
}
