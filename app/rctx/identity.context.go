package rctx

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/nicolashery/simply-shared-notes/app/db"
	"github.com/nicolashery/simply-shared-notes/app/identity"
	"github.com/nicolashery/simply-shared-notes/app/session"
)

func IdentityCtxMiddleware(logger *slog.Logger, queries *db.Queries) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			access := GetAccess(r.Context())
			if access.IsView() {
				identity := identity.Anonymous()
				ctx := context.WithValue(r.Context(), identityContextKey, &identity)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			sess := GetSession(r.Context())
			memberID, ok := sess.Values[session.IdentityKey].(int64)
			if !ok {
				delete(sess.Values, session.IdentityKey)
				err := sess.Save(r, w)
				if err != nil {
					logger.Error("failed to save session", slog.Any("error", err))
					http.Error(w, "internal server error", http.StatusInternalServerError)
					return
				}

				http.Redirect(w, r, fmt.Sprintf("/s/%s/identity", access.Token), http.StatusSeeOther)
				return
			}

			space := GetSpace(r.Context())

			member, err := queries.GetMemberByID(r.Context(), memberID)
			if err != nil || member.SpaceID != space.ID {
				delete(sess.Values, session.IdentityKey)
				err := sess.Save(r, w)
				if err != nil {
					logger.Error("failed to save session", slog.Any("error", err))
					http.Error(w, "internal server error", http.StatusInternalServerError)
					return
				}

				http.Redirect(w, r, fmt.Sprintf("/s/%s/identity", access.Token), http.StatusSeeOther)
				return
			}

			identity := identity.Member(&member)

			ctx := context.WithValue(r.Context(), identityContextKey, &identity)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetIdentity(ctx context.Context) *identity.Identity {
	identity, ok := ctx.Value(identityContextKey).(*identity.Identity)
	if !ok {
		panic("identity not found in context, make sure to use middleware")
	}

	return identity
}
