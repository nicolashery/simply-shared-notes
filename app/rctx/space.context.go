package rctx

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nicolashery/simply-shared-notes/app/db"
	"github.com/nicolashery/simply-shared-notes/app/helpers"
)

func SpaceCtxMiddleware(queries *db.Queries) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := chi.URLParam(r, "token")
			if !helpers.IsValidAccessToken(token) {
				http.Error(w, "space not found", http.StatusNotFound)
				return
			}

			space, err := queries.GetSpaceByAccessToken(r.Context(), token)
			if err != nil {
				http.Error(w, "space not found", http.StatusNotFound)
				return
			}

			ctx := context.WithValue(r.Context(), spaceContextKey, &space)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetSpace(ctx context.Context) *db.Space {
	space, ok := ctx.Value(spaceContextKey).(*db.Space)
	if !ok {
		panic("space not found in context, make sure to use middleware")
	}

	return space
}
