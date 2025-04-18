package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nicolashery/simply-shared-notes/app/db"
	"github.com/nicolashery/simply-shared-notes/app/handlers"
	"github.com/nicolashery/simply-shared-notes/app/middlewares"
)

func NewServer(logger *slog.Logger, queries *db.Queries, assetsConfig AssetsConfig) http.Handler {
	router := chi.NewRouter()

	router.Use(
		middleware.Logger,
		middleware.Recoverer,
		middlewares.ViteCtx(assetsConfig.ViteFragment),
	)

	handlers.RegisterRoutes(router, logger, queries)

	StaticDir(router, "/assets", assetsConfig.AssetsFS)
	StaticFile(router, "/robots.txt", assetsConfig.PublicFS)

	return router
}

func RunServer(ctx context.Context, handler http.Handler, logger *slog.Logger, port int) error {
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
