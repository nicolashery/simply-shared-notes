package routes

import (
	"log/slog"
	"net/http"
)

func handleHome(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("handling home")
		w.Write([]byte("Hello, world!"))
	}
}
