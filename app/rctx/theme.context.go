package rctx

import (
	"context"
	"net/http"

	"github.com/nicolashery/simply-shared-notes/app/session"
)

func ThemeCtxMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sess := GetSession(r.Context())
			theme, ok := sess.Values[session.ThemeKey].(string)
			if !ok {
				theme = ""
			}

			ctx := context.WithValue(r.Context(), themeContextKey, theme)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetTheme(ctx context.Context) string {
	theme, ok := ctx.Value(themeContextKey).(string)
	if !ok {
		return ""
	}
	return theme
}
