package logger

import (
	"go-starter/config"
	"log/slog"
	"os"
)

// New defines logger specifications based on application environment
func New(config *config.App) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if config.Env == "production" {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	slog.SetDefault(logger)
}
