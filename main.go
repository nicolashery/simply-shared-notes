package main

import (
	"log/slog"
	"os"
)

func NewLogger() *slog.Logger {
	h := slog.NewTextHandler(os.Stdout, nil)
	logger := slog.New(h)

	return logger
}

func main() {
	logger := NewLogger()

	logger.Info("Hello, world!", slog.String("name", "world"))
}
