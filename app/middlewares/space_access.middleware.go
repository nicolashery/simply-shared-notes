package middlewares

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nicolashery/simply-shared-notes/app/db"
	"github.com/nicolashery/simply-shared-notes/app/helpers"
)

func SpaceAccessCtx(logger *slog.Logger, queries *db.Queries) func(http.Handler) http.Handler {
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

			role, ok := helpers.GetTokenRole(&space, token)
			if !ok {
				logger.Error(
					"failed to get role for token",
					slog.Int64("space_id", space.ID),
					slog.String("token", token),
				)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}

			spaceAccess := helpers.SpaceAccess{
				Space: &space,
				Token: token,
				Role:  role,
			}

			ctx := context.WithValue(r.Context(), spaceAccessContextKey, &spaceAccess)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetSpaceAccess(ctx context.Context) *helpers.SpaceAccess {
	spaceAccess, ok := ctx.Value(spaceAccessContextKey).(*helpers.SpaceAccess)
	if !ok {
		panic("SpaceAccess not found in context, make sure to use middleware")
	}

	return spaceAccess
}
