package rctx

import (
	"context"
	"github.com/gorilla/sessions"
	"github.com/nicolashery/simply-shared-notes/app/session"
	"log/slog"
	"net/http"
)

func SessionCtxMiddleware(logger *slog.Logger, sessionStore *sessions.CookieStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sess, err := sessionStore.Get(r, session.CookieName)
			if err != nil {
				logger.Warn("failed to decode session", slog.Any("error", err))
			}

			ctx := context.WithValue(r.Context(), sessionContextKey, sess)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetSession(ctx context.Context) *sessions.Session {
	sess, ok := ctx.Value(sessionContextKey).(*sessions.Session)
	if !ok {
		panic("session not found in context, make sure to use middleware")
	}

	return sess
}
