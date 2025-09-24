package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/nicolashery/simply-shared-notes/app/access"
	"github.com/nicolashery/simply-shared-notes/app/config"
	"github.com/nicolashery/simply-shared-notes/app/db"
	"github.com/nicolashery/simply-shared-notes/app/forms"
	"github.com/nicolashery/simply-shared-notes/app/publicid"
	"github.com/nicolashery/simply-shared-notes/app/rctx"
	"github.com/nicolashery/simply-shared-notes/app/session"
	"github.com/nicolashery/simply-shared-notes/app/views/pages"
)

func handleSpacesNew(cfg *config.Config, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requiresCode := cfg.RequiresInvitationCode()

		var form forms.CreateSpace
		if requiresCode {
			form.Code = r.URL.Query().Get("code")
		}

		err := pages.SpacesNew(requiresCode, &form, forms.EmptyErrors()).Render(r.Context(), w)
		if err != nil {
			logger.Error(
				"failed to render template",
				slog.Any("error", err),
				slog.String("template", "SpacesNew"),
			)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}

func handleSpacesCreate(cfg *config.Config, logger *slog.Logger, sqlDB *sql.DB, queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requiresCode := cfg.RequiresInvitationCode()

		form, errors := forms.ParseCreateSpace(r, requiresCode)
		if errors != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			err := pages.SpacesNew(requiresCode, &form, errors).Render(r.Context(), w)
			if err != nil {
				logger.Error(
					"failed to render template",
					slog.Any("error", err),
					slog.String("template", "SpacesNew"),
				)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
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

		memberPublicID, err := publicid.Generate()
		if err != nil {
			logger.Error("error generating member public ID", slog.Any("error", err))
			http.Error(w, "error creating space", http.StatusInternalServerError)
			return
		}

		spaceActivityPublicID, err := publicid.Generate()
		if err != nil {
			logger.Error("error generating activity public ID", slog.Any("error", err))
			http.Error(w, "error creating space", http.StatusInternalServerError)
			return
		}

		memberActivityPublicID, err := publicid.Generate()
		if err != nil {
			logger.Error("error generating activity public ID", slog.Any("error", err))
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
			memberPublicID,
			spaceActivityPublicID,
			memberActivityPublicID,
		)
		if err != nil {
			logger.Error("error creating space and first member in database", slog.Any("error", err))
			http.Error(w, "error creating space", http.StatusInternalServerError)
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
	memberPublicID string,
	spaceActivityPublicID string,
	memberActivityPublicID string,
) (*db.Space, *db.Member, error) {
	tx, err := sqlDB.Begin()
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if rbErr := tx.Rollback(); rbErr != nil && err == nil {
			// Only overwrite main err if it's nil and rollback failed
			err = fmt.Errorf("tx rollback error: %w", rbErr)
		}
	}()
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
		PublicID:  memberPublicID,
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

	_, err = qtx.CreateActivity(ctx, db.CreateActivityParams{
		CreatedAt:  now,
		SpaceID:    space.ID,
		PublicID:   spaceActivityPublicID,
		MemberID:   sql.NullInt64{Int64: member.ID, Valid: true},
		Action:     db.ActivityAction_Create,
		EntityType: db.ActivityEntity_Space,
		EntityID:   sql.NullInt64{Int64: space.ID, Valid: true},
	})
	if err != nil {
		return nil, nil, err
	}

	_, err = qtx.CreateActivity(ctx, db.CreateActivityParams{
		CreatedAt:  now,
		SpaceID:    space.ID,
		PublicID:   memberActivityPublicID,
		MemberID:   sql.NullInt64{Int64: member.ID, Valid: true},
		Action:     db.ActivityAction_Create,
		EntityType: db.ActivityEntity_Member,
		EntityID:   sql.NullInt64{Int64: member.ID, Valid: true},
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

func handleSpacesShow(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := pages.SpacesShow().Render(r.Context(), w)
		if err != nil {
			logger.Error(
				"failed to render template",
				slog.Any("error", err),
				slog.String("template", "SpacesShow"),
			)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}

func handleSpacesEdit(logger *slog.Logger, queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		space := rctx.GetSpace(r.Context())
		form := forms.UpdateSpace{
			Name: space.Name,
		}

		memberIDs := collectCreatedUpdatedByIDsFromSpace(space)
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

		err = pages.SpacesEdit(&form, forms.EmptyErrors(), membersByID).Render(r.Context(), w)
		if err != nil {
			logger.Error(
				"failed to render template",
				slog.Any("error", err),
				slog.String("template", "SpacesEdit"),
			)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}

func handleSpacesUpdate(logger *slog.Logger, sqlDB *sql.DB, queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		form, errors := forms.ParseUpdateSpace(r)
		if errors != nil {
			space := rctx.GetSpace(r.Context())
			memberIDs := collectCreatedUpdatedByIDsFromSpace(space)
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
			err = pages.SpacesEdit(&form, errors, membersByID).Render(r.Context(), w)
			if err != nil {
				logger.Error(
					"failed to render template",
					slog.Any("error", err),
					slog.String("template", "SpacesEdit"),
				)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		activityPublicID, err := publicid.Generate()
		if err != nil {
			logger.Error("error generating activity public ID", slog.Any("error", err))
			http.Error(w, "error updating space", http.StatusInternalServerError)
			return
		}

		now := time.Now().UTC()
		identity := rctx.GetIdentity(r.Context())
		space := rctx.GetSpace(r.Context())

		// Start transaction to update space and create activity atomically
		tx, err := sqlDB.Begin()
		if err != nil {
			logger.Error("error starting transaction", slog.Any("error", err))
			http.Error(w, "error updating space", http.StatusInternalServerError)
			return
		}
		defer func() {
			if rbErr := tx.Rollback(); rbErr != nil && err == nil {
				logger.Error("error rolling back transaction", slog.Any("error", rbErr))
			}
		}()
		qtx := queries.WithTx(tx)

		spaceUpdated, err := qtx.UpdateSpace(
			r.Context(),
			db.UpdateSpaceParams{
				UpdatedAt: now,
				UpdatedBy: sql.NullInt64{Int64: identity.Member.ID, Valid: true},
				Name:      form.Name,
				SpaceID:   space.ID,
			},
		)
		if err != nil {
			logger.Error("error updating space in database", slog.Any("error", err))
			http.Error(w, "error updating space", http.StatusInternalServerError)
			return
		}

		_, err = qtx.CreateActivity(r.Context(), db.CreateActivityParams{
			CreatedAt:  now,
			SpaceID:    space.ID,
			PublicID:   activityPublicID,
			MemberID:   sql.NullInt64{Int64: identity.Member.ID, Valid: true},
			Action:     db.ActivityAction_Update,
			EntityType: db.ActivityEntity_Space,
			EntityID:   sql.NullInt64{Int64: space.ID, Valid: true},
		})
		if err != nil {
			logger.Error("error creating activity in database", slog.Any("error", err))
			http.Error(w, "error updating space", http.StatusInternalServerError)
			return
		}

		err = tx.Commit()
		if err != nil {
			logger.Error("error committing transaction", slog.Any("error", err))
			http.Error(w, "error updating space", http.StatusInternalServerError)
			return
		}

		sess := rctx.GetSession(r.Context())
		sess.AddFlash(session.FlashMessage{
			Type:    session.FlashType_Success,
			Content: fmt.Sprintf("%s updated successfully.", spaceUpdated.Name),
		})
		err = sess.Save(r, w)
		if err != nil {
			logger.Error("failed to save session", slog.Any("error", err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		access := rctx.GetAccess(r.Context())
		http.Redirect(w, r, fmt.Sprintf("/s/%s/settings", access.Token), http.StatusSeeOther)
	}
}

func handleTokensShow(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		space := rctx.GetSpace(r.Context())
		tokens := access.AccessTokens{
			AdminToken: space.AdminToken,
			EditToken:  space.EditToken,
			ViewToken:  space.ViewToken,
		}

		scheme := "http"
		if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
			scheme = "https"
		}
		baseURL := scheme + "://" + r.Host

		err := pages.TokensShow(baseURL, tokens).Render(r.Context(), w)
		if err != nil {
			logger.Error(
				"failed to render template",
				slog.Any("error", err),
				slog.String("template", "TokensShow"),
			)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}

func handleActivityList(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := pages.ActivityList().Render(r.Context(), w)
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
