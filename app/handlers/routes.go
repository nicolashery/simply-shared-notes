package handlers

import (
	"database/sql"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/nicolashery/simply-shared-notes/app/config"
	"github.com/nicolashery/simply-shared-notes/app/db"
	"github.com/nicolashery/simply-shared-notes/app/rctx"
)

func RegisterRoutes(r chi.Router, cfg *config.Config, logger *slog.Logger, conn *sql.DB, queries *db.Queries) {
	r.Get("/", handleHome())
	r.Get("/new", handleSpacesNew(cfg))
	r.Post("/new", handleSpacesCreate(cfg, logger, conn, queries))

	r.Route("/s/{token}", func(r chi.Router) {
		r.Use(rctx.SpaceCtxMiddleware(queries))
		r.Use(rctx.AccessCtxMiddleware(logger))

		r.Get("/", handleSpacesShow())
	})
}
