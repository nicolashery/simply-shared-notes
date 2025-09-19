package rctx

import (
	"context"
	"net/http"

	"github.com/nicolashery/simply-shared-notes/app/db"
)

func SpaceStatsCtxMiddleware(queries *db.Queries) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			space := GetSpace(r.Context())
			stats, err := queries.GetSpaceStats(r.Context(), space.ID)
			if err != nil {
				http.Error(w, "space not found", http.StatusNotFound)
				return
			}

			ctx := context.WithValue(r.Context(), spaceStatsContextKey, &stats)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetSpaceStats(ctx context.Context) *db.GetSpaceStatsRow {
	stats, ok := ctx.Value(spaceStatsContextKey).(*db.GetSpaceStatsRow)
	if !ok {
		panic("space stats not found in context, make sure to use middleware")
	}

	return stats
}
