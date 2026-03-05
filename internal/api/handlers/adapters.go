package handlers

import "net/http"

// AdapterConfigResponse captures user-facing adapter configuration.
type AdapterConfigResponse struct {
	EnabledViews        []string                  `json:"enabled_views"`
	EnabledSyncAdapters []string                  `json:"enabled_sync_adapters"`
	ViewSettings        map[string]map[string]any `json:"view_settings,omitempty"`
	SyncSettings        map[string]map[string]any `json:"sync_settings,omitempty"`
}

// AdapterStatusResponse represents one adapter's runtime status.
type AdapterStatusResponse struct {
	Kind    string `json:"kind"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	Healthy bool   `json:"healthy"`
	Message string `json:"message,omitempty"`
}

// AdapterRuntimeResponse mirrors the user-facing adapter runtime payload.
type AdapterRuntimeResponse struct {
	Config AdapterConfigResponse   `json:"config"`
	Status []AdapterStatusResponse `json:"status"`
	Ready  bool                    `json:"ready"`
	Errors []string                `json:"errors,omitempty"`
}

type AdapterHandler struct {
	runtime AdapterRuntimeResponse
}

func NewAdapterHandler(runtime AdapterRuntimeResponse) *AdapterHandler {
	return &AdapterHandler{runtime: runtime}
}

func (h *AdapterHandler) List(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.runtime)
}
