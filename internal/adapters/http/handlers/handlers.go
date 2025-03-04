package handlers

import (
	"go-starter/internal/domain/ports"
	"go-starter/internal/domain/services"
)

// Handlers holds all handler implementations for the application.
type Handlers struct {
	HealthHandler *HealthHandler
	AuthHandler   *AuthHandler
	UserHandler   *UserHandler
	MailerHandler *MailerHandler
}

// New creates and initializes a new Handlers instance with the provided dependencies.
func New(s *services.Services, errTracker ports.ErrTrackerAdapter) *Handlers {
	return &Handlers{
		HealthHandler: NewHealthHandler(),
		AuthHandler:   NewAuthHandler(s.AuthService),
		UserHandler:   NewUserHandler(s.UserService, errTracker),
		MailerHandler: NewMailerHandler(s.MailerService),
	}
}
