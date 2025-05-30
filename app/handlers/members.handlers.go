package handlers

import (
	"log/slog"
	"net/http"

	"github.com/nicolashery/simply-shared-notes/app/views/pages"
)

func handleMembersList(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := pages.MembersList().Render(r.Context(), w)
		if err != nil {
			logger.Error(
				"failed to render template",
				slog.Any("error", err),
				slog.String("template", "MembersList"),
			)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}

func handleMembersEdit(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := pages.MembersEdit().Render(r.Context(), w)
		if err != nil {
			logger.Error(
				"failed to render template",
				slog.Any("error", err),
				slog.String("template", "MembersEdit"),
			)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}
