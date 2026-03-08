---
# go-todo-md-k4yb
title: SSE event stream
status: todo
type: feature
created_at: 2026-03-08T19:36:16Z
updated_at: 2026-03-08T19:36:16Z
parent: go-todo-md-9yck
---

Add a server-sent events stream so clients — agents, web UI, TUI — can react to task changes in real time.

## Why SSE
- Zero new dependencies: Go's `http.Flusher` is all that's needed
- curl-native: `curl -N http://localhost:8080/v1/tasks/events` just works for agents
- One-directional push is sufficient — mutations already use REST
- Non-blocking fan-out: slow subscribers drop events, don't block the broker

## No schema impact
Nothing changes in tasks.jsonl. This is purely a server transport feature.

## Design
- In-process event broker: `sync.RWMutex` + buffered `chan Event` per subscriber
- `EventEmittingService` wrapper around `core.TaskService` — delegates all calls, publishes events after successful mutations. Zero changes to existing service code.
- Event id format: `<task_id>@<version>` — clients send `Last-Event-ID` on reconnect for deduplication

## Events
- `task.created`
- `task.updated`
- `task.deleted`
- `task.status_changed` (with old/new status — useful for agents watching for `open` transitions)

## Endpoint
`GET /v1/tasks/events` — SSE stream, `text/event-stream`

## One required fix
The logging middleware's `responseRecorder` wraps `ResponseWriter` but doesn't forward `Flush()`. SSE silently fails without:
```go
func (r *responseRecorder) Flush() {
    if f, ok := r.ResponseWriter.(http.Flusher); ok { f.Flush() }
}
```
