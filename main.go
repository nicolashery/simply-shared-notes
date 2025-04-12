package main

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"

	"github.com/nicolashery/simply-shared-notes/db"
	"github.com/nicolashery/simply-shared-notes/server"
)

//go:embed sql/pragmas.sql
var pragmasSQL string

//go:embed all:dist
var distFS embed.FS

func run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	h := slog.NewTextHandler(os.Stdout, nil)
	logger := slog.New(h)

	port := 3000
	if portStr := os.Getenv("PORT"); portStr != "" {
		var err error
		port, err = strconv.Atoi(portStr)
		if err != nil {
			return fmt.Errorf("invalid PORT environment variable: %w", err)
		}
	}

	isDev := false
	if isDevStr := os.Getenv("DEV"); isDevStr == "true" {
		isDev = true
	}

	var assetsConfig server.AssetsConfig
	var err error
	if isDev {
		assetsConfig, err = server.DevAssets()
		if err != nil {
			return err
		}
		logger.Info("using dev assets, make sure Vite is running")
	} else {
		assetsConfig, err = server.ProdAssets(distFS)
		if err != nil {
			return err
		}
		logger.Info("using prod assets")
	}

	dbPath := "data/app.sqlite"
	dbConn, err := db.InitDB(ctx, dbPath, pragmasSQL)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	queries := db.New(dbConn)

	s := server.NewServer(logger, queries, assetsConfig)

	return server.RunServer(ctx, s, logger, port)
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
