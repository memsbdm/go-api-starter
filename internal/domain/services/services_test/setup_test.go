// //go:build !integration

package services_test

import (
	"go-starter/config"
	"go-starter/internal/adapters/storage/postgres/repositories/mocks"
	"go-starter/internal/adapters/storage/redis"
	"go-starter/internal/domain/services"
)

var (
	authService  *services.AuthService
	cacheService *services.CacheService
	userService  *services.UserService
)

func init() {
	cacheRepo := redis.NewMock()
	cacheService = services.NewCacheService(cacheRepo)
	userRepo := mocks.MockUserRepository()
	userService = services.NewUserService(userRepo, cacheService)
	authService = services.NewAuthService(&config.Security{
		JWTSecret: []byte("secret"),
	}, userService)
}
