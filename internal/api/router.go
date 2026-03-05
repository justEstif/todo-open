package api

import (
	"net/http"

	"github.com/ebeyene/todo-open/internal/api/handlers"
	"github.com/ebeyene/todo-open/internal/core"
)

func NewRouter(taskService core.TaskService) http.Handler {
	mux := http.NewServeMux()
	tasks := handlers.NewTaskHandler(taskService)

	mux.HandleFunc("GET /healthz", handlers.Health)
	mux.HandleFunc("POST /v1/tasks", tasks.Create)
	mux.HandleFunc("GET /v1/tasks", tasks.List)
	mux.HandleFunc("GET /v1/tasks/{id}", tasks.Get)
	mux.HandleFunc("PATCH /v1/tasks/{id}", tasks.Update)
	mux.HandleFunc("DELETE /v1/tasks/{id}", tasks.Delete)
	return mux
}
