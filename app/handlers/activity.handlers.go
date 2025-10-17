package handlers

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nicolashery/simply-shared-notes/app/db"
	"github.com/nicolashery/simply-shared-notes/app/publicid"
	"github.com/nicolashery/simply-shared-notes/app/rctx"
	"github.com/nicolashery/simply-shared-notes/app/views/pages"
)

func handleActivityList(logger *slog.Logger, queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		space := rctx.GetSpace(r.Context())

		activityEntries, err := queries.ListActivity(r.Context(), db.ListActivityParams{
			SpaceID: space.ID,
			Limit:   50,
		})
		if err != nil {
			logger.Error(
				"error getting activity from database",
				slog.Any("error", err),
				slog.Int64("space_id", space.ID),
			)
			http.Error(w, "internal server err", http.StatusInternalServerError)
			return
		}

		memberIDs := collectMemberIDsFromActivity(activityEntries)
		members, err := queries.ListMembersByIDs(r.Context(), db.ListMembersByIDsParams{
			SpaceID:   space.ID,
			MemberIds: memberIDs,
		})
		if err != nil {
			logger.Error(
				"error getting members from database",
				slog.Any("error", err),
				slog.Int64("space_id", space.ID),
			)
			http.Error(w, "internal server err", http.StatusInternalServerError)
			return
		}
		membersByID := memberListToMap(members)

		noteIDs := collectNoteIDsFromActivity(activityEntries)
		notes, err := queries.ListNotesByIDs(r.Context(), db.ListNotesByIDsParams{
			SpaceID: space.ID,
			NoteIds: noteIDs,
		})
		if err != nil {
			logger.Error(
				"error getting notes from database",
				slog.Any("error", err),
				slog.Int64("space_id", space.ID),
			)
			http.Error(w, "internal server err", http.StatusInternalServerError)
			return
		}
		notesByID := noteListToMap(notes)

		err = pages.ActivityList(activityEntries, membersByID, notesByID).Render(r.Context(), w)
		if err != nil {
			logger.Error(
				"failed to render template",
				slog.Any("error", err),
				slog.String("template", "ActivityList"),
			)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}

func handleActivityShow(logger *slog.Logger, queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		activityPublicID := chi.URLParam(r, "activityPublicID")
		if !publicid.IsValidPublicID(activityPublicID) {
			http.Error(w, "activity not found", http.StatusNotFound)
			return
		}

		space := rctx.GetSpace(r.Context())
		activity, err := queries.GetActivityByPublicID(r.Context(), db.GetActivityByPublicIDParams{
			SpaceID:  space.ID,
			PublicID: activityPublicID,
		})
		if err != nil {
			http.Error(w, "activity not found", http.StatusNotFound)
			return
		}

		memberIDs := collectMemberIDsFromActivity([]db.Activity{activity})
		members, err := queries.ListMembersByIDs(r.Context(), db.ListMembersByIDsParams{
			SpaceID:   space.ID,
			MemberIds: memberIDs,
		})
		if err != nil {
			logger.Error(
				"error getting members from database",
				slog.Any("error", err),
				slog.Int64("space_id", space.ID),
			)
			http.Error(w, "internal server err", http.StatusInternalServerError)
			return
		}
		membersByID := memberListToMap(members)

		noteIDs := collectNoteIDsFromActivity([]db.Activity{activity})
		notes, err := queries.ListNotesByIDs(r.Context(), db.ListNotesByIDsParams{
			SpaceID: space.ID,
			NoteIds: noteIDs,
		})
		if err != nil {
			logger.Error(
				"error getting notes from database",
				slog.Any("error", err),
				slog.Int64("space_id", space.ID),
			)
			http.Error(w, "internal server err", http.StatusInternalServerError)
			return
		}
		notesByID := noteListToMap(notes)

		err = pages.ActivityShow(activity, membersByID, notesByID).Render(r.Context(), w)
		if err != nil {
			logger.Error(
				"failed to render template",
				slog.Any("error", err),
				slog.String("template", "ActivityShow"),
			)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}
