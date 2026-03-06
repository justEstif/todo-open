package plugin

import (
	"errors"
	"testing"
)

func TestValidateHandshakeOKView(t *testing.T) {
	t.Parallel()

	req := HandshakeRequest{
		ProtocolVersion: ProtocolVersion,
		ExpectedName:    "markdown",
		ExpectedKind:    AdapterKindView,
	}
	resp := HandshakeResponse{
		ProtocolVersion: ProtocolVersion,
		Name:            "markdown",
		Kind:            AdapterKindView,
		Capabilities:    []Capability{CapabilityRenderTasks},
		Health:          Health{State: HealthReady},
	}

	if err := ValidateHandshake(req, resp); err != nil {
		t.Fatalf("validate handshake: %v", err)
	}
}

func TestValidateHandshakeOKSync(t *testing.T) {
	t.Parallel()

	req := HandshakeRequest{
		ProtocolVersion: ProtocolVersion,
		ExpectedName:    "git",
		ExpectedKind:    AdapterKindSync,
	}
	resp := HandshakeResponse{
		ProtocolVersion: ProtocolVersion,
		Name:            "git",
		Kind:            AdapterKindSync,
		Capabilities: []Capability{
			CapabilityPull,
			CapabilityPush,
			CapabilityStatus,
		},
		Health: Health{State: HealthReady},
	}

	if err := ValidateHandshake(req, resp); err != nil {
		t.Fatalf("validate handshake: %v", err)
	}
}

func TestValidateHandshakeMismatches(t *testing.T) {
	t.Parallel()

	baseReq := HandshakeRequest{
		ProtocolVersion: ProtocolVersion,
		ExpectedName:    "markdown",
		ExpectedKind:    AdapterKindView,
	}
	baseResp := HandshakeResponse{
		ProtocolVersion: ProtocolVersion,
		Name:            "markdown",
		Kind:            AdapterKindView,
		Capabilities:    []Capability{CapabilityRenderTasks},
		Health:          Health{State: HealthReady},
	}

	tests := []struct {
		name string
		req  HandshakeRequest
		resp HandshakeResponse
		want error
	}{
		{
			name: "version mismatch",
			req:  baseReq,
			resp: func() HandshakeResponse {
				v := baseResp
				v.ProtocolVersion = "todoopen.plugin.v2"
				return v
			}(),
			want: ErrProtocolVersionMismatch,
		},
		{
			name: "name mismatch",
			req:  baseReq,
			resp: func() HandshakeResponse {
				v := baseResp
				v.Name = "other"
				return v
			}(),
			want: ErrAdapterNameMismatch,
		},
		{
			name: "kind mismatch",
			req:  baseReq,
			resp: func() HandshakeResponse {
				v := baseResp
				v.Kind = AdapterKindSync
				v.Capabilities = []Capability{CapabilityPull, CapabilityPush, CapabilityStatus}
				return v
			}(),
			want: ErrAdapterKindMismatch,
		},
		{
			name: "missing required capability",
			req:  baseReq,
			resp: func() HandshakeResponse {
				v := baseResp
				v.Capabilities = nil
				return v
			}(),
			want: ErrCapabilityMissing,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := ValidateHandshake(tc.req, tc.resp)
			if !errors.Is(err, tc.want) {
				t.Fatalf("err = %v, want %v", err, tc.want)
			}
		})
	}
}

func TestPluginErrorString(t *testing.T) {
	t.Parallel()

	err := (&PluginError{Code: ErrorCodeUnavailable, Message: "plugin offline"}).Error()
	if err != "unavailable: plugin offline" {
		t.Fatalf("error = %q", err)
	}

	err = (&PluginError{Code: ErrorCodeTimeout, Message: "request timed out", Detail: "render_tasks"}).Error()
	if err != "timeout: request timed out (render_tasks)" {
		t.Fatalf("error = %q", err)
	}
}
