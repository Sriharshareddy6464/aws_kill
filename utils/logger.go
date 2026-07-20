package utils

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func init() {
	// Initialize standard structured logger outputting to stderr
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	Logger = slog.New(handler)
}
