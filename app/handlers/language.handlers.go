package handlers

import (
	"net/http"

	"github.com/nicolashery/simply-shared-notes/app/views/pages"
)

func handleLanguageSelect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pages.LanguageSelect().Render(r.Context(), w)
	}
}
