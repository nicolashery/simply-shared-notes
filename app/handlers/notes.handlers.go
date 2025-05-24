package handlers

import (
	"net/http"

	"github.com/nicolashery/simply-shared-notes/app/views/pages"
)

func handleNotesList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pages.NotesList().Render(r.Context(), w)
	}
}
