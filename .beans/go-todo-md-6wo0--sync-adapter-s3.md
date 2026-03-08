---
# go-todo-md-6wo0
title: 'sync adapter: S3'
status: completed
type: feature
priority: normal
created_at: 2026-03-08T19:22:20Z
updated_at: 2026-03-08T21:58:50Z
---

Implement an S3 sync adapter that reads/writes tasks.jsonl to an object storage bucket.

Works with any S3-compatible provider: AWS S3, Cloudflare R2, Backblaze B2, MinIO.

## Scope
- Lives at `internal/sync/s3` (built-in) or `cmd/todoopen-plugin-sync-s3`
- Config: endpoint (for non-AWS providers), bucket, key/path, region
- Push: upload tasks.jsonl to the configured key
- Pull: download tasks.jsonl from the configured key
- Handle missing object (first push) gracefully
- Use ${VAR} env expansion for credentials (access key, secret) — never store in config file

## Config example
```toml
[adapters.s3]
  bin  = "todoopen-plugin-sync-s3"
  kind = "sync"

[adapters.s3.config]
  bucket   = "my-tasks"
  endpoint = "https://${R2_ACCOUNT}.r2.cloudflarestorage.com"
  region   = "auto"
  access_key = "${S3_ACCESS_KEY}"
  secret_key = "${S3_SECRET_KEY}"
```

## Priority
Build after git adapter. Broadest reach for non-developer users who want cloud backup without git infrastructure.

## Implementation Plan\n\n- [x] Create `adapters/todoopen-adapter-sync-s3/main.go`\n- [x] Implement handshake, push, pull, status handlers\n- [x] Use AWS SDK v2 (aws-sdk-go-v2) for S3 operations\n- [x] Support env var expansion for credentials\n- [x] Handle missing object (first push) gracefully\n- [x] Write README\n- [x] Build + test compile

## Summary of Changes\n\nCreated `adapters/todoopen-adapter-sync-s3/` with:\n- `main.go` — full plugin implementation (handshake + push/pull/status handlers)\n- `README.md` — config docs and provider examples\n\nAdded `aws-sdk-go-v2` (s3, config, credentials, smithy-go) to `go.mod`/`go.sum`.\n\nKey design decisions:\n- Mirrors the git adapter's plugin protocol pattern exactly\n- Uses `s3.EndpointResolverV2` + `UsePathStyle = true` for non-AWS providers\n- Credential chain falls back to AWS default (env, ~/.aws, IAM role) when access_key/secret_key are empty\n- First-push/pull with no remote object returns a no-op message, not an error
