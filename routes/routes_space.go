package routes

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nicolashery/simply-shared-notes/db"
	"github.com/nicolashery/simply-shared-notes/views"
)

func handleGetSpace(logger *slog.Logger, queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("handling get space")

		accessToken := chi.URLParam(r, "token")
		space, err := queries.GetSpaceByAccessToken(r.Context(), accessToken)
		if err != nil {
			http.Error(w, "space not found", http.StatusNotFound)
			return
		}

		views.PageSpace(&space).Render(r.Context(), w)
	}
}
