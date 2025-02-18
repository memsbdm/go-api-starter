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
	ErrTracker    ports.ErrorTracker
	MailerService ports.MailerService
}

// New creates and initializes a new Services instance with the provided dependencies.
func New(cfg *config.Container, a *adapters.Adapters) *Services {
	cacheSvc := NewCacheService(a.CacheRepository, a.ErrTracker)
	userSvc := NewUserService(a.UserRepository, cacheSvc, a.ErrTracker)
	tokenSvc := NewTokenService(cfg.Token, a.TokenRepository, cacheSvc, a.ErrTracker)
	mailerSvc := NewMailerService(cfg, a.MailerRepository, a.ErrTracker)
	authSvc := NewAuthService(userSvc, tokenSvc, a.ErrTracker)
	return &Services{
		CacheService:  cacheSvc,
		UserService:   userSvc,
		AuthService:   authSvc,
		TokenService:  tokenSvc,
		ErrTracker:    a.ErrTracker,
		MailerService: mailerSvc,
	}
}
