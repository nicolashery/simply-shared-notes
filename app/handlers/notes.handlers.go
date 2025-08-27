package handlers

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/nicolashery/simply-shared-notes/app/db"
	"github.com/nicolashery/simply-shared-notes/app/forms"
	"github.com/nicolashery/simply-shared-notes/app/publicid"
	"github.com/nicolashery/simply-shared-notes/app/rctx"
	"github.com/nicolashery/simply-shared-notes/app/session"
	"github.com/nicolashery/simply-shared-notes/app/views/pages"
)

func handleNotesList(logger *slog.Logger, queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		space := rctx.GetSpace(r.Context())
		notes, err := queries.ListNotes(r.Context(), space.ID)
		if err != nil {
			logger.Error(
				"error getting notes from database",
				slog.Any("error", err),
				slog.Int64("space_id", space.ID),
			)
			http.Error(w, "internal server err", http.StatusInternalServerError)
			return
		}

		err = pages.NotesList(notes).Render(r.Context(), w)
		if err != nil {
			logger.Error(
				"failed to render template",
				slog.Any("error", err),
				slog.String("template", "NotesList"),
			)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}

func handleNotesNew(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var form forms.CreateNote

		err := pages.NotesNew(&form, forms.EmptyErrors()).Render(r.Context(), w)
		if err != nil {
			logger.Error(
				"failed to render template",
				slog.Any("error", err),
				slog.String("template", "NotesNew"),
			)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}

func handleNotesCreate(logger *slog.Logger, queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form, errors := forms.ParseCreateNote(r)
		if errors != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			err := pages.NotesNew(&form, errors).Render(r.Context(), w)
			if err != nil {
				logger.Error(
					"failed to render template",
					slog.Any("error", err),
					slog.String("template", "NotesNew"),
				)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		notePublicID, err := publicid.Generate()
		if err != nil {
			logger.Error("error generating note public ID", slog.Any("error", err))
			http.Error(w, "error creating note", http.StatusInternalServerError)
			return
		}
		identity := rctx.GetIdentity(r.Context())
		space := rctx.GetSpace(r.Context())
		now := time.Now().UTC()

		note, err := queries.CreateNote(
			r.Context(),
			db.CreateNoteParams{
				CreatedAt: now,
				UpdatedAt: now,
				CreatedBy: sql.NullInt64{Int64: identity.Member.ID, Valid: true},
				UpdatedBy: sql.NullInt64{Int64: identity.Member.ID, Valid: true},
				SpaceID:   space.ID,
				PublicID:  notePublicID,
				Title:     form.Title,
				Content:   form.Content,
			},
		)
		if err != nil {
			logger.Error("error creating note in database", slog.Any("error", err))
			http.Error(w, "error creating note", http.StatusInternalServerError)
			return
		}

		sess := rctx.GetSession(r.Context())
		sess.AddFlash(session.FlashMessage{
			Type:    session.FlashType_Info,
			Content: fmt.Sprintf("Created new note: %s", note.Title),
		})
		err = sess.Save(r, w)
		if err != nil {
			logger.Error("failed to save session", slog.Any("error", err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		access := rctx.GetAccess(r.Context())

		http.Redirect(w, r, fmt.Sprintf("/s/%s/notes", access.Token), http.StatusSeeOther)
	}
}
