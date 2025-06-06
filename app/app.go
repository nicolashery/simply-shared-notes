package app

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/nicolashery/simply-shared-notes/app/config"
	"github.com/nicolashery/simply-shared-notes/app/db"
	"github.com/nicolashery/simply-shared-notes/app/server"
	"github.com/nicolashery/simply-shared-notes/app/session"
	"github.com/nicolashery/simply-shared-notes/app/vite"
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

	vite, err := vite.New(logger, cfg.IsDev, distFS)
	if err != nil {
		return err
	}

	sqlDB, err := db.InitDB(ctx, cfg.DatabasePath(), pragmasSQL)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	queries := db.New(sqlDB)

	sessionStore := session.InitStore(cfg.CookieSecret, cfg.IsDev)

	s := server.New(cfg, logger, sqlDB, queries, vite, sessionStore)

	return server.Run(ctx, s, logger, cfg.Port)
}
