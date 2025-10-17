package server

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/nicolashery/simply-shared-notes/app/config"
	"github.com/nicolashery/simply-shared-notes/app/rctx"
	"github.com/nicolashery/simply-shared-notes/app/views/pages"
	"github.com/nicolashery/simply-shared-notes/app/vite"
)

func NewMaintenance(
	cfg *config.Config,
	logger *slog.Logger,
	vite *vite.Vite,
	sessionStore *sessions.CookieStore,
	i18nBundle *i18n.Bundle,
) http.Handler {
	logger.Warn("running in maintenance mode")

	r := chi.NewRouter()

	r.Use(
		middleware.Logger,
		middleware.Recoverer,
	)

	StaticDir(r, "/assets", vite.AssetsFS)
	StaticFile(r, "/robots.txt", vite.PublicFS)

	r.Group(func(r chi.Router) {
		r.Use(rctx.ViteCtxMiddleware(vite))
		r.Use(rctx.SessionCtxMiddleware(logger, sessionStore))
		r.Use(rctx.ThemeCtxMiddleware())
		r.Use(rctx.IntlCtxMiddleware(logger, i18nBundle))

		r.Get("/", handleMaintenance(logger))
		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/", http.StatusFound)
		})
	})

	return r
}

func handleMaintenance(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)

		err := pages.Maintenance().Render(r.Context(), w)
		if err != nil {
			logger.Error(
				"failed to render template",
				slog.Any("error", err),
				slog.String("template", "Maintenance"),
			)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}
