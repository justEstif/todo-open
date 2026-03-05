package app

import (
	"net/http"
	"time"

	"github.com/justEstif/todo-open/internal/api"
	"github.com/justEstif/todo-open/internal/core"
	"github.com/justEstif/todo-open/internal/store/memory"
)

func NewServer(addr string) *http.Server {
	repo := memory.NewTaskRepo()
	taskService := core.NewService(repo, time.Now, nil)

	return &http.Server{
		Addr:    addr,
		Handler: api.NewRouter(taskService),
	}
}
