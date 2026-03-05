package api

import (
	"net/http"

	"github.com/justEstif/todo-open/internal/api/handlers"
	"github.com/justEstif/todo-open/internal/api/web"
	"github.com/justEstif/todo-open/internal/core"
)

func NewRouter(taskService core.TaskService, adapterRuntime handlers.AdapterRuntimeResponse) http.Handler {
	mux := http.NewServeMux()
	tasks := handlers.NewTaskHandler(taskService)
	adapters := handlers.NewAdapterHandler(adapterRuntime)
	assets := web.AssetsHandler()
	index := web.IndexHandler()

	mux.HandleFunc("GET /healthz", handlers.Health)
	mux.HandleFunc("GET /v1/adapters", adapters.List)
	mux.HandleFunc("POST /v1/tasks", tasks.Create)
	mux.HandleFunc("GET /v1/tasks", tasks.List)
	mux.HandleFunc("GET /v1/tasks/{id}", tasks.Get)
	mux.HandleFunc("PATCH /v1/tasks/{id}", tasks.Update)
	mux.HandleFunc("DELETE /v1/tasks/{id}", tasks.Delete)

	mux.Handle("GET /static/", http.StripPrefix("/static/", assets))
	mux.Handle("GET /", index)
	return withRequestLogging(mux)
}
