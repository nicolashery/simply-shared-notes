package handlers

import (
	"net/http"

	"github.com/nicolashery/simply-shared-notes/views/pages"
)

func handleHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pages.Home().Render(r.Context(), w)
	}
}
