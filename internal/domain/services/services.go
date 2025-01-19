package services

import (
	"go-starter/config"
	"go-starter/internal/adapters"
	"go-starter/internal/domain/ports"
)

type Services struct {
	CacheService ports.CacheService
	UserService  ports.UserService
	AuthService  ports.AuthService
	TokenService ports.TokenService
}

func New(cfg *config.Container, adapters *adapters.Adapters) *Services {
	cacheSvc := NewCacheService(adapters.CacheRepository)
	userSvc := NewUserService(adapters.UserRepository, cacheSvc)
	tokenSvc := NewTokenService(cfg.Token, adapters.TokenRepository, cacheSvc)
	authSvc := NewAuthService(userSvc, tokenSvc)
	return &Services{
		CacheService: cacheSvc,
		UserService:  userSvc,
		AuthService:  authSvc,
		TokenService: tokenSvc,
	}
}
