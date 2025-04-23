package server

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"
	"github.com/nicolashery/simply-shared-notes/app/assets"
	"github.com/nicolashery/simply-shared-notes/app/config"
	"github.com/nicolashery/simply-shared-notes/app/db"
	"github.com/nicolashery/simply-shared-notes/app/handlers"
	"github.com/nicolashery/simply-shared-notes/app/rctx"
)

func New(cfg *config.Config, logger *slog.Logger, sqlDB *sql.DB, queries *db.Queries, assetsConfig assets.AssetsConfig, sessionStore *sessions.CookieStore) http.Handler {
	router := chi.NewRouter()

	router.Use(
		middleware.Logger,
		middleware.Recoverer,
		rctx.ViteCtxMiddleware(assetsConfig.ViteFragment),
	)

	handlers.RegisterRoutes(router, cfg, logger, sqlDB, queries, sessionStore)

	StaticDir(router, "/assets", assetsConfig.AssetsFS)
	StaticFile(router, "/robots.txt", assetsConfig.PublicFS)

	return router
}

func Run(ctx context.Context, handler http.Handler, logger *slog.Logger, port int) error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}

	go func() {
		<-ctx.Done()
		logger.Info("shutting down")
		srv.Shutdown(context.Background())
	}()

	logger.Info(fmt.Sprintf("listening on port %d", port))

	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
