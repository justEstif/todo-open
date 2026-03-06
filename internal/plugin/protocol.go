package plugin

import (
	"errors"
	"fmt"
)

var (
	ErrProtocolVersionMismatch = errors.New("plugin protocol version mismatch")
	ErrAdapterNameMismatch     = errors.New("plugin adapter name mismatch")
	ErrAdapterKindMismatch     = errors.New("plugin adapter kind mismatch")
	ErrCapabilityMissing       = errors.New("plugin required capability missing")
)

const ProtocolVersion = "todoopen.plugin.v1"

type AdapterKind string

const (
	AdapterKindView AdapterKind = "view"
	AdapterKindSync AdapterKind = "sync"
)

type Capability string

const (
	CapabilityRenderTasks Capability = "render_tasks"
	CapabilityPull        Capability = "pull"
	CapabilityPush        Capability = "push"
	CapabilityStatus      Capability = "status"
)

type HealthState string

const (
	HealthReady     HealthState = "ready"
	HealthDegraded  HealthState = "degraded"
	HealthUnhealthy HealthState = "unhealthy"
)

type ErrorCode string

const (
	ErrorCodeInvalidInput ErrorCode = "invalid_input"
	ErrorCodeUnavailable  ErrorCode = "unavailable"
	ErrorCodeTimeout      ErrorCode = "timeout"
	ErrorCodeNotSupported ErrorCode = "not_supported"
	ErrorCodeInternal     ErrorCode = "internal"
)

// HandshakeRequest is sent by the host when the plugin process starts.
type HandshakeRequest struct {
	ProtocolVersion string      `json:"protocol_version"`
	ExpectedName    string      `json:"expected_name"`
	ExpectedKind    AdapterKind `json:"expected_kind"`
}

// HandshakeResponse is returned by the plugin.
type HandshakeResponse struct {
	ProtocolVersion string       `json:"protocol_version"`
	Name            string       `json:"name"`
	Kind            AdapterKind  `json:"kind"`
	Capabilities    []Capability `json:"capabilities,omitempty"`
	Health          Health       `json:"health"`
}

// Health describes current plugin readiness.
type Health struct {
	State   HealthState `json:"state"`
	Message string      `json:"message,omitempty"`
}

// RequestEnvelope/ResponseEnvelope are transport-neutral request/response wrappers.
type RequestEnvelope struct {
	ID      string         `json:"id"`
	Method  string         `json:"method"`
	Payload map[string]any `json:"payload,omitempty"`
}

type ResponseEnvelope struct {
	ID      string         `json:"id"`
	Payload map[string]any `json:"payload,omitempty"`
	Error   *PluginError   `json:"error,omitempty"`
}

// PluginError is a structured plugin-side error contract.
type PluginError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Detail  string    `json:"detail,omitempty"`
}

func (e *PluginError) Error() string {
	if e == nil {
		return ""
	}
	if e.Detail == "" {
		return fmt.Sprintf("%s: %s", e.Code, e.Message)
	}
	return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Detail)
}

func ValidateHandshake(req HandshakeRequest, resp HandshakeResponse) error {
	if req.ProtocolVersion != resp.ProtocolVersion {
		return fmt.Errorf("host=%q plugin=%q: %w", req.ProtocolVersion, resp.ProtocolVersion, ErrProtocolVersionMismatch)
	}
	if req.ExpectedName != resp.Name {
		return fmt.Errorf("host=%q plugin=%q: %w", req.ExpectedName, resp.Name, ErrAdapterNameMismatch)
	}
	if req.ExpectedKind != resp.Kind {
		return fmt.Errorf("host=%q plugin=%q: %w", req.ExpectedKind, resp.Kind, ErrAdapterKindMismatch)
	}

	required := requiredCapabilities(req.ExpectedKind)
	capSet := make(map[Capability]struct{}, len(resp.Capabilities))
	for _, cap := range resp.Capabilities {
		capSet[cap] = struct{}{}
	}
	for _, cap := range required {
		if _, ok := capSet[cap]; !ok {
			return fmt.Errorf("kind=%s capability=%s: %w", req.ExpectedKind, cap, ErrCapabilityMissing)
		}
	}

	return nil
}

func requiredCapabilities(kind AdapterKind) []Capability {
	switch kind {
	case AdapterKindView:
		return []Capability{CapabilityRenderTasks}
	case AdapterKindSync:
		return []Capability{CapabilityPull, CapabilityPush, CapabilityStatus}
	default:
		return nil
	}
}
