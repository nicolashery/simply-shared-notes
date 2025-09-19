package rctx

type contextKey string

const (
	viteContextKey       contextKey = "vite"
	sessionContextKey    contextKey = "session"
	spaceContextKey      contextKey = "space"
	spaceStatsContextKey contextKey = "space_stats"
	accessContextKey     contextKey = "access"
	identityContextKey   contextKey = "identity"
	flashContextKey      contextKey = "flash"
	memberContextKey     contextKey = "member"
	noteContextKey       contextKey = "note"
)
