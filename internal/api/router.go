package api

import (
	"net/http"

	"github.com/justEstif/todo-open/internal/adapters"
	"github.com/justEstif/todo-open/internal/api/handlers"
	"github.com/justEstif/todo-open/internal/api/web"
	"github.com/justEstif/todo-open/internal/core"
	"github.com/justEstif/todo-open/internal/events"
)

func NewRouter(taskService core.TaskService, adapterRuntime adapters.Runtime, broker *events.Broker) http.Handler {
	mux := http.NewServeMux()
	tasks := handlers.NewTaskHandler(taskService)
	adapters := handlers.NewAdapterHandler(adapterRuntime)
	eventsH := handlers.NewEventHandler(broker)
	assets := web.AssetsHandler()
	index := web.IndexHandler()

	mux.HandleFunc("GET /healthz", handlers.Health)
	mux.HandleFunc("GET /v1/adapters", adapters.List)
	mux.HandleFunc("GET /v1/tasks/events", eventsH.Stream)
	mux.HandleFunc("GET /v1/tasks/next", tasks.NextTask)
	mux.HandleFunc("POST /v1/tasks", tasks.Create)
	mux.HandleFunc("GET /v1/tasks", tasks.List)
	mux.HandleFunc("GET /v1/tasks/{id}", tasks.Get)
	mux.HandleFunc("PATCH /v1/tasks/{id}", tasks.Update)
	mux.HandleFunc("DELETE /v1/tasks/{id}", tasks.Delete)
	mux.HandleFunc("POST /v1/tasks/{id}/complete", tasks.Complete)
	mux.HandleFunc("POST /v1/tasks/{id}/claim", tasks.Claim)
	mux.HandleFunc("POST /v1/tasks/{id}/heartbeat", tasks.Heartbeat)
	mux.HandleFunc("POST /v1/tasks/{id}/release", tasks.Release)

	mux.Handle("GET /static/", http.StripPrefix("/static/", assets))
	mux.Handle("GET /", index)
	return withRequestLogging(mux)
}
