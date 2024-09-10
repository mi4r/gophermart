package logger

import (
	"log/slog"
	"os"
)

func InitLogger(level string) {
	// var logger *slog.Logger
	var opts PrettyHandlerOptions

	switch level {
	case "info", "production", "prod":
		opts = PrettyHandlerOptions{
			SlogOpts: slog.HandlerOptions{
				Level: slog.LevelInfo,
			},
		}
	default:
		opts = PrettyHandlerOptions{
			SlogOpts: slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		}
	}
	handler := newPrettyHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
