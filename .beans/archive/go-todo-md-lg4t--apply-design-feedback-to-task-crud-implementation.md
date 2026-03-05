---
# go-todo-md-lg4t
title: Apply design feedback to task CRUD implementation
status: completed
type: task
priority: normal
created_at: 2026-03-05T20:33:01Z
updated_at: 2026-03-05T20:35:33Z
parent: go-todo-md-y0ga
---

Address review feedback from CRUD implementation:\n\n## Todo\n- [x] Unify repository interface ownership to avoid duplicated contracts\n- [x] Align task domain model/lifecycle with schema.md (or document MVP subset clearly)\n- [x] Tighten JSON request decoding (unknown fields/trailing payload)\n- [x] Re-run tests and update bean summary

## Summary of Changes
- Replaced handler-local service interface with `core.TaskService` to reduce duplicate abstractions.
- Unified repository contract ownership by aliasing `internal/store.TaskRepository` to canonical `core.TaskRepository`.
- Expanded `core.Task` status and fields toward schema v1 (status set, priority, timestamps, version, metadata fields).
- Updated service/repository mutation behavior to initialize and increment `version` on create/update/delete.
- Hardened JSON request parsing with `DisallowUnknownFields` and trailing-data checks in task handlers.
- Extended integration tests to cover strict JSON decoding and version increments.
- Verified changes with `go test ./...`.
