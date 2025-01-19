package handlers

import (
	"go-starter/internal/domain/services"
)

// Handlers holds all handler implementations for the application.
type Handlers struct {
	HealthHandler *HealthHandler
	AuthHandler   *AuthHandler
	UserHandler   *UserHandler
}

// New creates and initializes a new Handlers instance with the provided dependencies.
func New(s *services.Services) *Handlers {
	return &Handlers{
		HealthHandler: NewHealthHandler(),
		AuthHandler:   NewAuthHandler(s.AuthService),
		UserHandler:   NewUserHandler(s.UserService),
	}
}
