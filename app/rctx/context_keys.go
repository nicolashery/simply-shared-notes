package rctx

type contextKey string

const (
	viteContextKey       contextKey = "vite"
	sessionContextKey    contextKey = "session"
	themeContextKey      contextKey = "theme"
	spaceContextKey      contextKey = "space"
	spaceStatsContextKey contextKey = "space_stats"
	accessContextKey     contextKey = "access"
	identityContextKey   contextKey = "identity"
	flashContextKey      contextKey = "flash"
	memberContextKey     contextKey = "member"
	noteContextKey       contextKey = "note"
)
