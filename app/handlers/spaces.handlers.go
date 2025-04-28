package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/nicolashery/simply-shared-notes/app/access"
	"github.com/nicolashery/simply-shared-notes/app/config"
	"github.com/nicolashery/simply-shared-notes/app/db"
	"github.com/nicolashery/simply-shared-notes/app/forms"
	"github.com/nicolashery/simply-shared-notes/app/publicid"
	"github.com/nicolashery/simply-shared-notes/app/session"
	"github.com/nicolashery/simply-shared-notes/app/views/pages"
)

func handleSpacesNew(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requiresCode := cfg.RequiresInvitationCode()

		var form forms.CreateSpace
		if requiresCode {
			form.Code = r.URL.Query().Get("code")
		}

		pages.SpacesNew(requiresCode, &form, forms.EmptyErrors()).Render(r.Context(), w)
	}
}

func handleSpacesCreate(cfg *config.Config, logger *slog.Logger, sqlDB *sql.DB, queries *db.Queries, sessionStore *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requiresCode := cfg.RequiresInvitationCode()

		form, errors := forms.ParseCreateSpace(r, requiresCode)
		if errors != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			pages.SpacesNew(requiresCode, &form, errors).Render(r.Context(), w)
			return
		}

		if cfg.RequiresInvitationCode() && form.Code != cfg.InvitationCode {
			http.Error(w, "invitation code required", http.StatusBadRequest)
			return
		}

		tokens, err := access.GenerateAccessTokens()
		if err != nil {
			logger.Error("error generating access tokens", slog.Any("error", err))
			http.Error(w, "error creating space", http.StatusInternalServerError)
			return
		}

		memberPublicId, err := publicid.Generate()
		if err != nil {
			logger.Error("error generating member public ID", slog.Any("error", err))
			http.Error(w, "error creating space", http.StatusInternalServerError)
			return
		}

		now := time.Now().UTC()

		space, member, err := createSpaceAndFirstMember(
			r.Context(),
			sqlDB,
			queries,
			form,
			now,
			tokens,
			memberPublicId,
		)
		if err != nil {
			logger.Error("error creating space and first member in database", slog.Any("error", err))
			http.Error(w, "error creating space", http.StatusInternalServerError)
			return
		}

		sess, err := sessionStore.Get(r, session.CookieName)
		if err != nil {
			logger.Error("failed to get session", slog.Any("error", err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

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

		http.Redirect(w, r, fmt.Sprintf("/s/%s", space.AdminToken), http.StatusSeeOther)
	}
}

func createSpaceAndFirstMember(
	ctx context.Context,
	sqlDB *sql.DB,
	queries *db.Queries,
	form forms.CreateSpace,
	now time.Time,
	tokens access.AccessTokens,
	memberPublicId string,
) (*db.Space, *db.Member, error) {
	tx, err := sqlDB.Begin()
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback()
	qtx := queries.WithTx(tx)

	space, err := qtx.CreateSpace(ctx, db.CreateSpaceParams{
		CreatedAt:  now,
		UpdatedAt:  now,
		Name:       form.Name,
		Email:      form.Email,
		AdminToken: tokens.AdminToken,
		EditToken:  tokens.EditToken,
		ViewToken:  tokens.ViewToken,
	})
	if err != nil {
		return nil, nil, err
	}

	member, err := qtx.CreateMember(ctx, db.CreateMemberParams{
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: sql.NullInt64{Valid: false},
		UpdatedBy: sql.NullInt64{Valid: false},
		SpaceID:   space.ID,
		PublicID:  memberPublicId,
		Name:      form.Identity,
	})
	if err != nil {
		return nil, nil, err
	}

	err = qtx.UpdateSpaceCreatedBy(ctx, db.UpdateSpaceCreatedByParams{
		SpaceID:   space.ID,
		CreatedBy: sql.NullInt64{Int64: member.ID, Valid: true},
	})
	if err != nil {
		return nil, nil, err
	}

	err = qtx.UpdateMemberCreatedBy(ctx, db.UpdateMemberCreatedByParams{
		MemberID:  member.ID,
		CreatedBy: sql.NullInt64{Int64: member.ID, Valid: true},
	})
	if err != nil {
		return nil, nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, nil, err
	}

	return &space, &member, nil
}

func handleSpacesShow() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pages.SpacesShow().Render(r.Context(), w)
	}
}
