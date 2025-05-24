package handlers

import (
	"net/http"

	"github.com/nicolashery/simply-shared-notes/app/views/pages"
)

func handleMembersList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pages.MembersList().Render(r.Context(), w)
	}
}

func handleMembersEdit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pages.MembersEdit().Render(r.Context(), w)
	}
}
