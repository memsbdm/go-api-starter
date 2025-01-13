// //go:build !integration

package services_test

import (
	"go-starter/config"
	"go-starter/internal/adapters/storage/postgres/repositories/mocks"
	"go-starter/internal/adapters/storage/redis"
	"go-starter/internal/adapters/timegen"
	"go-starter/internal/adapters/token"
	"go-starter/internal/domain/services"
	"time"
)

var (
	authService  *services.AuthService
	cacheService *services.CacheService
	userService  *services.UserService
)

func init() {
	timeGenerator := timegen.NewRealTimeGenerator()
	cacheRepo := redis.NewMock(timeGenerator)
	cacheService = services.NewCacheService(cacheRepo)
	tokenRepo := token.NewTokenRepository()
	tokenService := services.NewTokenService(&config.Token{
		TokenDuration:        10 * time.Minute,
		RefreshTokenDuration: 1 * time.Hour,
		JWTSecret:            []byte("secret"),
	}, tokenRepo, cacheService)
	userRepo := mocks.MockUserRepository()
	userService = services.NewUserService(userRepo, cacheService)
	authService = services.NewAuthService(userService, tokenService)
}
