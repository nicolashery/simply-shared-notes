package handlers

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/nicolashery/simply-shared-notes/db"
)

func RegisterRoutes(r chi.Router, logger *slog.Logger, queries *db.Queries) {
	r.Get("/", handleHome())
	r.Get("/new", handleSpacesNew())
	r.Post("/new", handleSpacesCreate(logger, queries))
	r.Get("/s/{token}", handleSpacesShow(queries))
}
