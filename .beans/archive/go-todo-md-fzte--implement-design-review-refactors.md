---
# go-todo-md-fzte
title: Implement design review refactors
status: completed
type: task
priority: normal
created_at: 2026-03-05T23:28:08Z
updated_at: 2026-03-05T23:32:35Z
---

Apply agreed high-impact design refactors: shared adapter registry component, consistent Create semantics across repos with tests, remove duplicate adapter runtime DTOs, and simplify CLI task command dispatch.

## Summary of Changes

- Added shared generic adapter registry: internal/adapterregistry/registry.go.
- Refactored internal/view and internal/sync registries to delegate to the shared registry while preserving existing public behavior.
- Consolidated adapter runtime DTOs into internal/adapters/types.go and updated router/handlers/app/tests to use shared types.
- Made memory repository Create semantics consistent with JSONL by rejecting duplicate IDs with core.ErrInvalidInput.
- Added cross-repo contract test: internal/store/repository_contract_test.go (duplicate-ID create behavior enforced for memory + JSONL).
- Refactored CLI task subcommand dispatch in cmd/todoopen/main.go into a command map + dedicated handlers for lower cognitive load.
- Ran gofmt and validated with go test ./... (all passing).
