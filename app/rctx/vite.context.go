package rctx

import (
	"context"
	"net/http"

	"github.com/nicolashery/simply-shared-notes/app/vite"
)

func ViteCtxMiddleware(vite *vite.Vite) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), viteContextKey, vite)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetVite(ctx context.Context) *vite.Vite {
	vite, ok := ctx.Value(viteContextKey).(*vite.Vite)
	if !ok {
		panic("vite not found in context, make sure to use middleware")
	}

	return vite
}
