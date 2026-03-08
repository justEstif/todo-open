package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/justEstif/todo-open/internal/plugin"
	"github.com/justEstif/todo-open/internal/sync/git"
)

func main() {
	// Write handshake response immediately
	handshake := plugin.HandshakeResponse{
		ProtocolVersion: plugin.ProtocolVersion,
		Name:            "git",
		Kind:            plugin.AdapterKindSync,
		Capabilities: []plugin.Capability{
			plugin.CapabilityPull,
			plugin.CapabilityPush,
			plugin.CapabilityStatus,
		},
		Health: plugin.Health{
			State:   plugin.HealthReady,
			Message: "",
		},
	}

	if err := writeJSONLine(handshake); err != nil {
		log.Fatalf("failed to write handshake: %v", err)
	}

	// Start processing requests
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var req plugin.RequestEnvelope
		if err := json.Unmarshal(scanner.Bytes(), &req); err != nil {
			writeJSONLine(plugin.ResponseEnvelope{
				ID: "",
				Error: &plugin.PluginError{
					Code:    plugin.ErrorCodeInternal,
					Message: fmt.Sprintf("failed to decode request: %v", err),
				},
			})
			continue
		}

		response := handleRequest(req)
		if err := writeJSONLine(response); err != nil {
			log.Fatalf("failed to write response: %v", err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("stdin error: %v", err)
	}
}

func handleRequest(req plugin.RequestEnvelope) plugin.ResponseEnvelope {
	switch req.Method {
	case "push":
		return handlePush(req)
	case "pull":
		return handlePull(req)
	case "status":
		return handleStatus(req)
	default:
		return plugin.ResponseEnvelope{
			ID: req.ID,
			Error: &plugin.PluginError{
				Code:    plugin.ErrorCodeNotSupported,
				Message: fmt.Sprintf("unknown method: %s", req.Method),
			},
		}
	}
}

func handlePush(req plugin.RequestEnvelope) plugin.ResponseEnvelope {
	workspaceRoot, ok := req.Payload["workspace_root"].(string)
	if !ok {
		return plugin.ResponseEnvelope{
			ID: req.ID,
			Error: &plugin.PluginError{
				Code:    plugin.ErrorCodeInvalidInput,
				Message: "workspace_root is required and must be a string",
			},
		}
	}

	configMap, ok := req.Payload["config"].(map[string]interface{})
	if !ok {
		configMap = make(map[string]interface{})
	}

	config := git.Config{
		Remote: getConfigString(configMap, "remote", "origin"),
		Branch: getConfigString(configMap, "branch", ""),
	}

	adapter := git.NewGitAdapter(workspaceRoot, config)
	result, err := adapter.Push()
	if err != nil {
		return plugin.ResponseEnvelope{
			ID:    req.ID,
			Error: toPluginError(err),
		}
	}

	return plugin.ResponseEnvelope{
		ID:      req.ID,
		Payload: result,
	}
}

func handlePull(req plugin.RequestEnvelope) plugin.ResponseEnvelope {
	workspaceRoot, ok := req.Payload["workspace_root"].(string)
	if !ok {
		return plugin.ResponseEnvelope{
			ID: req.ID,
			Error: &plugin.PluginError{
				Code:    plugin.ErrorCodeInvalidInput,
				Message: "workspace_root is required and must be a string",
			},
		}
	}

	configMap, ok := req.Payload["config"].(map[string]interface{})
	if !ok {
		configMap = make(map[string]interface{})
	}

	config := git.Config{
		Remote: getConfigString(configMap, "remote", "origin"),
		Branch: getConfigString(configMap, "branch", ""),
	}

	adapter := git.NewGitAdapter(workspaceRoot, config)
	result, err := adapter.Pull()
	if err != nil {
		return plugin.ResponseEnvelope{
			ID:    req.ID,
			Error: toPluginError(err),
		}
	}

	return plugin.ResponseEnvelope{
		ID:      req.ID,
		Payload: result,
	}
}

func handleStatus(req plugin.RequestEnvelope) plugin.ResponseEnvelope {
	workspaceRoot, ok := req.Payload["workspace_root"].(string)
	if !ok {
		return plugin.ResponseEnvelope{
			ID: req.ID,
			Error: &plugin.PluginError{
				Code:    plugin.ErrorCodeInvalidInput,
				Message: "workspace_root is required and must be a string",
			},
		}
	}

	adapter := git.NewGitAdapter(workspaceRoot, git.Config{})
	result, err := adapter.Status()
	if err != nil {
		return plugin.ResponseEnvelope{
			ID:    req.ID,
			Error: toPluginError(err),
		}
	}

	return plugin.ResponseEnvelope{
		ID:      req.ID,
		Payload: result,
	}
}

func getConfigString(config map[string]interface{}, key, defaultValue string) string {
	if val, ok := config[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

func toPluginError(err error) *plugin.PluginError {
	if pluginErr, ok := err.(*git.PluginError); ok {
		return &plugin.PluginError{
			Code:    plugin.ErrorCode(pluginErr.Code),
			Message: pluginErr.Message,
		}
	}
	return &plugin.PluginError{
		Code:    plugin.ErrorCodeInternal,
		Message: err.Error(),
	}
}

func writeJSONLine(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}
	data = append(data, '\n')
	_, err = os.Stdout.Write(data)
	return err
}
