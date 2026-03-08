// Package info builds the agent-info payload for todoopen --agent-info.
package info

// AgentInfo is the top-level structure emitted by todoopen --agent-info.
type AgentInfo struct {
	Tool      ToolInfo     `json:"tool"`
	Server    ServerInfo   `json:"server"`
	Endpoints EndpointMap  `json:"endpoints"`
	Schema    SchemaInfo   `json:"schema"`
	Workflow  WorkflowInfo `json:"workflow"`
}

// ToolInfo describes the binary itself.
type ToolInfo struct {
	Name        string          `json:"name"`
	Version     string          `json:"version"`
	Description string          `json:"description"`
	Usage       string          `json:"usage"`
	Flags       map[string]Flag `json:"flags"`
}

// Flag describes a single CLI flag.
type Flag struct {
	Description string `json:"description"`
	Default     string `json:"default,omitempty"`
}

// ServerInfo describes how to reach the server.
type ServerInfo struct {
	BaseURL      string `json:"base_url"`
	HealthURL    string `json:"health_url"`
	Capabilities string `json:"capabilities_url"`
}

// EndpointMap groups the REST + SSE surface by concern.
type EndpointMap struct {
	Tasks  TaskEndpoints  `json:"tasks"`
	Agent  AgentEndpoints `json:"agent"`
	Events string         `json:"events"`
}

// TaskEndpoints are the CRUD endpoints.
type TaskEndpoints struct {
	List   string `json:"list"`
	Create string `json:"create"`
	Get    string `json:"get"`
	Update string `json:"update"`
	Delete string `json:"delete"`
	Next   string `json:"next"`
}

// AgentEndpoints are the agent-coordination endpoints.
type AgentEndpoints struct {
	Claim     string `json:"claim"`
	Heartbeat string `json:"heartbeat"`
	Complete  string `json:"complete"`
	Release   string `json:"release"`
}

// SchemaInfo documents the constrained field values.
type SchemaInfo struct {
	Statuses   map[string]string `json:"statuses"`
	Priorities map[string]string `json:"priorities"`
	Notes      string            `json:"notes"`
}

// WorkflowInfo documents the recommended agent task loop.
type WorkflowInfo struct {
	Steps           []string `json:"steps"`
	LeaseTTLSeconds int      `json:"lease_ttl_seconds"`
	IdempotencyNote string   `json:"idempotency_note"`
	HeartbeatNote   string   `json:"heartbeat_note"`
}

// Build assembles an AgentInfo for the given server base URL and version.
func Build(version, baseURL string) AgentInfo {
	return AgentInfo{
		Tool: ToolInfo{
			Name:    "todoopen",
			Version: version,
			Description: "A local-first task server with an open API. " +
				"Tasks live on disk as tasks.jsonl; any tool — CLI, agent, or script — " +
				"talks to the same HTTP endpoint.",
			Usage: "todoopen [command] [flags]",
			Flags: map[string]Flag{
				"--agent-info": {
					Description: "Print this agent-info JSON and exit. No server required.",
				},
				"--server": {
					Description: "Base URL of a running todo.open server.",
					Default:     "http://127.0.0.1:8080",
				},
			},
		},
		Server: ServerInfo{
			BaseURL:      baseURL,
			HealthURL:    baseURL + "/healthz",
			Capabilities: baseURL + "/v1/capabilities",
		},
		Endpoints: EndpointMap{
			Tasks: TaskEndpoints{
				List:   "GET    " + baseURL + "/v1/tasks",
				Create: "POST   " + baseURL + "/v1/tasks",
				Get:    "GET    " + baseURL + "/v1/tasks/{id}",
				Update: "PATCH  " + baseURL + "/v1/tasks/{id}",
				Delete: "DELETE " + baseURL + "/v1/tasks/{id}",
				Next:   "GET    " + baseURL + "/v1/tasks/next",
			},
			Agent: AgentEndpoints{
				Claim:     "POST " + baseURL + "/v1/tasks/{id}/claim",
				Heartbeat: "POST " + baseURL + "/v1/tasks/{id}/heartbeat",
				Complete:  "POST " + baseURL + "/v1/tasks/{id}/complete",
				Release:   "POST " + baseURL + "/v1/tasks/{id}/release",
			},
			Events: "GET  " + baseURL + "/v1/tasks/events  (SSE)",
		},
		Schema: SchemaInfo{
			Statuses: map[string]string{
				"pending":     "Waiting for trigger_ids dependencies to complete.",
				"open":        "Ready to be claimed.",
				"in_progress": "Claimed by an agent or user; a lease is active.",
				"done":        "Completed successfully.",
				"archived":    "Removed from active queues; retained in history.",
			},
			Priorities: map[string]string{
				"low":      "Below normal; work when nothing else is queued.",
				"normal":   "Default priority.",
				"high":     "Prefer over normal tasks.",
				"critical": "Urgent; address immediately if possible.",
			},
			Notes: "Tasks are stored as tasks.jsonl (one JSON object per line). " +
				"All timestamps are RFC3339 UTC. Extension data belongs under the 'ext' key.",
		},
		Workflow: WorkflowInfo{
			Steps: []string{
				"GET /v1/tasks/next            — fetch the highest-priority open task",
				"POST /v1/tasks/{id}/claim     — acquire a lease (X-Idempotency-Key recommended)",
				"POST /v1/tasks/{id}/heartbeat — renew the lease every ~60 s while working",
				"POST /v1/tasks/{id}/complete  — mark done when finished",
				"POST /v1/tasks/{id}/release   — release the lease without completing (on error/abort)",
			},
			LeaseTTLSeconds: 300,
			IdempotencyNote: "Pass X-Idempotency-Key on claim and other mutating requests. " +
				"Replaying the same key returns the original response without side-effects.",
			HeartbeatNote: "If no heartbeat is received within the lease TTL the task reverts " +
				"to open and becomes available for other agents to claim.",
		},
	}
}
