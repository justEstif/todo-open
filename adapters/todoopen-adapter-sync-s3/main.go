// todoopen-adapter-sync-s3 syncs a todo-open workspace to/from an S3-compatible object store.
//
// Works with AWS S3, Cloudflare R2, Backblaze B2, MinIO, and any S3-compatible provider.
//
// Required config (in .todoopen/config.toml):
//
//	[adapters.s3]
//	  bin  = "todoopen-adapter-sync-s3"
//	  kind = "sync"
//
//	[adapters.s3.config]
//	  bucket     = "my-tasks"
//	  key        = "tasks.jsonl"                                        # optional, default: tasks.jsonl
//	  endpoint   = "https://${R2_ACCOUNT}.r2.cloudflarestorage.com"    # optional, AWS default if omitted
//	  region     = "auto"                                               # optional, default: us-east-1
//	  access_key = "${S3_ACCESS_KEY}"
//	  secret_key = "${S3_SECRET_KEY}"
//
// Credentials are resolved from the config fields (with ${VAR} expansion).
// If access_key/secret_key are empty, the standard AWS credential chain is used.
//
// Implements the todoopen.plugin.v1 protocol over stdin/stdout.
package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
)

// --- plugin protocol types (inlined; no shared internal dependency) ---

const protocolVersion = "todoopen.plugin.v1"

type handshakeResponse struct {
	ProtocolVersion string   `json:"protocol_version"`
	Name            string   `json:"name"`
	Kind            string   `json:"kind"`
	Capabilities    []string `json:"capabilities"`
	Health          health   `json:"health"`
}

type health struct {
	State   string `json:"state"`
	Message string `json:"message,omitempty"`
}

type requestEnvelope struct {
	ID      string         `json:"id"`
	Method  string         `json:"method"`
	Payload map[string]any `json:"payload,omitempty"`
}

type responseEnvelope struct {
	ID      string         `json:"id"`
	Payload map[string]any `json:"payload,omitempty"`
	Error   *pluginError   `json:"error,omitempty"`
}

type pluginError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

// --- s3 adapter config ---

type adapterConfig struct {
	Bucket    string
	Key       string
	Endpoint  string
	Region    string
	AccessKey string
	SecretKey string
}

// expandEnv replaces ${VAR} and $VAR patterns using os.Getenv.
func expandEnv(s string) string {
	return os.ExpandEnv(s)
}

func configFromPayload(payload map[string]any) (adapterConfig, error) {
	cfg := adapterConfig{
		Key:    "tasks.jsonl",
		Region: "us-east-1",
	}
	raw, ok := payload["config"]
	if !ok {
		return cfg, fmt.Errorf("config section is required")
	}
	m, ok := raw.(map[string]any)
	if !ok {
		return cfg, fmt.Errorf("config must be an object")
	}

	if v, ok := m["bucket"].(string); ok {
		cfg.Bucket = expandEnv(v)
	}
	if v, ok := m["key"].(string); ok && v != "" {
		cfg.Key = expandEnv(v)
	}
	if v, ok := m["endpoint"].(string); ok {
		cfg.Endpoint = expandEnv(v)
	}
	if v, ok := m["region"].(string); ok && v != "" {
		cfg.Region = expandEnv(v)
	}
	if v, ok := m["access_key"].(string); ok {
		cfg.AccessKey = expandEnv(v)
	}
	if v, ok := m["secret_key"].(string); ok {
		cfg.SecretKey = expandEnv(v)
	}

	if cfg.Bucket == "" {
		return cfg, fmt.Errorf("config.bucket is required")
	}
	return cfg, nil
}

func workspaceRoot(payload map[string]any) (string, error) {
	v, ok := payload["workspace_root"].(string)
	if !ok || v == "" {
		return "", fmt.Errorf("workspace_root is required")
	}
	return v, nil
}

// --- s3 client construction ---

func newS3Client(ctx context.Context, cfg adapterConfig) (*s3.Client, error) {
	var optFns []func(*awsconfig.LoadOptions) error

	optFns = append(optFns, awsconfig.WithRegion(cfg.Region))

	if cfg.AccessKey != "" || cfg.SecretKey != "" {
		optFns = append(optFns, awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, ""),
		))
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, optFns...)
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}

	var clientOptFns []func(*s3.Options)
	if cfg.Endpoint != "" {
		clientOptFns = append(clientOptFns,
			s3.WithEndpointResolverV2(&staticEndpointResolver{endpoint: cfg.Endpoint}),
			func(o *s3.Options) {
				// Required for path-style addressing used by MinIO / non-AWS providers.
				o.UsePathStyle = true
			},
		)
	}

	return s3.NewFromConfig(awsCfg, clientOptFns...), nil
}

// staticEndpointResolver satisfies s3.EndpointResolverV2 for custom endpoints.
type staticEndpointResolver struct {
	endpoint string
}

func (r *staticEndpointResolver) ResolveEndpoint(_ context.Context, _ s3.EndpointParameters) (
	smithyendpoints.Endpoint, error,
) {
	u, err := url.Parse(r.endpoint)
	if err != nil {
		return smithyendpoints.Endpoint{}, fmt.Errorf("parse endpoint URL: %w", err)
	}
	return smithyendpoints.Endpoint{URI: *u}, nil
}

// --- handlers ---

