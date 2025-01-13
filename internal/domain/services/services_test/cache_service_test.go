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

	tests := []struct {
		name          string
		input         string
		expectedValue []byte
		expectedErr   error
	}{
		{
			name:          "get an existing key",
			input:         key,
			expectedValue: value,
			expectedErr:   nil,
		},
		{
			name:          "get an non existing key",
			input:         "non-existing",
			expectedValue: nil,
			expectedErr:   domain.ErrCacheNotFound,
		},
	}

	// Act & Assert
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			value, err := builder.CacheService.Get(ctx, tt.input)
			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
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

	tests := []struct {
		name        string
		input       string
		expectedErr error
	}{
		{
			name:        "delete an existing key",
			input:       keyValueToStore.key,
			expectedErr: nil,
		},
		{
			name:        "delete an non existing key",
			input:       "non-existing",
			expectedErr: nil,
		},
	}

	// Act & Assert
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := builder.CacheService.Delete(ctx, tt.input)
			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
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
	const otherKey = "other"
	otherValue := []byte("value")

	err := builder.CacheService.Set(ctx, otherKey, otherValue, time.Hour)
	if err != nil {
		t.Fatalf("failed to set cache: %v", err)
	}

	err = builder.CacheService.DeleteByPrefix(ctx, testedPrefix)
	if err != nil {
		t.Fatalf("failed to delete cache by prefix: %v", err)
	}

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

	// Arrange
	ctx := context.Background()
	timeGenerator := timegen.NewFakeTimeGenerator(time.Now())
	builder := NewTestBuilder().WithTimeGenerator(timeGenerator).Build()
	const key = "example"
	err := builder.CacheService.Set(ctx, key, []byte("value"), time.Hour)
	if err != nil {
		t.Fatalf("failed to set cache: %v", err)
	}
	_, err = builder.CacheService.Get(ctx, key)
	if err != nil {
		t.Fatalf("failed to get cache: %v", err)
	}

	timeGenerator.Advance(time.Hour)

	// Act & Assert
	_, err = builder.CacheService.Get(ctx, key)
	if err == nil {
		t.Errorf("expected error %v, nil", domain.ErrCacheNotFound)
	}
}
