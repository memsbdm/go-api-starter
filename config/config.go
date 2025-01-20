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
		ErrTracker  *ErrTracker
	}

	// App contains all the environment variables for the application.
	App struct {
		Env string
	}

	// DB contains all the environment variables for the database.
	DB struct {
		Addr         string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
	}

	// HTTP contains all the environment variables for the http server.
	HTTP struct {
		Port int
	}

	// Redis contains all the environment variables for the cache service.
	Redis struct {
		Addr     string
		Password string
	}

	// Token contains all the environment variables for the token service.
	Token struct {
		AccessTokenSignature  []byte
		RefreshTokenSignature []byte
		AccessTokenDuration   time.Duration
		RefreshTokenDuration  time.Duration
	}

	// ErrTracker contains all the environment variables for the error tracking.
	ErrTracker struct {
		DSN              string
		TracesSampleRate float64
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
		AccessTokenSignature:  []byte(env.GetString("ACCESS_TOKEN_SIGNATURE")),
		RefreshTokenSignature: []byte(env.GetString("REFRESH_TOKEN_SIGNATURE")),
		AccessTokenDuration:   env.GetDuration("ACCESS_TOKEN_DURATION"),
		RefreshTokenDuration:  env.GetDuration("REFRESH_TOKEN_DURATION"),
	}

	errTracker := &ErrTracker{
		DSN:              env.GetString("SENTRY_DSN"),
		TracesSampleRate: env.GetFloat64("SENTRY_TRACES_SAMPLE_RATE"),
	}

	return &Container{
		Application: app,
		DB:          db,
		HTTP:        http,
		Redis:       redis,
		Token:       token,
		ErrTracker:  errTracker,
	}
}
