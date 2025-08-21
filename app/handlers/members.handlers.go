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

func handleMembersNew(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var form forms.CreateMember

		err := pages.MembersNew(&form, forms.EmptyErrors()).Render(r.Context(), w)
		if err != nil {
			logger.Error(
				"failed to render template",
				slog.Any("error", err),
				slog.String("template", "MembersNew"),
			)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}

func handleMembersCreate(logger *slog.Logger, queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form, errors := forms.ParseCreateMember(r)
		if errors != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			err := pages.MembersNew(&form, errors).Render(r.Context(), w)
			if err != nil {
				logger.Error(
					"failed to render template",
					slog.Any("error", err),
					slog.String("template", "MembersNew"),
				)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		memberPublicID, err := publicid.Generate()
		if err != nil {
			logger.Error("error generating member public ID", slog.Any("error", err))
			http.Error(w, "error creating member", http.StatusInternalServerError)
			return
		}
		identity := rctx.GetIdentity(r.Context())
		space := rctx.GetSpace(r.Context())
		now := time.Now().UTC()

		member, err := queries.CreateMember(
			r.Context(),
			db.CreateMemberParams{
				CreatedAt: now,
				UpdatedAt: now,
				CreatedBy: sql.NullInt64{Int64: identity.Member.ID, Valid: true},
				UpdatedBy: sql.NullInt64{Int64: identity.Member.ID, Valid: true},
				SpaceID:   space.ID,
				PublicID:  memberPublicID,
				Name:      form.Name,
			},
		)
		if err != nil {
			logger.Error("error creating member in database", slog.Any("error", err))
			http.Error(w, "error creating member", http.StatusInternalServerError)
			return
		}

		sess := rctx.GetSession(r.Context())
		sess.Values[session.IdentityKey] = member.ID
		sess.AddFlash(session.FlashMessage{
			Type:    session.FlashType_Info,
			Content: fmt.Sprintf("Added new member: %s", member.Name),
		})
		err = sess.Save(r, w)
		if err != nil {
			logger.Error("failed to save session", slog.Any("error", err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		access := rctx.GetAccess(r.Context())

		http.Redirect(w, r, fmt.Sprintf("/s/%s/members", access.Token), http.StatusSeeOther)
	}
}

func handleMembersEdit(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		member := rctx.GetMember(r.Context())
		form := forms.UpdateMember{
			Name: member.Name,
		}

		err := pages.MembersEdit(&form, forms.EmptyErrors()).Render(r.Context(), w)
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

func handleMembersUpdate(logger *slog.Logger, queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form, errors := forms.ParseUpdateMember(r)
		if errors != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			err := pages.MembersEdit(&form, errors).Render(r.Context(), w)
			if err != nil {
				logger.Error(
					"failed to render template",
					slog.Any("error", err),
					slog.String("template", "MembersEdit"),
				)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		now := time.Now().UTC()
		identity := rctx.GetIdentity(r.Context())
		member := rctx.GetMember(r.Context())
		memberUpdated, err := queries.UpdateMember(
			r.Context(),
			db.UpdateMemberParams{
				UpdatedAt: now,
				UpdatedBy: sql.NullInt64{Int64: identity.Member.ID, Valid: true},
				Name:      form.Name,
				MemberID:  member.ID,
			},
		)
		if err != nil {
			logger.Error("error updating member in database", slog.Any("error", err))
			http.Error(w, "error updating member", http.StatusInternalServerError)
			return
		}

		sess := rctx.GetSession(r.Context())
		sess.AddFlash(session.FlashMessage{
			Type:    session.FlashType_Success,
			Content: fmt.Sprintf("Saved changes for: %s", memberUpdated.Name),
		})
		err = sess.Save(r, w)
		if err != nil {
			logger.Error("failed to save session", slog.Any("error", err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		access := rctx.GetAccess(r.Context())
		http.Redirect(w, r, fmt.Sprintf("/s/%s/members", access.Token), http.StatusSeeOther)
	}
}
