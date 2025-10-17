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
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/nicolashery/simply-shared-notes/app/config"
	"github.com/nicolashery/simply-shared-notes/app/db"
	"github.com/nicolashery/simply-shared-notes/app/email"
	"github.com/nicolashery/simply-shared-notes/app/handlers"
	"github.com/nicolashery/simply-shared-notes/app/vite"
)

func New(
	cfg *config.Config,
	logger *slog.Logger,
	sqlDB *sql.DB,
	queries *db.Queries,
	vite *vite.Vite,
	sessionStore *sessions.CookieStore,
	email *email.Email,
	i18nBundle *i18n.Bundle,
) http.Handler {
	r := chi.NewRouter()

	r.Use(
		middleware.Logger,
		middleware.Recoverer,
	)

	StaticDir(r, "/assets", vite.AssetsFS)
	StaticFile(r, "/robots.txt", vite.PublicFS)

	r.Group(func(r chi.Router) {
		handlers.RegisterRoutes(r, cfg, logger, sqlDB, queries, vite, sessionStore, email, i18nBundle)
	})

	return r
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
