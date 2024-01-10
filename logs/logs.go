package logs

import (
	"log/slog"
	"os"
)

// NewFromEnvironment creates a logger from the environment variables
func NewFromEnvironment() *slog.Logger {
	// set the level
	level := slog.LevelInfo
	if v, ok := os.LookupEnv("LOG_LEVEL"); ok {
		switch v {
		case "debug":
			level = slog.LevelDebug
		case "warn":
			level = slog.LevelWarn
		case "error":
			level = slog.LevelError
		default:
			level = slog.LevelInfo
		}
	}

	// set the output
	output := os.Stdout
	if v, ok := os.LookupEnv("LOG_OUTPUT"); ok {
		switch v {
		case "stderr":
			output = os.Stderr
		default:
			output = os.Stdout
		}
	}

	// set the handler
	var handler slog.Handler
	if v, ok := os.LookupEnv("LOG_FORMAT"); ok {
		switch v {
		case "json":
			handler = slog.NewJSONHandler(output, &slog.HandlerOptions{
				Level: level,
			})
		default:
			handler = slog.NewTextHandler(output, &slog.HandlerOptions{
				Level: level,
			})
		}
	}

	// if the handler is nil, then we will use the default text handler
	if handler == nil {
		handler = slog.NewTextHandler(output, &slog.HandlerOptions{
			Level: level,
		})
	}

	return slog.New(handler)
}
