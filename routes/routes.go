package routes

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, logger *slog.Logger) {
	r.Get("/", handleHome(logger))
}
