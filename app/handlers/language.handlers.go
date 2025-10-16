package handlers

import (
	"log/slog"
	"net/http"
	"slices"
	"strings"

	"github.com/nicolashery/simply-shared-notes/app/intl"
	"github.com/nicolashery/simply-shared-notes/app/rctx"
	"github.com/nicolashery/simply-shared-notes/app/session"
	"github.com/nicolashery/simply-shared-notes/app/views/pages"
	"golang.org/x/text/language"
)

func handleLanguageSelect(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		redirect := redirectFromReferer(r)

		err := pages.LanguageSelect(redirect).Render(r.Context(), w)
		if err != nil {
			logger.Error(
				"failed to render template",
				slog.Any("error", err),
				slog.String("template", "LanguageSelect"),
			)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
	}
}

type SelectLanguageForm struct {
	Language string
	Redirect string
}

func parseSelectLanguageForm(r *http.Request, f *SelectLanguageForm) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	f.Language = strings.Trim(r.Form.Get("language"), " ")
	f.Redirect = strings.Trim(r.Form.Get("redirect"), " ")

	return nil
}

func handleLanguageSet(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var form SelectLanguageForm
		err := parseSelectLanguageForm(r, &form)
		if err != nil {
			http.Error(w, "failed to parse form", http.StatusBadRequest)
			return
		}

		sess := rctx.GetSession(r.Context())

		if form.Language == "" {
			delete(sess.Values, session.LangKey)
		} else {
			tag, err := language.Parse(form.Language)
			if err != nil {
				http.Error(w, "invalid language", http.StatusBadRequest)
				return
			}

			if !slices.Contains(intl.SupportedLangs, tag) {
				http.Error(w, "unsupported language", http.StatusBadRequest)
				return
			}

			sess.Values[session.LangKey] = form.Language
		}

		err = sess.Save(r, w)
		if err != nil {
			logger.Error("failed to save session", slog.Any("error", err))
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		redirect := safeRedirect(form.Redirect)
		if redirect == "" || redirect == "/theme" {
			redirect = "/"
		}

		http.Redirect(w, r, redirect, http.StatusSeeOther)
	}
}
