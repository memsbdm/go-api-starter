// //go:build !integration

package services_test

import (
	"go-starter/config"
	"go-starter/internal/adapters/auth"
	"go-starter/internal/adapters/storage/postgres/repositories/mocks"
	"go-starter/internal/adapters/storage/redis"
	"go-starter/internal/domain/services"
	"time"
)

var (
	authService  *services.AuthService
	cacheService *services.CacheService
	userService  *services.UserService
)

func init() {
	tokenService := auth.NewTokenService(&config.Token{
		Duration:  10 * time.Minute,
		JWTSecret: []byte("secret"),
	})
	cacheRepo := redis.NewMock()
	cacheService = services.NewCacheService(cacheRepo)
	userRepo := mocks.MockUserRepository()
	userService = services.NewUserService(userRepo, cacheService)
	authService = services.NewAuthService(userService, tokenService)
}
