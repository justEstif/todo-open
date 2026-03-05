---
# go-todo-md-trr1
title: Use mise go backend for golangci-lint task execution
status: completed
type: task
priority: normal
created_at: 2026-03-05T21:55:09Z
updated_at: 2026-03-05T21:55:55Z
---

Switch lint task to use mise Go backend so golangci-lint is built with the current Go toolchain context.

## Summary of Changes
- Updated `mise.toml` lint task to run golangci-lint via mise Go backend:
  - `mise x go:github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.5.0 -- golangci-lint run`
- Verified `mise run lint` passes.
- Verified full parity still passes with `mise run ci`.
