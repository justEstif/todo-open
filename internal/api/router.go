package api

import (
	"net/http"

	"github.com/ebeyene/todo-open/internal/api/handlers"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handlers.Health)
	return mux
}
