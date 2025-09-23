package handlers

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/nicolashery/simply-shared-notes/app/rctx"
	"github.com/nicolashery/simply-shared-notes/app/session"
	"github.com/nicolashery/simply-shared-notes/app/views/pages"
)

func handleThemeSelect(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		redirect := redirectFromReferer(r)

		err := pages.ThemeSelect(redirect).Render(r.Context(), w)
		if err != nil {
			logger.Error(
				"failed to render template",
				slog.Any("error", err),
				slog.String("template", "ThemeSelect"),
			)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}

type SelectThemeForm struct {
	Theme    string
	Redirect string
}

func parseSelectThemeForm(r *http.Request, f *SelectThemeForm) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	f.Theme = strings.Trim(r.Form.Get("theme"), " ")
	f.Redirect = strings.Trim(r.Form.Get("redirect"), " ")

	return nil
}

func handleThemeSet(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var form SelectThemeForm
		err := parseSelectThemeForm(r, &form)
		if err != nil {
			http.Error(w, "failed to parse form", http.StatusBadRequest)
			return
		}

		validThemes := map[string]bool{
			"light": true,
			"dark":  true,
			"":      true, // system default
		}

		if !validThemes[form.Theme] {
			http.Error(w, "invalid theme", http.StatusBadRequest)
			return
		}

		sess := rctx.GetSession(r.Context())

		if form.Theme == "" {
			delete(sess.Values, session.ThemeKey)
		} else {
			sess.Values[session.ThemeKey] = form.Theme
		}

		err = sess.Save(r, w)
		if err != nil {
			logger.Error("failed to save session", slog.Any("error", err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		redirect := safeRedirect(form.Redirect)
		if redirect == "" {
			redirect = "/"
		}

		http.Redirect(w, r, redirect, http.StatusSeeOther)
	}
}
