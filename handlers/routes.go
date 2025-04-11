package handlers

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/nicolashery/simply-shared-notes/db"
)

func RegisterRoutes(r chi.Router, logger *slog.Logger, queries *db.Queries) {
	r.Get("/", handleHome(logger))
	r.Get("/s/{token}", handleGetSpace(logger, queries))
}
