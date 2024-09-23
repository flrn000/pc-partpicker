package logging

import (
	"log/slog"
)

func NewLogger(handler slog.Handler) *slog.Logger {
	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}
