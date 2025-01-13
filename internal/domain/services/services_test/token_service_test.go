//go:build !integration

package services_test

import (
	"context"
	"errors"
	"go-starter/internal/adapters/timegen"
	"go-starter/internal/domain"
	"go-starter/internal/domain/entities"
	"testing"
	"time"
)

func TestTokenService_ValidateToken(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := context.Background()
	timeGenerator := timegen.NewFakeTimeGenerator(time.Now().Add(-time.Hour))
	builder := NewTestBuilder().WithTimeGenerator(timeGenerator).Build()

	userToCreate := &entities.User{
		Username: "example",
		Password: "secret123",
	}

	_, err := builder.UserService.Register(ctx, userToCreate)
	if err != nil {
		t.Errorf("Failed to register user: %v", err)
	}

	accessToken, _, err := builder.AuthService.Login(ctx, userToCreate.Username, userToCreate.Password)
	if err != nil {
		t.Errorf("Failed to login: %v", err)
	}

	_, err = builder.TokenService.ValidateToken(accessToken)
	if !errors.Is(err, domain.ErrInvalidToken) {
		t.Errorf("ValidateToken() error = %v, wantErr %v", err, domain.ErrInvalidToken)
	}
}
