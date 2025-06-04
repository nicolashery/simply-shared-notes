package rctx

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/nicolashery/simply-shared-notes/app/session"
)

func FlashCtxMiddleware(logger *slog.Logger, sessionStore *sessions.CookieStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sess := GetSession(r.Context())
			rawMessages := sess.Flashes()
			messages := []session.FlashMessage{}
			for _, rawMessage := range rawMessages {
				message, ok := rawMessage.(*session.FlashMessage)
				if !ok {
					logger.Warn("could not type cast flash message", slog.Any("message", rawMessage))
					continue
				}

				messages = append(messages, *message)
			}

			err := sess.Save(r, w)
			if err != nil {
				logger.Error("failed to save session", slog.Any("error", err))
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), flashContextKey, messages)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetFlashMessages(ctx context.Context) []session.FlashMessage {
	messages, ok := ctx.Value(flashContextKey).([]session.FlashMessage)
	if !ok {
		panic("flash not found in context, make sure to use middleware")
	}

	return messages
}
