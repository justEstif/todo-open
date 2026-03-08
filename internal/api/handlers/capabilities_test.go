package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/justEstif/todo-open/internal/version"
)

func TestCapabilities(t *testing.T) {
	req := httptest.NewRequest("GET", "/v1/capabilities", nil)
	w := httptest.NewRecorder()

	Capabilities(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}

	var capabilities map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&capabilities); err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}

	schemaVersion, ok := capabilities["schema_version"].(float64)
	if !ok {
		t.Fatal("schema_version should be a number")
	}
	if schemaVersion != 1 {
		t.Errorf("expected schema_version 1, got %v", schemaVersion)
	}

	serverVersion, ok := capabilities["server_version"].(string)
	if !ok {
		t.Fatal("server_version should be a string")
	}
	if serverVersion != version.Version {
		t.Errorf("expected server_version %s, got %s", version.Version, serverVersion)
	}

	agent, ok := capabilities["agent"].(map[string]any)
	if !ok {
		t.Fatal("agent should be an object")
	}

	endpoints, ok := agent["endpoints"].(map[string]any)
	if !ok {
		t.Fatal("agent.endpoints should be an object")
	}

	nextEndpoint, ok := endpoints["next"].(string)
	if !ok {
		t.Fatal("agent.endpoints.next should be a string")
	}
	if nextEndpoint != "GET  /v1/tasks/next" {
		t.Errorf("expected agent.endpoints.next 'GET  /v1/tasks/next', got %s", nextEndpoint)
	}
}