package main

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/nicolashery/simply-shared-notes/app"
	"github.com/nicolashery/simply-shared-notes/app/db"
)

//go:embed sql/pragmas.sql
var pragmasSQL string

//go:embed all:dist
var distFS embed.FS

func run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	config, err := app.NewConfig()
	if err != nil {
		return err
	}

	var h slog.Handler
	if config.IsDev {
		h = slog.NewTextHandler(os.Stdout, nil)
	} else {
		h = slog.NewJSONHandler(os.Stdout, nil)
	}
	logger := slog.New(h)

	var assetsConfig app.AssetsConfig
	if config.IsDev {
		assetsConfig, err = app.DevAssets()
		if err != nil {
			return err
		}
		logger.Info("using dev assets, make sure Vite is running")
	} else {
		assetsConfig, err = app.ProdAssets(distFS)
		if err != nil {
			return err
		}
		logger.Info("using prod assets")
	}

	dbConn, err := db.InitDB(ctx, config.DatabasePath(), pragmasSQL)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	queries := db.New(dbConn)

	s := app.NewServer(logger, queries, assetsConfig)

	return app.RunServer(ctx, s, logger, config.Port)
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
