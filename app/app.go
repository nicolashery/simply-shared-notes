package app

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/nicolashery/simply-shared-notes/app/assets"
	"github.com/nicolashery/simply-shared-notes/app/config"
	"github.com/nicolashery/simply-shared-notes/app/db"
	"github.com/nicolashery/simply-shared-notes/app/server"
)

func Run(ctx context.Context, distFS embed.FS, pragmasSQL string) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	cfg, err := config.New()
	if err != nil {
		return err
	}

	var h slog.Handler
	if cfg.IsDev {
		h = slog.NewTextHandler(os.Stdout, nil)
	} else {
		h = slog.NewJSONHandler(os.Stdout, nil)
	}
	logger := slog.New(h)

	var assetsConfig assets.AssetsConfig
	if cfg.IsDev {
		assetsConfig, err = assets.DevAssets()
		if err != nil {
			return err
		}
		logger.Info("using dev assets, make sure Vite is running")
	} else {
		assetsConfig, err = assets.ProdAssets(distFS)
		if err != nil {
			return err
		}
		logger.Info("using prod assets")
	}

	dbConn, err := db.InitDB(ctx, cfg.DatabasePath(), pragmasSQL)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	queries := db.New(dbConn)

	s := server.New(cfg, logger, dbConn, queries, assetsConfig)

	return server.Run(ctx, s, logger, cfg.Port)
}
