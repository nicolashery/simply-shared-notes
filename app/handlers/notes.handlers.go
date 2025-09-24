package handlers

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/nicolashery/simply-shared-notes/app/db"
	"github.com/nicolashery/simply-shared-notes/app/forms"
	"github.com/nicolashery/simply-shared-notes/app/markdown"
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

		memberIDs := collectCreatedUpdatedByIDsFromNotes(notes)
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

		err = pages.NotesList(notes, membersByID).Render(r.Context(), w)
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

func handleNotesCreate(logger *slog.Logger, sqlDB *sql.DB, queries *db.Queries) http.HandlerFunc {
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

		activityPublicID, err := publicid.Generate()
		if err != nil {
			logger.Error("error generating activity public ID", slog.Any("error", err))
			http.Error(w, "error creating note", http.StatusInternalServerError)
			return
		}

		identity := rctx.GetIdentity(r.Context())
		space := rctx.GetSpace(r.Context())
		now := time.Now().UTC()

		tx, err := sqlDB.Begin()
		if err != nil {
			logger.Error("error starting transaction", slog.Any("error", err))
			http.Error(w, "error creating note", http.StatusInternalServerError)
			return
		}
		defer func() {
			if rbErr := tx.Rollback(); rbErr != nil && err == nil {
				logger.Error("error rolling back transaction", slog.Any("error", rbErr))
			}
		}()
		qtx := queries.WithTx(tx)

		note, err := qtx.CreateNote(
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

		_, err = qtx.CreateActivity(r.Context(), db.CreateActivityParams{
			CreatedAt:  now,
			SpaceID:    space.ID,
			PublicID:   activityPublicID,
			MemberID:   sql.NullInt64{Int64: identity.Member.ID, Valid: true},
			Action:     db.ActivityAction_Create,
			EntityType: db.ActivityEntity_Note,
			EntityID:   sql.NullInt64{Int64: note.ID, Valid: true},
		})
		if err != nil {
			logger.Error("error creating activity in database", slog.Any("error", err))
			http.Error(w, "error creating note", http.StatusInternalServerError)
			return
		}

		err = tx.Commit()
		if err != nil {
			logger.Error("error committing transaction", slog.Any("error", err))
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

func handleNotesShow(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		note := rctx.GetNote(r.Context())
		renderedMarkdown, err := markdown.Render(note.Content)
		if err != nil {
			logger.Error("failed to render Markdown", slog.Any("error", err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		err = pages.NotesShow(renderedMarkdown).Render(r.Context(), w)
		if err != nil {
			logger.Error(
				"failed to render template",
				slog.Any("error", err),
				slog.String("template", "NotesShow"),
			)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}

func handleNotesEdit(logger *slog.Logger, queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		note := rctx.GetNote(r.Context())
		form := forms.UpdateNote{
			Title:   note.Title,
			Content: note.Content,
		}

		space := rctx.GetSpace(r.Context())
		memberIDs := collectCreatedUpdatedByIDsFromNotes([]db.Note{*note})
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

		err = pages.NotesEdit(&form, forms.EmptyErrors(), membersByID).Render(r.Context(), w)
		if err != nil {
			logger.Error(
				"failed to render template",
				slog.Any("error", err),
				slog.String("template", "NotesEdit"),
			)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}

func handleNotesUpdate(logger *slog.Logger, sqlDB *sql.DB, queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form, errors := forms.ParseUpdateNote(r)
		if errors != nil {
			space := rctx.GetSpace(r.Context())
			note := rctx.GetNote(r.Context())
			memberIDs := collectCreatedUpdatedByIDsFromNotes([]db.Note{*note})
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

			w.WriteHeader(http.StatusUnprocessableEntity)
			err = pages.NotesEdit(&form, errors, membersByID).Render(r.Context(), w)
			if err != nil {
				logger.Error(
					"failed to render template",
					slog.Any("error", err),
					slog.String("template", "NotesEdit"),
				)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		activityPublicID, err := publicid.Generate()
		if err != nil {
			logger.Error("error generating activity public ID", slog.Any("error", err))
			http.Error(w, "error updating note", http.StatusInternalServerError)
			return
		}

		now := time.Now().UTC()
		identity := rctx.GetIdentity(r.Context())
		note := rctx.GetNote(r.Context())
		space := rctx.GetSpace(r.Context())

		tx, err := sqlDB.Begin()
		if err != nil {
			logger.Error("error starting transaction", slog.Any("error", err))
			http.Error(w, "error updating note", http.StatusInternalServerError)
			return
		}
		defer func() {
			if rbErr := tx.Rollback(); rbErr != nil && err == nil {
				logger.Error("error rolling back transaction", slog.Any("error", rbErr))
			}
		}()
		qtx := queries.WithTx(tx)

		_, err = qtx.UpdateNote(
			r.Context(),
			db.UpdateNoteParams{
				UpdatedAt: now,
				UpdatedBy: sql.NullInt64{Int64: identity.Member.ID, Valid: true},
				Title:     form.Title,
				Content:   form.Content,
				NoteID:    note.ID,
			},
		)
		if err != nil {
			logger.Error("error updating note in database", slog.Any("error", err))
			http.Error(w, "error updating note", http.StatusInternalServerError)
			return
		}

		_, err = qtx.CreateActivity(r.Context(), db.CreateActivityParams{
			CreatedAt:  now,
			SpaceID:    space.ID,
			PublicID:   activityPublicID,
			MemberID:   sql.NullInt64{Int64: identity.Member.ID, Valid: true},
			Action:     db.ActivityAction_Update,
			EntityType: db.ActivityEntity_Note,
			EntityID:   sql.NullInt64{Int64: note.ID, Valid: true},
		})
		if err != nil {
			logger.Error("error creating activity in database", slog.Any("error", err))
			http.Error(w, "error updating note", http.StatusInternalServerError)
			return
		}

		err = tx.Commit()
		if err != nil {
			logger.Error("error committing transaction", slog.Any("error", err))
			http.Error(w, "error updating note", http.StatusInternalServerError)
			return
		}

		sess := rctx.GetSession(r.Context())
		sess.AddFlash(session.FlashMessage{
			Type:    session.FlashType_Success,
			Content: "Changes saved successfully",
		})
		err = sess.Save(r, w)
		if err != nil {
			logger.Error("failed to save session", slog.Any("error", err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		access := rctx.GetAccess(r.Context())
		http.Redirect(w, r, fmt.Sprintf("/s/%s/notes/%s", access.Token, note.PublicID), http.StatusSeeOther)
	}
}

func handleNotesDeleteConfirm(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := pages.NotesDelete().Render(r.Context(), w)
		if err != nil {
			logger.Error(
				"failed to render template",
				slog.Any("error", err),
				slog.String("template", "NotesDelete"),
			)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}

func handleNotesDelete(logger *slog.Logger, sqlDB *sql.DB, queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now().UTC()
		identity := rctx.GetIdentity(r.Context())
		note := rctx.GetNote(r.Context())
		space := rctx.GetSpace(r.Context())

		activityPublicID, err := publicid.Generate()
		if err != nil {
			logger.Error("error generating activity public ID", slog.Any("error", err))
			http.Error(w, "error deleting note", http.StatusInternalServerError)
			return
		}

		tx, err := sqlDB.Begin()
		if err != nil {
			logger.Error("error starting transaction", slog.Any("error", err))
			http.Error(w, "error deleting note", http.StatusInternalServerError)
			return
		}
		defer func() {
			if rbErr := tx.Rollback(); rbErr != nil && err == nil {
				logger.Error("error rolling back transaction", slog.Any("error", rbErr))
			}
		}()
		qtx := queries.WithTx(tx)

		err = qtx.DeleteNote(r.Context(), note.ID)
		if err != nil {
			logger.Error("error deleting note in database", slog.Any("error", err))
			http.Error(w, "error deleting note", http.StatusInternalServerError)
			return
		}

		_, err = qtx.CreateActivity(r.Context(), db.CreateActivityParams{
			CreatedAt:  now,
			SpaceID:    space.ID,
			PublicID:   activityPublicID,
			MemberID:   sql.NullInt64{Int64: identity.Member.ID, Valid: true},
			Action:     db.ActivityAction_Delete,
			EntityType: db.ActivityEntity_Note,
			EntityID:   sql.NullInt64{Valid: false},
		})
		if err != nil {
			logger.Error("error creating activity in database", slog.Any("error", err))
			http.Error(w, "error deleting note", http.StatusInternalServerError)
			return
		}

		err = qtx.SetActivityEntityIDToNull(r.Context(), db.SetActivityEntityIDToNullParams{
			EntityType: db.ActivityEntity_Note,
			EntityID:   sql.NullInt64{Int64: note.ID, Valid: true},
		})
		if err != nil {
			logger.Error("error updating activity in database", slog.Any("error", err))
			http.Error(w, "error deleting note", http.StatusInternalServerError)
			return
		}

		err = tx.Commit()
		if err != nil {
			logger.Error("error committing transaction", slog.Any("error", err))
			http.Error(w, "error deleting note", http.StatusInternalServerError)
			return
		}

		sess := rctx.GetSession(r.Context())
		sess.AddFlash(session.FlashMessage{
			Type:    session.FlashType_Success,
			Content: fmt.Sprintf("Deleted note: %s", note.Title),
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
