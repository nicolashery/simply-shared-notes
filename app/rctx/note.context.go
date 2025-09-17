package rctx

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nicolashery/simply-shared-notes/app/db"
	"github.com/nicolashery/simply-shared-notes/app/publicid"
)

func NoteCtxMiddleware(queries *db.Queries) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			notePublicID := chi.URLParam(r, "notePublicID")
			if !publicid.IsValidPublicID(notePublicID) {
				http.Error(w, "note not found", http.StatusNotFound)
				return
			}

			space := GetSpace(r.Context())
			note, err := queries.GetNoteByPublicID(r.Context(), db.GetNoteByPublicIDParams{
				SpaceID:  space.ID,
				PublicID: notePublicID,
			})
			if err != nil {
				http.Error(w, "member not found", http.StatusNotFound)
				return
			}

			ctx := context.WithValue(r.Context(), noteContextKey, &note)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetNote(ctx context.Context) *db.Note {
	note, ok := ctx.Value(noteContextKey).(*db.Note)
	if !ok {
		panic("note not found in context, make sure to use middleware")
	}

	return note
}
