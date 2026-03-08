---
# go-todo-md-6wo0
title: 'sync adapter: S3'
status: todo
type: feature
created_at: 2026-03-08T19:22:20Z
updated_at: 2026-03-08T19:22:20Z
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
