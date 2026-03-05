---
# go-todo-md-ftnr
title: Implement concrete adapter contracts for views and sync
status: completed
type: task
priority: normal
created_at: 2026-03-05T20:56:01Z
updated_at: 2026-03-05T22:00:46Z
parent: go-todo-md-y0ga
---

Convert documented extension points into concrete runtime interfaces/contracts in code for views and sync adapters.

## Todo
- [x] Define Go interfaces/types for view adapters and sync adapters
- [x] Add adapter registration/loading strategy
- [x] Provide one reference implementation and tests
- [x] Document contract usage and extension workflow

## Summary of Changes
- Added `internal/view` and `internal/sync` adapter contracts with registries and error semantics.
- Added app-level adapter loading helpers in `internal/app/adapters.go` for built-in adapter registration.
- Added reference adapters: `internal/view/json` and `internal/sync/noop`.
- Added unit tests for registries, reference adapters, and app loader wiring.
- Added `docs/adapters.md` documenting contracts and extension workflow.
- Verified `go test ./...` and `mise run ci` pass.
