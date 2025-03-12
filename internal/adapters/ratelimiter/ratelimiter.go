package ratelimiter

import (
	"context"
	"fmt"
	"go-starter/internal/domain/ports"
	"time"
)

// fixed window rate limiter lua script
const rateLimiterScript = `
local key = KEYS[1]
local limit = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local current = redis.call('INCR', key)

if current == 1 then
  redis.call('EXPIRE', key, window)
end

local ttl = redis.call('TTL', key)
return {current, limit, ttl}
`

// RateLimiter provides rate limiting functionality using the cache repository
type RateLimiter struct {
	cache     ports.CacheRepository
	keyPrefix string
	name      string
}

// New creates a new RateLimiter instance
func New(cache ports.CacheRepository, name string) *RateLimiter {
	keyPrefix := "rate_limit:"

	return &RateLimiter{
		cache:     cache,
		keyPrefix: keyPrefix,
		name:      name,
	}
}

// Result represents the outcome of a rate limit check
type Result struct {
	Allowed    bool
	Current    int64
	Limit      int64
	ResetAfter time.Duration
}

// Check verifies if the request should be allowed based on the rate limit
func (rl *RateLimiter) Check(ctx context.Context, key string, limit int64, window time.Duration) (*Result, error) {
	// Format: rate_limit:limiter_name:IP
	redisKey := fmt.Sprintf("%s%s:%s", rl.keyPrefix, rl.name, key)

	// Execute the Lua script
	res, err := rl.cache.Eval(ctx, rateLimiterScript, []string{redisKey}, limit, int(window.Seconds()))
	if err != nil {
		return nil, fmt.Errorf("failed to execute rate limit script: %w", err)
	}

	// Parse the result
	values, ok := res.([]interface{})
	if !ok || len(values) != 3 {
		return nil, fmt.Errorf("unexpected result format from rate limit script")
	}

	// Extract current count and limit
	current, ok := values[0].(int64)
	if !ok {
		// Try to convert from float64 (Redis sometimes returns numbers as float64)
		if floatVal, isFloat := values[0].(float64); isFloat {
			current = int64(floatVal)
		} else {
			return nil, fmt.Errorf("unexpected type for current count")
		}
	}

	limitVal, ok := values[1].(int64)
	if !ok {
		// Try to convert from float64
		if floatVal, isFloat := values[1].(float64); isFloat {
			limitVal = int64(floatVal)
		} else {
			return nil, fmt.Errorf("unexpected type for limit value")
		}
	}

	// Extract TTL
	ttlVal, ok := values[2].(int64)
	if !ok {
		if floatVal, isFloat := values[2].(float64); isFloat {
			ttlVal = int64(floatVal)
		} else {
			ttlVal = int64(window.Seconds()) // Fallback
		}
	}
	resetAfter := time.Duration(ttlVal) * time.Second

	return &Result{
		Allowed:    current <= limitVal,
		Current:    current,
		Limit:      limitVal,
		ResetAfter: resetAfter,
	}, nil
}
