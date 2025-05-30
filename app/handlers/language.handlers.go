package handlers

import (
	"log/slog"
	"net/http"

	"github.com/nicolashery/simply-shared-notes/app/views/pages"
)

func handleLanguageSelect(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := pages.LanguageSelect().Render(r.Context(), w)
		if err != nil {
			logger.Error(
				"failed to render template",
				slog.Any("error", err),
				slog.String("template", "LanguageSelect"),
			)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}
