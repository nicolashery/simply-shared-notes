package handlers

import (
	"log/slog"
	"net/http"

	"github.com/nicolashery/simply-shared-notes/app/db"
	"github.com/nicolashery/simply-shared-notes/app/rctx"
	"github.com/nicolashery/simply-shared-notes/app/views/pages"
)

func handleMembersList(logger *slog.Logger, queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		space := rctx.GetSpace(r.Context())
		members, err := queries.ListMembers(r.Context(), space.ID)
		if err != nil {
			logger.Error(
				"error getting space members from database",
				slog.Any("error", err),
				slog.Int64("space_id", space.ID),
			)
			http.Error(w, "internal server err", http.StatusInternalServerError)
			return
		}

		err = pages.MembersList(members).Render(r.Context(), w)
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
