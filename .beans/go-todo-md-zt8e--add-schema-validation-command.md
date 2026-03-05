---
# go-todo-md-zt8e
title: Add schema validation command
status: completed
type: task
priority: normal
created_at: 2026-03-05T20:55:52Z
updated_at: 2026-03-05T21:04:39Z
parent: go-todo-md-y0ga
---

Add a validation command to verify task records against the core schema and extension-field rules.

## Todo
- [x] Define validation rules aligned with docs/schema.md
- [x] Implement validation command in CLI and/or server endpoint
- [x] Surface actionable validation errors with line/context
- [x] Add tests for valid/invalid records

## Summary of Changes

- Added a schema validator in `internal/core/validation.go` covering required fields, types, enums, RFC3339 UTC timestamps, `done`/`completed_at` invariants, integer bounds, tag uniqueness, and strict-vs-compat unknown field handling.
- Added CLI subcommand `todoopen validate` in `cmd/todoopen/main.go` with `--file` and `--mode (strict|compat)` flags.
- Validation output now reports actionable issues with line number, field, and raw line context.
- Added tests in `internal/core/validation_test.go` for valid records, invalid records (including line/context assertions), and compat mode behavior.
- Verified with `go test ./...`.
