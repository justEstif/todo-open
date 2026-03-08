package api

import (
	"net/http"

	"github.com/justEstif/todo-open/internal/adapters"
	"github.com/justEstif/todo-open/internal/api/handlers"
	"github.com/justEstif/todo-open/internal/api/middleware"
	"github.com/justEstif/todo-open/internal/api/web"
	"github.com/justEstif/todo-open/internal/core"
	"github.com/justEstif/todo-open/internal/events"
)

func NewRouter(taskService core.TaskService, adapterRuntime adapters.Runtime, broker *events.Broker) http.Handler {
	return NewRouterWithIdempotency(taskService, adapterRuntime, broker, middleware.NewIdempotencyStore(nil))
}

func NewRouterWithIdempotency(taskService core.TaskService, adapterRuntime adapters.Runtime, broker *events.Broker, idem *middleware.IdempotencyStore) http.Handler {
	mux := http.NewServeMux()
	tasks := handlers.NewTaskHandler(taskService)
	adapterH := handlers.NewAdapterHandler(adapterRuntime)
	eventsH := handlers.NewEventHandler(broker)
	assets := web.AssetsHandler()
	index := web.IndexHandler()

	mux.HandleFunc("GET /healthz", handlers.Health)
	mux.HandleFunc("GET /v1/capabilities", handlers.Capabilities)
	mux.HandleFunc("GET /v1/adapters", adapterH.List)
	mux.HandleFunc("GET /v1/tasks/events", eventsH.Stream)
	mux.HandleFunc("GET /v1/tasks/next", tasks.NextTask)
	mux.HandleFunc("POST /v1/tasks", idem.Middleware(http.HandlerFunc(tasks.Create)).ServeHTTP)
	mux.HandleFunc("GET /v1/tasks", tasks.List)
	mux.HandleFunc("GET /v1/tasks/{id}", tasks.Get)
	mux.HandleFunc("PUT /v1/tasks/{id}", tasks.Upsert)
	mux.HandleFunc("PATCH /v1/tasks/{id}", tasks.Update)
	mux.HandleFunc("DELETE /v1/tasks/{id}", tasks.Delete)
	mux.HandleFunc("POST /v1/tasks/{id}/complete", idem.Middleware(http.HandlerFunc(tasks.Complete)).ServeHTTP)
	mux.HandleFunc("POST /v1/tasks/{id}/claim", idem.Middleware(http.HandlerFunc(tasks.Claim)).ServeHTTP)
	mux.HandleFunc("POST /v1/tasks/{id}/heartbeat", idem.Middleware(http.HandlerFunc(tasks.Heartbeat)).ServeHTTP)
	mux.HandleFunc("POST /v1/tasks/{id}/release", idem.Middleware(http.HandlerFunc(tasks.Release)).ServeHTTP)

	mux.Handle("GET /static/", http.StripPrefix("/static/", assets))
	mux.Handle("GET /", index)
	return withRequestLogging(mux)
}
