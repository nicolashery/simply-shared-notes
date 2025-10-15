package rctx

import (
	"context"
	"log/slog"
	"net/http"
	"slices"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/nicolashery/simply-shared-notes/app/intl"
	"github.com/nicolashery/simply-shared-notes/app/session"
	"golang.org/x/text/language"
)

func IntlCtxMiddleware(logger *slog.Logger, i18nBundle *i18n.Bundle) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lang, preferredLang := selectLanguage(r)
			tz := selectTimezone(r)

			intl := intl.New(logger, i18nBundle, lang, preferredLang, tz)

			ctx := context.WithValue(r.Context(), intlContextKey, intl)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func selectLanguage(r *http.Request) (language.Tag, *language.Tag) {
	var preferredLang *language.Tag

	// Try session cookie first
	sess := GetSession(r.Context())
	if cookieLang, ok := sess.Values[session.LangKey].(string); ok {
		if tag, err := language.Parse(cookieLang); err == nil {
			if slices.Contains(intl.SupportedLangs, tag) {
				preferredLang = &tag
				return *preferredLang, preferredLang
			}
		}
	}

	// Try Accept-Language header
	if acceptLang := r.Header.Get("Accept-Language"); acceptLang != "" {
		if tags, _, err := language.ParseAcceptLanguage(acceptLang); err == nil {
			// Find the first tag that matches one of our supported languages
			for _, tag := range tags {
				if slices.Contains(intl.SupportedLangs, tag) {
					return tag, preferredLang
				}
			}
		}
	}

	// Fall back to default
	return intl.DefaultLang, preferredLang
}

func selectTimezone(r *http.Request) *time.Location {
	tz := intl.DefaultTimezone

	c, err := r.Cookie("tz")
	if err == nil && c.Value != "" {
		l, err := time.LoadLocation(c.Value)
		if err == nil {
			tz = l
		}
	}

	return tz
}

func GetIntl(ctx context.Context) *intl.Intl {
	intl, ok := ctx.Value(intlContextKey).(*intl.Intl)
	if !ok {
		panic("intl not found in context, make sure to use middleware")
	}

	return intl
}
