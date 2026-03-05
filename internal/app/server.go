package app

import (
	"net/http"

	"github.com/ebeyene/todo-open/internal/api"
)

func NewServer(addr string) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: api.NewRouter(),
	}
}
