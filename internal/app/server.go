package app

import (
	"net/http"
	"time"

	"github.com/ebeyene/todo-open/internal/api"
	"github.com/ebeyene/todo-open/internal/core"
	"github.com/ebeyene/todo-open/internal/store/memory"
)

func NewServer(addr string) *http.Server {
	repo := memory.NewTaskRepo()
	taskService := core.NewService(repo, time.Now, nil)

	return &http.Server{
		Addr:    addr,
		Handler: api.NewRouter(taskService),
	}
}
