package adapters

// Config captures user-facing adapter configuration.
type Config struct {
	EnabledViews        []string                  `json:"enabled_views"`
	EnabledSyncAdapters []string                  `json:"enabled_sync_adapters"`
	ViewSettings        map[string]map[string]any `json:"view_settings,omitempty"`
	SyncSettings        map[string]map[string]any `json:"sync_settings,omitempty"`
}

// Status represents runtime visibility for one adapter.
type Status struct {
	Kind    string `json:"kind"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	Healthy bool   `json:"healthy"`
	Message string `json:"message,omitempty"`
}

// Runtime represents validated adapter configuration and startup health.
type Runtime struct {
	Config Config   `json:"config"`
	Status []Status `json:"status"`
	Ready  bool     `json:"ready"`
	Errors []string `json:"errors,omitempty"`
}
