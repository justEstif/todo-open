---
# go-todo-md-pf54
title: Implement JSONL store behind repository interfaces
status: completed
type: task
priority: normal
created_at: 2026-03-05T19:37:01Z
updated_at: 2026-03-05T20:44:32Z
parent: go-todo-md-y0ga
---

Build the initial persistence layer behind clean interfaces in internal/store and internal/core boundaries.\n\n## Todo\n- [x] Define repository interfaces for task CRUD and list operations\n- [x] Implement JSONL read/write/update with atomic file writes\n- [x] Add metadata file handling and schema version checks\n- [x] Add unit tests for happy path and corruption/recovery scenarios

## Summary of Changes
- Implemented a functional JSONL repository in `internal/store/jsonl/repo.go` behind the core repository contract.
- Added task CRUD persistence over `tasks.jsonl` with atomic writes via temp-file rename.
- Added workspace metadata bootstrap/validation in `.todoopen/meta.json` with schema-version checks (`todo.open.task.v1`).
- Implemented JSONL parsing with corruption detection (line-numbered decode errors).
- Added unit tests in `internal/store/jsonl/repo_test.go` for CRUD happy path, metadata bootstrap, corruption handling, and schema mismatch rejection.
- Verified with `go test ./...`.
