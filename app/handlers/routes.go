package handlers

import (
	"database/sql"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/nicolashery/simply-shared-notes/app/config"
	"github.com/nicolashery/simply-shared-notes/app/db"
	"github.com/nicolashery/simply-shared-notes/app/rctx"
)

func RegisterRoutes(r chi.Router, cfg *config.Config, logger *slog.Logger, sqlDB *sql.DB, queries *db.Queries, sessionStore *sessions.CookieStore) {
	r.Get("/", handleHome(logger))

	r.Get("/new", handleSpacesNew(cfg, logger))
	r.Post("/new", handleSpacesCreate(cfg, logger, sqlDB, queries, sessionStore))

	r.Get("/language", handleLanguageSelect(logger))

	r.Route("/s/{token}", func(r chi.Router) {
		r.Use(rctx.SpaceCtxMiddleware(queries))
		r.Use(rctx.AccessCtxMiddleware(logger))

		r.Get("/identity", handleIdentitySelect(logger, queries))
		r.Post("/identity", handleIdentitySet(logger, queries, sessionStore))
		r.Post("/identity/delete", handleIdentityClear(logger, sessionStore))

		r.Group(func(r chi.Router) {
			r.Use(rctx.IdentityCtxMiddleware(logger, sessionStore, queries))
			r.Use(rctx.FlashCtxMiddleware(logger, sessionStore))

			r.Get("/", handleSpacesShow(logger))
			r.Get("/settings", handleSpacesEdit(logger))

			r.Get("/share", handleTokensShow(logger))

			r.Get("/notes", handleNotesList(logger))

			r.Get("/members", handleMembersList(logger))
			r.Get("/members/{memberId}/edit", handleMembersEdit(logger))

			r.Get("/activity", handleActivityList(logger))
		})
	})
}
