package server

import (
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func StaticDir(r chi.Router, path string, root fs.FS) {
	pattern := path + "/*"

	r.Get(pattern, func(w http.ResponseWriter, r *http.Request) {
		fs := http.StripPrefix(path, http.FileServer(http.FS(root)))
		fs.ServeHTTP(w, r)
	})
}

func StaticFile(r chi.Router, path string, root fs.FS) {
	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		fs := http.FileServer(http.FS(root))
		fs.ServeHTTP(w, r)
	})
}
