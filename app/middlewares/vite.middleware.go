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
	viteFragment, ok := ctx.Value(viteContextKey).(*vite.Fragment)
	if !ok {
		panic("vite.Fragment not found in context, make sure to use middleware")
	}

	return viteFragment
}
