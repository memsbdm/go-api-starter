//go:build !integration

package services_test

import (
	"context"
	"errors"
	"fmt"
	"go-starter/internal/adapters/timegen"
	"go-starter/internal/domain"
	"reflect"
	"testing"
	"time"
)

func TestCacheService_Get(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()
	builder := NewTestBuilder().Build()

	key := "example"
	value := []byte("value")

	err := builder.CacheService.Set(ctx, key, value, time.Hour)
	if err != nil {
		t.Fatalf("failed to set cache: %v", err)
	}

	tests := map[string]struct {
		input         string
		expectedValue []byte
		expectedErr   error
	}{
		"get an existing key": {
			input:         key,
			expectedValue: value,
			expectedErr:   nil,
		},
		"get an non existing key": {
			input:         "non-existing",
			expectedValue: nil,
			expectedErr:   domain.ErrCacheNotFound,
		},
	}

	// Act & Assert
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			value, err := builder.CacheService.Get(ctx, tt.input)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if !reflect.DeepEqual(value, tt.expectedValue) {
				t.Errorf("expected value %v, got %v", tt.expectedValue, value)
			}
		})
	}
}

func TestCacheService_Delete(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()
	builder := NewTestBuilder().Build()

	type setRequest struct {
		key   string
		value []byte
	}
	keyValueToStore := setRequest{
		key:   "example",
		value: []byte("value"),
	}
	err := builder.CacheService.Set(ctx, keyValueToStore.key, keyValueToStore.value, time.Hour)
	if err != nil {
		t.Fatalf("failed to set cache: %v", err)
	}

	tests := map[string]struct {
		input       string
		expectedErr error
	}{
		"delete an existing key": {
			input:       keyValueToStore.key,
			expectedErr: nil,
		},
		"delete an non existing key": {
			input:       "non-existing",
			expectedErr: nil,
		},
	}

	// Act & Assert
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := builder.CacheService.Delete(ctx, tt.input)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			value, err := builder.CacheService.Get(ctx, tt.input)
			if !errors.Is(err, domain.ErrCacheNotFound) {
				t.Errorf("expected error %v, got %v", domain.ErrCacheNotFound, err)
			}
			if value != nil {
				t.Errorf("expected value to be nil, got %v", value)
			}
		})
	}
}

func TestCacheService_DeleteByPrefix(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()
	builder := NewTestBuilder().Build()

	const testedPrefix = "example"

	for i := range 10 {
		err := builder.CacheService.Set(ctx, fmt.Sprintf("%s%d", testedPrefix, i), []byte("value"), time.Hour)
		if err != nil {
			t.Fatalf("failed to set cache for key %d: %v", i, err)
		}
	}
	otherKey := "other"
	otherValue := []byte("value")

	err := builder.CacheService.Set(ctx, otherKey, otherValue, time.Hour)
	if err != nil {
		t.Fatalf("failed to set cache: %v", err)
	}

	err = builder.CacheService.DeleteByPrefix(ctx, testedPrefix)
	if err != nil {
		t.Fatalf("failed to delete cache by prefix: %v", err)
	}

	// Act & Assert
	for i := range 10 {
		value, err := builder.CacheService.Get(ctx, fmt.Sprintf("%s%d", testedPrefix, i))
		if err == nil || !errors.Is(err, domain.ErrCacheNotFound) {
			t.Errorf("expected error %v, got %v", domain.ErrCacheNotFound, err)
		}
		if value != nil {
			t.Errorf("expected value to be nil, got %v", value)
		}
	}

	value, _ := builder.CacheService.Get(ctx, otherKey)
	if !reflect.DeepEqual(value, otherValue) {
		t.Errorf("expected value to be %v, got %v", otherValue, value)
	}
}

func TestCacheService_CacheExpiration(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		duration    time.Duration
		advance     time.Duration
		expectedErr error
	}{
		"cache not expired": {
			duration:    time.Hour,
			advance:     time.Minute,
			expectedErr: nil,
		},
		"cache expired": {
			duration:    time.Hour,
			advance:     time.Hour + time.Second,
			expectedErr: domain.ErrCacheNotFound,
		},
		"cache near expiration": {
			duration:    time.Hour,
			advance:     time.Hour - time.Second,
			expectedErr: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			ctx := context.Background()
			key := "example"
			timeGenerator := timegen.NewFakeTimeGenerator(time.Now())
			builder := NewTestBuilder().WithTimeGenerator(timeGenerator).Build()

			err := builder.CacheService.Set(ctx, key, []byte("value"), tt.duration)
			if err != nil {
				t.Fatalf("failed to set cache: %v", err)
			}

			_, err = builder.CacheService.Get(ctx, key)
			if err != nil {
				t.Fatalf("failed to get cache: %v", err)
			}

			timeGenerator.Advance(tt.advance)

			// Act & Assert
			_, err = builder.CacheService.Get(ctx, key)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}
