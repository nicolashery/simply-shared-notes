package rctx

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nicolashery/simply-shared-notes/app/db"
	"github.com/nicolashery/simply-shared-notes/app/publicid"
)

func MemberCtxMiddleware(queries *db.Queries) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			memberPublicID := chi.URLParam(r, "memberPublicID")
			if !publicid.IsValidPublicID(memberPublicID) {
				http.Error(w, "member not found", http.StatusNotFound)
				return
			}

			space := GetSpace(r.Context())
			member, err := queries.GetMemberByPublicID(r.Context(), db.GetMemberByPublicIDParams{
				SpaceID:  space.ID,
				PublicID: memberPublicID,
			})
			if err != nil {
				http.Error(w, "member not found", http.StatusNotFound)
				return
			}

			ctx := context.WithValue(r.Context(), memberContextKey, &member)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetMember(ctx context.Context) *db.Member {
	member, ok := ctx.Value(memberContextKey).(*db.Member)
	if !ok {
		panic("member not found in context, make sure to use middleware")
	}

	return member
}
