package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/nicolashery/simply-shared-notes/server"
)

func run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	h := slog.NewTextHandler(os.Stdout, nil)
	logger := slog.New(h)

	port := 3000

	s := server.NewServer(logger)

	return server.RunServer(ctx, s, logger, port)
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
