package rctx

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nicolashery/simply-shared-notes/app/access"
)

func AccessCtxMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := chi.URLParam(r, "token")
			if !access.IsValidAccessToken(token) {
				http.Error(w, "space not found", http.StatusNotFound)
				return
			}

			space := GetSpace(r.Context())

			role, ok := access.GetTokenRole(space, token)
			if !ok {
				logger.Error(
					"failed to get role for token",
					slog.Int64("space_id", space.ID),
					slog.String("token", token),
				)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}

			access := access.Access{
				Token: token,
				Role:  role,
			}

			ctx := context.WithValue(r.Context(), accessContextKey, &access)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetAccess(ctx context.Context) *access.Access {
	access, ok := ctx.Value(accessContextKey).(*access.Access)
	if !ok {
		panic("access not found in context, make sure to use middleware")
	}

	return access
}
