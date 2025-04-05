package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nicolashery/simply-shared-notes/routes"
)

func NewServer(logger *slog.Logger) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Logger, middleware.Recoverer)

	routes.RegisterRoutes(router, logger)

	return router
}

func RunServer(ctx context.Context, handler http.Handler, logger *slog.Logger, port int) error {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}

	go func() {
		<-ctx.Done()
		srv.Shutdown(context.Background())
	}()

	logger.Info(fmt.Sprintf("listening on port %d", port))

	return srv.ListenAndServe()
}
