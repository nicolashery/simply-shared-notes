package routes

import (
	"log/slog"
	"net/http"

	"github.com/nicolashery/simply-shared-notes/views"
)

func handleHome(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("rendering home")
		views.PageHome().Render(r.Context(), w)
	}
}