func handlePush(ctx context.Context, payload map[string]any) (map[string]any, *pluginError) {
	root, err := workspaceRoot(payload)
	if err != nil {
		return nil, &pluginError{Code: "invalid_input", Message: err.Error()}
	}
	cfg, err := configFromPayload(payload)
	if err != nil {
		return nil, &pluginError{Code: "invalid_input", Message: err.Error()}
	}

	tasksPath := filepath.Join(root, "tasks.jsonl")
	data, err := os.ReadFile(tasksPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return map[string]any{"message": "nothing to push: tasks.jsonl does not exist"}, nil
		}
		return nil, &pluginError{Code: "internal", Message: fmt.Sprintf("read tasks.jsonl: %s", err)}
	}

	client, err := newS3Client(ctx, cfg)
	if err != nil {
		return nil, &pluginError{Code: "internal", Message: err.Error()}
	}

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(cfg.Bucket),
		Key:         aws.String(cfg.Key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String("application/x-ndjson"),
	})
	if err != nil {
		return nil, &pluginError{Code: "internal", Message: fmt.Sprintf("s3 put: %s", err)}
	}

	return map[string]any{"message": "pushed", "bytes": len(data)}, nil
}

func handlePull(ctx context.Context, payload map[string]any) (map[string]any, *pluginError) {
	root, err := workspaceRoot(payload)
	if err != nil {
		return nil, &pluginError{Code: "invalid_input", Message: err.Error()}
	}
	cfg, err := configFromPayload(payload)
	if err != nil {
		return nil, &pluginError{Code: "invalid_input", Message: err.Error()}
	}

	client, err := newS3Client(ctx, cfg)
	if err != nil {
		return nil, &pluginError{Code: "internal", Message: err.Error()}
	}

	resp, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(cfg.Bucket),
		Key:    aws.String(cfg.Key),
	})
	if err != nil {
		var nsk *types.NoSuchKey
		if errors.As(err, &nsk) {
			return map[string]any{"message": "nothing to pull: object does not exist yet"}, nil
		}
		return nil, &pluginError{Code: "internal", Message: fmt.Sprintf("s3 get: %s", err)}
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &pluginError{Code: "internal", Message: fmt.Sprintf("read s3 body: %s", err)}
	}

	tasksPath := filepath.Join(root, "tasks.jsonl")
	if err := os.WriteFile(tasksPath, data, 0o644); err != nil {
		return nil, &pluginError{Code: "internal", Message: fmt.Sprintf("write tasks.jsonl: %s", err)}
	}

	return map[string]any{"message": "pulled", "bytes": len(data)}, nil
}

func handleStatus(ctx context.Context, payload map[string]any) (map[string]any, *pluginError) {
	root, err := workspaceRoot(payload)
	if err != nil {
		return nil, &pluginError{Code: "invalid_input", Message: err.Error()}
	}
	cfg, err := configFromPayload(payload)
	if err != nil {
		return nil, &pluginError{Code: "invalid_input", Message: err.Error()}
	}

	client, err := newS3Client(ctx, cfg)
	if err != nil {
		return nil, &pluginError{Code: "internal", Message: err.Error()}
	}

	// Get remote object metadata.
	headResp, err := client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(cfg.Bucket),
		Key:    aws.String(cfg.Key),
	})

	status := map[string]any{}

	if err != nil {
		var nf *types.NotFound
		if errors.As(err, &nf) {
			status["remote_exists"] = false
		} else {
			return nil, &pluginError{Code: "internal", Message: fmt.Sprintf("s3 head: %s", err)}
		}
	} else {
		status["remote_exists"] = true
		if headResp.LastModified != nil {
			status["remote_last_modified"] = headResp.LastModified.UTC().Format("2006-01-02T15:04:05Z")
		}
		if headResp.ContentLength != nil {
			status["remote_size_bytes"] = *headResp.ContentLength
		}
	}

	// Get local file info.
	tasksPath := filepath.Join(root, "tasks.jsonl")
	fi, err := os.Stat(tasksPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			status["local_exists"] = false
		} else {
			return nil, &pluginError{Code: "internal", Message: fmt.Sprintf("stat tasks.jsonl: %s", err)}
		}
	} else {
		status["local_exists"] = true
		status["local_last_modified"] = fi.ModTime().UTC().Format("2006-01-02T15:04:05Z")
		status["local_size_bytes"] = fi.Size()
	}

	return status, nil
}

// --- main ---

func writeJSON(v any) {
	b, _ := json.Marshal(v)
	os.Stdout.Write(b)
	os.Stdout.Write([]byte("\n"))
}

func checkHealth() health {
	// We verify by attempting to load default config; real connectivity is checked on first operation.
	// The adapter is always considered "ready" at startup — errors surface on push/pull/status.
	return health{State: "ready"}
}

func main() {
	writeJSON(handshakeResponse{
		ProtocolVersion: protocolVersion,
		Name:            "s3",
		Kind:            "sync",
		Capabilities:    []string{"pull", "push", "status"},
		Health:          checkHealth(),
	})

	scanner := bufio.NewScanner(os.Stdin)
	// Increase scanner buffer for large payloads.
	scanner.Buffer(make([]byte, 1<<20), 1<<20)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var req requestEnvelope
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			continue
		}

		ctx := context.Background()

		var payload map[string]any
		var perr *pluginError

		switch req.Method {
		case "push":
			payload, perr = handlePush(ctx, req.Payload)
		case "pull":
			payload, perr = handlePull(ctx, req.Payload)
		case "status":
			payload, perr = handleStatus(ctx, req.Payload)
		default:
			perr = &pluginError{Code: "not_supported", Message: "unknown method: " + req.Method}
		}

		writeJSON(responseEnvelope{ID: req.ID, Payload: payload, Error: perr})
	}
}
