package session

import (
	"encoding/gob"

	"github.com/gorilla/sessions"
)

const (
	CookieName  = "simplysharednotes_session"
	ThemeKey    = "theme"
	IdentityKey = "identity"
	RedirectKey = "redirect_url"
)

func InitStore(secret string, isDev bool) *sessions.CookieStore {
	store := sessions.NewCookieStore([]byte(secret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 90, // 90 days
		HttpOnly: true,
		Secure:   !isDev,
	}

	gob.Register(&FlashMessage{})

	return store
}
