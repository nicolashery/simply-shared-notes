package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/nicolashery/simply-shared-notes/app/db"
	"github.com/nicolashery/simply-shared-notes/app/rctx"
	"github.com/nicolashery/simply-shared-notes/app/session"
	"github.com/nicolashery/simply-shared-notes/app/views/pages"
)

func handleIdentitySelect(logger *slog.Logger, queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		access := rctx.GetAccess(r.Context())
		if access.IsView() {
			http.Redirect(w, r, fmt.Sprintf("/s/%s/identity", access.Token), http.StatusSeeOther)
			return
		}

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

		err = pages.IdentitySelect(members).Render(r.Context(), w)
		if err != nil {
			logger.Error(
				"failed to render template",
				slog.Any("error", err),
				slog.String("template", "IdentitySelect"),
			)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}

type SelectIdentityForm struct {
	MemberPublicId string
}

func parseSelectIdentityForm(r *http.Request, f *SelectIdentityForm) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	f.MemberPublicId = strings.Trim(r.Form.Get("member"), " ")

	return nil
}

func handleIdentitySet(logger *slog.Logger, queries *db.Queries, sessionStore *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		access := rctx.GetAccess(r.Context())
		if access.IsView() {
			http.Redirect(w, r, fmt.Sprintf("/s/%s", access.Token), http.StatusSeeOther)
			return
		}

		var form SelectIdentityForm
		err := parseSelectIdentityForm(r, &form)
		if err != nil {
			http.Error(w, "failed to parse form", http.StatusBadRequest)
			return
		}

		space := rctx.GetSpace(r.Context())
		member, err := queries.GetMemberByPublicID(r.Context(), db.GetMemberByPublicIDParams{
			SpaceID:  space.ID,
			PublicID: form.MemberPublicId,
		})
		if err != nil {
			http.Error(w, "member not found", http.StatusBadRequest)
			return
		}

		sess := rctx.GetSession(r.Context())
		sess.Values[session.IdentityKey] = member.ID
		sess.AddFlash(session.FlashMessage{
			Type:    session.FlashType_Info,
			Content: fmt.Sprintf("%s, welcome to the space %s!", member.Name, space.Name),
		})
		err = sess.Save(r, w)
		if err != nil {
			logger.Error("failed to save session", slog.Any("error", err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/s/%s", access.Token), http.StatusSeeOther)
	}
}

func handleIdentityClear(logger *slog.Logger, sessionStore *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess := rctx.GetSession(r.Context())
		delete(sess.Values, session.IdentityKey)
		err := sess.Save(r, w)
		if err != nil {
			logger.Error("failed to save session", slog.Any("error", err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		access := rctx.GetAccess(r.Context())

		http.Redirect(w, r, fmt.Sprintf("/s/%s/identity", access.Token), http.StatusSeeOther)
	}
}
