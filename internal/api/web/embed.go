package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed static/*
var staticFS embed.FS

func AssetsHandler() http.Handler {
	sub, err := fs.Sub(staticFS, "static")
	if err != nil {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "static assets unavailable", http.StatusInternalServerError)
		})
	}
	return http.FileServer(http.FS(sub))
}

func IndexHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, staticFS, "static/index.html")
	})
}
