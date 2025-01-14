package config

import (
	"go-starter/pkg/env"
	"time"
)

type (
	// Container contains environment variables for the application, database and http server
	Container struct {
		Application *App
		DB          *DB
		HTTP        *HTTP
		Redis       *Redis
		Token       *Token
	}

	// App contains all the environment variables for the application
	App struct {
		Env string
	}

	// DB contains all the environment variables for the database
	DB struct {
		Addr         string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
	}

	// HTTP contains all the environment variables for the http server
	HTTP struct {
		Port int
	}

	// Redis contains all the environment variables for the cache service
	Redis struct {
		Addr     string
		Password string
	}

	// Token contains all the environment variables for the token service
	Token struct {
		JWTSecret            []byte
		AccessTokenDuration  time.Duration
		RefreshTokenDuration time.Duration
	}
)

// New creates a new container instance
func New() *Container {
	app := &App{
		Env: env.GetString("ENVIRONMENT"),
	}

	db := &DB{
		Addr:         env.GetString("DB_ADDR"),
		MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS"),
		MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS"),
		MaxIdleTime:  env.GetString("DB_MAX_IDLE_TIME"),
	}

	http := &HTTP{
		Port: env.GetInt("HTTP_PORT"),
	}

	redis := &Redis{
		Addr:     env.GetString("REDIS_ADDR"),
		Password: env.GetString("REDIS_PASSWORD"),
	}

	token := &Token{
		JWTSecret:            []byte(env.GetString("JWT_SECRET")),
		AccessTokenDuration:  env.GetDuration("ACCESS_TOKEN_DURATION"),
		RefreshTokenDuration: env.GetDuration("REFRESH_TOKEN_DURATION"),
	}

	return &Container{
		Application: app,
		DB:          db,
		HTTP:        http,
		Redis:       redis,
		Token:       token,
	}
}
