---
# go-todo-md-qh7k
title: Implement CLI task CRUD commands
status: completed
type: task
priority: normal
created_at: 2026-03-05T20:55:52Z
updated_at: 2026-03-05T21:14:40Z
parent: go-todo-md-y0ga
---

Implement CLI commands to create/list/get/update/delete tasks via the server API so CLI can manage tasks end-to-end.

## Todo
- [x] Add CLI subcommands for task CRUD
- [x] Extend internal API client methods for task endpoints
- [x] Add output and error handling for CLI UX
- [x] Add tests for CLI command behavior

## Summary of Changes

- Added `task` CLI subcommands in `cmd/todoopen/main.go`: `create`, `list`, `get`, `update`, and `delete`.
- Refactored CLI entrypoint to `run(args, stdout, stderr)` for testability and clearer command dispatch.
- Extended API client in `internal/client/api/client.go` with CRUD methods: `CreateTask`, `ListTasks`, `GetTask`, `UpdateTask`, and `DeleteTask`.
- Improved CLI UX with structured success output (JSON or tabular list) and explicit command-level error messages.
- Added command behavior tests in `cmd/todoopen/main_test.go` covering CRUD happy path and an error-path network failure.
- Verified with `go test ./...`.
