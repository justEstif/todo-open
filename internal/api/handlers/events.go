package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/justEstif/todo-open/internal/events"
)

// EventHandler streams task events over SSE.
type EventHandler struct {
	broker *events.Broker
}

// NewEventHandler returns a new EventHandler.
func NewEventHandler(broker *events.Broker) *EventHandler {
	return &EventHandler{broker: broker}
}

// Stream handles GET /v1/tasks/events.
func (h *EventHandler) Stream(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	ch, unsub := h.broker.Subscribe(64)
	defer unsub()

	for {
		select {
		case <-r.Context().Done():
			return
		case e, ok := <-ch:
			if !ok {
				return
			}
			data, err := json.Marshal(e)
			if err != nil {
				continue
			}
			// Build SSE id from task id + version when available.
			sseID := e.At.Format("20060102T150405Z")
			if e.Task != nil {
				sseID = fmt.Sprintf("%s@%d", e.Task.ID, e.Task.Version)
			}
			fmt.Fprintf(w, "id: %s\nevent: %s\ndata: %s\n\n", sseID, e.Type, data)
			flusher.Flush()
		}
	}
}
