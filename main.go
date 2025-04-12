package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"os/signal"
	"strconv"

	"github.com/nicolashery/simply-shared-notes/db"
	"github.com/nicolashery/simply-shared-notes/server"
	"github.com/olivere/vite"
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

	isDevAssets := false
	if isDevStr := os.Getenv("DEV"); isDevStr == "true" {
		isDevAssets = true
	}
	var viteFS fs.FS
	var viteAssetsFS fs.FS
	if !isDevAssets {
		var err error
		viteFS, err = fs.Sub(distFS, "dist")
		if err != nil {
			return fmt.Errorf("failed to create sub-filesystem for 'dist' directory: %w", err)
		}
		viteAssetsFS, err = fs.Sub(viteFS, "assets")
		if err != nil {
			return fmt.Errorf("failed to create sub-filesystem for 'dist/assets' directory: %w", err)
		}

	}
	viteFragment, err := vite.HTMLFragment(vite.Config{
		FS:           viteFS,
		IsDev:        isDevAssets,
		ViteTemplate: vite.None,
		ViteEntry:    "assets/app.js",
	})
	if err != nil {
		return fmt.Errorf("failed to instantiate Vite fragment: %w", err)
	}

	var assetsConfig server.AssetsConfig
	if isDevAssets {
		assetsConfig = server.AssetsConfig{
			AssetsFS:     os.DirFS("./assets"),
			AssetsPath:   "/assets",
			PublicFS:     os.DirFS("./public"),
			ViteFragment: viteFragment,
		}
		logger.Info("using dev assets, make sure Vite is running")
	} else {
		assetsConfig = server.AssetsConfig{
			AssetsFS:     viteAssetsFS,
			AssetsPath:   "/assets",
			PublicFS:     viteFS,
			ViteFragment: viteFragment,
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
