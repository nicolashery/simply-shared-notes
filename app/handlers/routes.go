package handlers

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/nicolashery/simply-shared-notes/app/db"
	"github.com/nicolashery/simply-shared-notes/app/rctx"
)

func RegisterRoutes(r chi.Router, logger *slog.Logger, queries *db.Queries) {
	r.Get("/", handleHome())
	r.Get("/new", handleSpacesNew())
	r.Post("/new", handleSpacesCreate(logger, queries))

	r.Route("/s/{token}", func(r chi.Router) {
		r.Use(rctx.SpaceCtxMiddleware(queries))
		r.Use(rctx.AccessCtxMiddleware(logger))

		r.Get("/", handleSpacesShow())
	})
}
