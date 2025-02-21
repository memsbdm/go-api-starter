package services

import (
	"go-starter/config"
	"go-starter/internal/adapters"
	"go-starter/internal/domain/ports"
)

// Services holds all service implementations for the application.
type Services struct {
	CacheService  ports.CacheService
	UserService   ports.UserService
	AuthService   ports.AuthService
	TokenService  ports.TokenService
	ErrTracker    ports.ErrTrackerAdapter
	MailerService ports.MailerService
}

// New creates and initializes a new Services instance with the provided dependencies.
func New(cfg *config.Container, a *adapters.Adapters) *Services {
	cacheSvc := NewCacheService(a.CacheRepository, a.ErrTrackerAdapter)
	userSvc := NewUserService(a.UserRepository, cacheSvc)
	tokenSvc := NewTokenService(cfg.Token, a.TokenRepository, cacheSvc)
	mailerSvc := NewMailerService(cfg, a.MailerAdapter, a.ErrTrackerAdapter)
	authSvc := NewAuthService(userSvc, tokenSvc, a.ErrTrackerAdapter)
	return &Services{
		CacheService:  cacheSvc,
		UserService:   userSvc,
		AuthService:   authSvc,
		TokenService:  tokenSvc,
		ErrTracker:    a.ErrTrackerAdapter,
		MailerService: mailerSvc,
	}
}
