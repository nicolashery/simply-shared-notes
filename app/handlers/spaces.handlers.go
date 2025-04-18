package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/nicolashery/simply-shared-notes/app/db"
	"github.com/nicolashery/simply-shared-notes/app/helpers"
	"github.com/nicolashery/simply-shared-notes/app/views/pages"
)

func handleSpacesNew() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pages.SpacesNew().Render(r.Context(), w)
	}
}

type CreateSpaceForm struct {
	Name  string
	Email string
}

func parseCreateSpaceForm(r *http.Request, f *CreateSpaceForm) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	f.Name = strings.Trim(r.Form.Get("name"), " ")
	f.Email = strings.Trim(r.Form.Get("email"), " ")

	return nil
}

func handleSpacesCreate(logger *slog.Logger, queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var form CreateSpaceForm
		err := parseCreateSpaceForm(r, &form)
		if err != nil {
			http.Error(w, "failed to parse form", http.StatusBadRequest)
			return
		}

		tokens, err := helpers.GenerateAccessTokens()
		if err != nil {
			logger.Error("error generating access tokens", slog.Any("error", err))
			http.Error(w, "error creating space", http.StatusInternalServerError)
			return
		}

		now := time.Now().UTC()
		space, err := queries.CreateSpace(r.Context(), db.CreateSpaceParams{
			CreatedAt:  now,
			UpdatedAt:  now,
			Name:       form.Name,
			Email:      form.Email,
			AdminToken: tokens.AdminToken,
			EditToken:  tokens.EditToken,
			ViewToken:  tokens.ViewToken,
		})
		if err != nil {
			logger.Error("error inserting space into database", slog.Any("error", err))
			http.Error(w, "error creating space", http.StatusInternalServerError)
			return
		}

		spaceUrl := fmt.Sprintf("/s/%s", space.AdminToken)
		http.Redirect(w, r, spaceUrl, http.StatusSeeOther)
	}
}

func handleSpacesShow() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pages.SpacesShow().Render(r.Context(), w)
	}
}
