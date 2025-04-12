package middlewares

import (
	"context"
	"net/http"

	"github.com/olivere/vite"
)

func ViteCtx(viteFragment *vite.Fragment) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), viteContextKey, viteFragment)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetVite(ctx context.Context) *vite.Fragment {
	if viteFragment, ok := ctx.Value(viteContextKey).(*vite.Fragment); ok {
		return viteFragment
	}
	return &vite.Fragment{}
}
