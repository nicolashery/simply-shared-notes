package rctx

type contextKey string

const (
	viteContextKey     contextKey = "vite"
	spaceContextKey    contextKey = "space"
	accessContextKey   contextKey = "access"
	identityContextKey contextKey = "identity"
	flashContextKey    contextKey = "flash"
)
