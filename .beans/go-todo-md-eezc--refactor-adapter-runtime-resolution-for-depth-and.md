---
# go-todo-md-eezc
title: Refactor adapter runtime resolution for depth and context propagation
status: completed
type: task
priority: normal
created_at: 2026-03-06T00:58:04Z
updated_at: 2026-03-06T00:59:45Z
---

Apply design feedback: extract plugin probing into a deeper abstraction and pass context through runtime resolution.

## Deliverables
- [x] Introduce dedicated plugin probe/resolver abstraction under internal/app
- [x] Refactor BuildAdapterRuntimeFromMeta to orchestrate via new abstraction
- [x] Add context propagation instead of context.Background in runtime resolution path
- [x] Keep behavior equivalent and update tests/callers
- [x] Run go test ./...

## Summary of Changes

- Refactored `internal/app/adapterconfig.go` to introduce a deeper runtime resolution abstraction.
- Added `adapterRuntimeResolver` to own plugin-backed adapter resolution flow.
- Added `pluginProbe` to encapsulate plugin loader probing behavior.
- Updated `BuildAdapterRuntimeFromMeta` to accept `context.Context` and orchestrate via the resolver.
- Updated call sites in `internal/app/server.go` and `cmd/todoopen/main.go`.
- Updated `internal/app/adapterconfig_test.go` for the new context-aware signature.
- Ran `gofmt` and `go test ./...` (all passing).
