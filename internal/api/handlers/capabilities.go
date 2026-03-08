package handlers

import (
	"net/http"

	"github.com/justEstif/todo-open/internal/version"
)

// CapabilitiesResponse represents the server capabilities and schema information
type CapabilitiesResponse struct {
	SchemaVersion int    `json:"schema_version"`
	ServerVersion string `json:"server_version"`
	Task          Task   `json:"task"`
	Agent         Agent  `json:"agent"`
	Events        Events `json:"events"`
	ETag          ETag   `json:"etag"`
}

// Task describes task-related capabilities
type Task struct {
	Statuses      []string `json:"statuses"`
	Priorities    []string `json:"priorities"`
	ExtNamespaces []string `json:"ext_namespaces"`
}

// Agent describes agent-specific capabilities
type Agent struct {
	Endpoints struct {
		Next      string `json:"next"`
		Claim     string `json:"claim"`
		Heartbeat string `json:"heartbeat"`
		Release   string `json:"release"`
		Complete  string `json:"complete"`
	} `json:"endpoints"`
	LeaseTTLDefaultSeconds int    `json:"lease_ttl_default_seconds"`
	IdempotencyHeader      string `json:"idempotency_header"`
	IdempotencyTTSeconds   int    `json:"idempotency_ttl_seconds"`
}

// Events describes event streaming capabilities
type Events struct {
	Endpoint string   `json:"endpoint"`
	Types    []string `json:"types"`
}

// ETag describes ETag capabilities
type ETag struct {
	HeaderRequest  string `json:"header_request"`
	HeaderResponse string `json:"header_response"`
	Format         string `json:"format"`
}

// Capabilities returns the server capabilities and schema information
func Capabilities(w http.ResponseWriter, _ *http.Request) {
	capabilities := CapabilitiesResponse{
		SchemaVersion: 1,
		ServerVersion: version.Version,
		Task: Task{
			Statuses:      []string{"pending", "open", "in_progress", "done", "archived"},
			Priorities:    []string{"low", "normal", "high", "critical"},
			ExtNamespaces: []string{"agent"},
		},
		Agent: Agent{
			Endpoints: struct {
				Next      string `json:"next"`
				Claim     string `json:"claim"`
				Heartbeat string `json:"heartbeat"`
				Release   string `json:"release"`
				Complete  string `json:"complete"`
			}{
				Next:      "GET  /v1/tasks/next",
				Claim:     "POST /v1/tasks/{id}/claim",
				Heartbeat: "POST /v1/tasks/{id}/heartbeat",
				Release:   "POST /v1/tasks/{id}/release",
				Complete:  "POST /v1/tasks/{id}/complete",
			},
			LeaseTTLDefaultSeconds: 300,
			IdempotencyHeader:      "X-Idempotency-Key",
			IdempotencyTTSeconds:   300,
		},
		Events: Events{
			Endpoint: "GET /v1/tasks/events",
			Types:    []string{"task.created", "task.updated", "task.deleted", "task.status_changed"},
		},
		ETag: ETag{
			HeaderRequest:  "If-Match",
			HeaderResponse: "ETag",
			Format:         "\"<integer version>\"",
		},
	}

	writeJSON(w, http.StatusOK, capabilities)
}
