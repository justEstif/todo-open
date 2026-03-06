---
# go-todo-md-b3na
title: Make JSONL persistence the default server backend
status: completed
type: feature
priority: critical
created_at: 2026-03-05T23:13:03Z
updated_at: 2026-03-05T23:21:23Z
parent: go-todo-md-ua64
---

Wire internal/app.NewServer to durable JSONL storage by default (workspace-root aware), keep memory store only for tests/dev, and add restart-persistence integration coverage.

## Summary of Changes

- Switched `internal/app.NewServer` to use JSONL-backed persistence by default (`internal/store/jsonl`) using a resolved workspace root.
- Added workspace-root awareness via `TODOOPEN_WORKSPACE_ROOT` (fallback to current working directory).
- Added optional `TODOOPEN_STORE=memory` override to keep in-memory storage available for dev/test scenarios.
- Added restart persistence integration coverage in `internal/app/server_integration_test.go` to verify tasks survive server restarts.
