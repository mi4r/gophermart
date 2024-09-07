package logger

import (
	"log/slog"
	"os"
)

func InitLogger(level string) {
	var logger *slog.Logger
	switch level {
	case "info", "production", "prod":
		logger = slog.New(slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelInfo,
			},
		))
	default:
		logger = slog.New(slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelInfo,
			},
		))
	}
	slog.SetDefault(logger)
}
