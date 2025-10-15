package handlers

import (
	"database/sql"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/nicolashery/simply-shared-notes/app/access"
	"github.com/nicolashery/simply-shared-notes/app/config"
	"github.com/nicolashery/simply-shared-notes/app/db"
	"github.com/nicolashery/simply-shared-notes/app/email"
	"github.com/nicolashery/simply-shared-notes/app/rctx"
)

func RegisterRoutes(
	r chi.Router,
	cfg *config.Config,
	logger *slog.Logger,
	sqlDB *sql.DB,
	queries *db.Queries,
	sessionStore *sessions.CookieStore,
	email *email.Email,
	i18nBundle *i18n.Bundle,
) {
	r.Use(rctx.SessionCtxMiddleware(logger, sessionStore))
	r.Use(rctx.ThemeCtxMiddleware())
	r.Use(rctx.IntlCtxMiddleware(logger, i18nBundle))

	r.Get("/", handleHome(logger))

	r.Get("/new", handleSpacesNew(cfg, logger))
	r.Post("/new", handleSpacesCreate(cfg, logger, sqlDB, queries, email))
	r.With(rctx.FlashCtxMiddleware(logger)).Get("/new/success", handleSpacesNewSuccess(logger))

	r.Get("/language", handleLanguageSelect(logger))
	r.Post("/language", handleLanguageSet(logger))
	r.Get("/theme", handleThemeSelect(logger))
	r.Post("/theme", handleThemeSet(logger))

	r.Route("/s/{token}", func(r chi.Router) {
		r.Use(rctx.SpaceCtxMiddleware(queries))
		r.Use(rctx.AccessCtxMiddleware(logger))

		r.Get("/identity", handleIdentitySelect(logger, queries))
		r.Post("/identity", handleIdentitySet(logger, queries))
		r.Post("/identity/delete", handleIdentityClear(logger))

		r.Group(func(r chi.Router) {
			r.Use(rctx.IdentityCtxMiddleware(logger, queries))
			r.Use(rctx.SpaceStatsCtxMiddleware(queries))
			r.Use(rctx.FlashCtxMiddleware(logger))

			r.Get("/", handleSpacesShow(logger))
			r.With(Authorize(access.Action_UpdateSpace)).Group(func(r chi.Router) {
				r.Get("/settings", handleSpacesEdit(logger, queries))
				r.Post("/settings", handleSpacesUpdate(logger, sqlDB, queries))
			})

			r.With(Authorize(access.Action_ViewTokens)).
				Get("/share", handleTokensShow(logger))

			r.Get("/notes", handleNotesList(logger, queries))
			r.With(Authorize(access.Action_CreateNote)).Group(func(r chi.Router) {
				r.Get("/notes/new", handleNotesNew(logger))
				r.Post("/notes/new", handleNotesCreate(logger, sqlDB, queries))
			})
			r.Route("/notes/{notePublicID}", func(r chi.Router) {
				r.Use(rctx.NoteCtxMiddleware(queries))

				r.Get("/", handleNotesShow(logger))

				r.With(Authorize(access.Action_UpdateNote)).Group(func(r chi.Router) {
					r.Get("/edit", handleNotesEdit(logger, queries))
					r.Post("/edit", handleNotesUpdate(logger, sqlDB, queries))
				})

				r.With(Authorize(access.Action_DeleteNote)).Group(func(r chi.Router) {
					r.Get("/delete", handleNotesDeleteConfirm(logger))
					r.Post("/delete", handleNotesDelete(logger, sqlDB, queries))
				})
			})

			r.Get("/members", handleMembersList(logger, queries))
			r.With(Authorize(access.Action_CreateMember)).Group(func(r chi.Router) {
				r.Get("/members/new", handleMembersNew(logger))
				r.Post("/members/new", handleMembersCreate(logger, sqlDB, queries))
			})
			r.Route("/members/{memberPublicID}", func(r chi.Router) {
				r.Use(rctx.MemberCtxMiddleware(queries))

				r.With(Authorize(access.Action_UpdateMember)).Group(func(r chi.Router) {
					r.Get("/edit", handleMembersEdit(logger, queries))
					r.Post("/edit", handleMembersUpdate(logger, sqlDB, queries))
				})

				r.With(Authorize(access.Action_DeleteMember)).Group(func(r chi.Router) {
					r.Get("/delete", handleMembersDeleteConfirm(logger))
					r.Post("/delete", handleMembersDelete(logger, sqlDB, queries))
				})
			})

			r.Get("/activity", handleActivityList(logger, queries))
			r.Get("/activity/{activityPublicID}", handleActivityShow(logger, queries))
		})
	})
}
