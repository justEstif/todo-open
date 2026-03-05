package handlers

import (
	"net/http"

	"github.com/justEstif/todo-open/internal/adapters"
)

type AdapterHandler struct {
	runtime adapters.Runtime
}

func NewAdapterHandler(runtime adapters.Runtime) *AdapterHandler {
	return &AdapterHandler{runtime: runtime}
}

func (h *AdapterHandler) List(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.runtime)
}
