package config

import (
	"fmt"
	"go-starter/pkg/env"
	"time"
)

const (
	EnvProduction  = "production"
	EnvStaging     = "staging"
	EnvDevelopment = "development"
)

type (
	// Container contains environment variables for the application, database, http server, ...
	Container struct {
		Application *App
		DB          *DB
		HTTP        *HTTP
		Redis       *Redis
		Token       *Token
		ErrTracker  *ErrTracker
		Mailer      *Mailer
		FileUpload  *FileUpload
	}

	// App contains all the environment variables for the application.
	App struct {
		Env     string
		BaseURL string
	}

	// DB contains all the environment variables for the database.
	DB struct {
		Addr         string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  time.Duration
	}

	// HTTP contains all the environment variables for the http server.
	HTTP struct {
		Port int
	}

	// Redis contains all the environment variables for the cache service.
	Redis struct {
		Addr     string
		Password string
		DB       int
	}

	// Token contains all the environment variables for the token service.
	Token struct {
		AccessTokenDuration            time.Duration
		EmailVerificationTokenDuration time.Duration
		PasswordResetTokenDuration     time.Duration
	}

	// ErrTracker contains all the environment variables for the error tracking.
	ErrTracker struct {
		DSN              string
		TracesSampleRate float64
	}

	// Mailer contains all the environment variables for the mailer.
	Mailer struct {
		Region    string
		AccessKey string
		SecretKey string
		From      string
		DebugTo   string
	}

	// FileUpload contains all the environment variables for the file uploader.
	FileUpload struct {
		Region    string
		AccessKey string
		SecretKey string
		Bucket    string
	}
)

// New creates a new container instance.
func New() *Container {
	app := &App{
		Env:     env.GetString("ENVIRONMENT"),
		BaseURL: env.GetString("BASE_URL"),
	}

	db := &DB{
		Addr:         env.GetString("DB_ADDR"),
		MaxOpenConns: env.GetOptionalInt("DB_MAX_OPEN_CONNS", 30),
		MaxIdleConns: env.GetOptionalInt("DB_MAX_IDLE_CONNS", 30),
		MaxIdleTime:  env.GetOptionalDuration("DB_MAX_IDLE_TIME", 15*time.Minute),
	}

	http := &HTTP{
		Port: env.GetOptionalInt("HTTP_PORT", 8080),
	}

	redis := &Redis{
		Addr:     env.GetString("REDIS_ADDR"),
		Password: env.GetOptionalString("REDIS_PASSWORD", ""),
		DB:       env.GetOptionalInt("REDIS_DB", 0),
	}

	token := &Token{
		AccessTokenDuration:            env.GetOptionalDuration("ACCESS_TOKEN_DURATION", 15*time.Minute),
		EmailVerificationTokenDuration: env.GetOptionalDuration("EMAIL_VERIFICATION_TOKEN_DURATION", 24*time.Hour),
		PasswordResetTokenDuration:     env.GetOptionalDuration("PASSWORD_RESET_TOKEN_DURATION", 15*time.Minute),
	}

	errTracker := &ErrTracker{
		DSN:              env.GetOptionalString("SENTRY_DSN", ""),
		TracesSampleRate: env.GetOptionalFloat64("SENTRY_TRACES_SAMPLE_RATE", 1.0),
	}

	mailer := &Mailer{
		Region:    env.GetString("SES_REGION"),
		AccessKey: env.GetString("SES_ACCESS_KEY"),
		SecretKey: env.GetString("SES_SECRET_KEY"),
		From:      env.GetString("SES_FROM"),
		DebugTo:   env.GetString("SES_DEBUG_TO"),
	}

	fileUpload := &FileUpload{
		Region:    env.GetString("S3_REGION"),
		AccessKey: env.GetString("S3_ACCESS_KEY"),
		SecretKey: env.GetString("S3_SECRET_KEY"),
		Bucket:    env.GetString("S3_BUCKET"),
	}

	c := &Container{
		Application: app,
		DB:          db,
		HTTP:        http,
		Redis:       redis,
		Token:       token,
		ErrTracker:  errTracker,
		Mailer:      mailer,
		FileUpload:  fileUpload,
	}

	err := c.validate()
	if err != nil {
		panic(err)
	}

	return c
}

// validate validates the container.
func (c *Container) validate() error {
	// Application
	if c.Application.Env != EnvDevelopment && c.Application.Env != EnvProduction && c.Application.Env != EnvStaging {
		return fmt.Errorf("invalid environment variable: %s", "ENVIRONMENT")
	}

	// DB
	if c.DB.MaxIdleTime < 0 {
		return fmt.Errorf("invalid environment variable: %s", "DB_MAX_IDLE_TIME")
	}

	if c.DB.MaxIdleConns < 0 {
		return fmt.Errorf("invalid environment variable: %s", "DB_MAX_IDLE_CONNS")
	}

	if c.DB.MaxOpenConns < 0 {
		return fmt.Errorf("invalid environment variable: %s", "DB_MAX_OPEN_CONNS")
	}

	if c.DB.MaxOpenConns < c.DB.MaxIdleConns {
		return fmt.Errorf("invalid environment variables: %s should be greater or equal to %s", "DB_MAX_OPEN_CONNS", "DB_MAX_IDLE_CONNS")
	}

	// HTTP
	if c.HTTP.Port <= 0 {
		return fmt.Errorf("invalid environment variable: %s", "HTTP_PORT")
	}

	// Token
	if c.Token.AccessTokenDuration < 0 {
		return fmt.Errorf("invalid environment variable: %s", "ACCESS_TOKEN_DURATION")
	}

	// ErrTracker
	if c.ErrTracker.TracesSampleRate < 0 || c.ErrTracker.TracesSampleRate > 1.0 {
		return fmt.Errorf("invalid environment variable: %s should be between 0 and 1", "SENTRY_TRACES_SAMPLE_RATE")
	}

	return nil
}
