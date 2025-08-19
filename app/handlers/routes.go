package handlers

import (
	"database/sql"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/nicolashery/simply-shared-notes/app/access"
	"github.com/nicolashery/simply-shared-notes/app/config"
	"github.com/nicolashery/simply-shared-notes/app/db"
	"github.com/nicolashery/simply-shared-notes/app/rctx"
)

func RegisterRoutes(r chi.Router, cfg *config.Config, logger *slog.Logger, sqlDB *sql.DB, queries *db.Queries, sessionStore *sessions.CookieStore) {
	r.Use(rctx.SessionCtxMiddleware(logger, sessionStore))

	r.Get("/", handleHome(logger))

	r.Get("/new", handleSpacesNew(cfg, logger))
	r.Post("/new", handleSpacesCreate(cfg, logger, sqlDB, queries))

	r.Get("/language", handleLanguageSelect(logger))

	r.Route("/s/{token}", func(r chi.Router) {
		r.Use(rctx.SpaceCtxMiddleware(queries))
		r.Use(rctx.AccessCtxMiddleware(logger))

		r.Get("/identity", handleIdentitySelect(logger, queries))
		r.Post("/identity", handleIdentitySet(logger, queries))
		r.Post("/identity/delete", handleIdentityClear(logger))

		r.Group(func(r chi.Router) {
			r.Use(rctx.IdentityCtxMiddleware(logger, queries))
			r.Use(rctx.FlashCtxMiddleware(logger))

			r.Get("/", handleSpacesShow(logger))
			r.With(Authorize(access.Action_UpdateSpace)).Group(func(r chi.Router) {
				r.Get("/settings", handleSpacesEdit(logger))
				r.Post("/settings", handleSpacesUpdate(logger, queries))
			})

			r.With(Authorize(access.Action_ViewTokens)).
				Get("/share", handleTokensShow(logger))

			r.Get("/notes", handleNotesList(logger))

			r.Get("/members", handleMembersList(logger, queries))
			r.With(Authorize(access.Action_CreateMember)).
				Get("/members/new", handleMembersNew(logger))
			r.With(Authorize(access.Action_CreateMember)).
				Post("/members", handleMembersCreate(logger, queries))
			r.With(Authorize(access.Action_UpdateMember)).
				Get("/members/{memberId}/edit", handleMembersEdit(logger))

			r.Get("/activity", handleActivityList(logger))
		})
	})
}
