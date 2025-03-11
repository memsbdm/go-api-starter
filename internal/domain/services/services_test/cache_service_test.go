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

	const key = "key"
	const value = "value"
	const cacheDuration = time.Hour

	tests := map[string]struct {
		initCache   func(*TestBuilder)
		advance     time.Duration
		expectedErr error
	}{
		"get an existing key": {
			initCache: func(builder *TestBuilder) {
				err := builder.CacheService.Set(context.Background(), key, []byte(value), cacheDuration)
				if err != nil {
					t.Fatalf("failed to set cache: %v", err)
				}
			},
			advance:     0,
			expectedErr: nil,
		},
		"get a non existing key": {
			advance:     0,
			expectedErr: domain.ErrCacheNotFound,
		},
		"get a key that is expired": {
			initCache: func(builder *TestBuilder) {
				err := builder.CacheService.Set(context.Background(), key, []byte(value), cacheDuration)
				if err != nil {
					t.Fatalf("failed to set cache: %v", err)
				}
			},
			advance:     cacheDuration,
			expectedErr: domain.ErrCacheNotFound,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			timeGenerator := timegen.NewTimeGeneratorMock(time.Now())
			builder := NewTestBuilder().WithTimeGenerator(timeGenerator).Build()
			if tt.initCache != nil {
				tt.initCache(builder)
			}

			advanceTime(t, builder.TimeGenerator, tt.advance)

			// Act
			val, err := builder.CacheService.Get(context.Background(), key)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
			if tt.expectedErr == nil && !reflect.DeepEqual(val, []byte(value)) {
				t.Errorf("expected value %v, got %v", val, []byte(value))
			}
		})
	}
}

func TestCacheService_Delete(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()
	builder := NewTestBuilder().Build()

	const key = "key"
	const value = "value"

	err := builder.CacheService.Set(ctx, key, []byte(value), time.Hour)
	if err != nil {
		t.Fatalf("failed to set cache: %v", err)
	}

	// Act
	err = builder.CacheService.Delete(ctx, key)
	if err != nil {
		t.Fatalf("failed to delete cache: %v", err)
	}

	// Assert
	_, err = builder.CacheService.Get(ctx, key)
	if err == nil || !errors.Is(err, domain.ErrCacheNotFound) {
		t.Errorf("expected error %v, got %v", domain.ErrCacheNotFound, err)
	}
}

func TestCacheService_DeleteByPrefix(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()
	builder := NewTestBuilder().Build()

	const prefix = "prefix"
	const value = "value"

	for i := range 4 {
		err := builder.CacheService.Set(ctx, fmt.Sprintf("%s%d", prefix, i), []byte(value), time.Hour)
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

	// Act
	err = builder.CacheService.DeleteByPrefix(ctx, prefix)
	if err != nil {
		t.Fatalf("failed to delete cache by prefix: %v", err)
	}

	// Assert
	for i := range 4 {
		_, err := builder.CacheService.Get(ctx, fmt.Sprintf("%s%d", prefix, i))
		if err == nil || !errors.Is(err, domain.ErrCacheNotFound) {
			t.Errorf("expected error %v, got %v", domain.ErrCacheNotFound, err)
		}
	}

	cachedOtherValue, err := builder.CacheService.Get(ctx, otherKey)
	if err != nil {
		t.Fatalf("failed to get cached value: %v", err)
	}
	if !reflect.DeepEqual(cachedOtherValue, otherValue) {
		t.Errorf("expected value to be %v, got %v", otherValue, cachedOtherValue)
	}
}
