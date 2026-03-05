---
# go-todo-md-jtb3
title: Decide runtime, architecture, and module layout
status: completed
type: task
priority: high
created_at: 2026-03-05T18:16:08Z
updated_at: 2026-03-05T18:38:10Z
parent: go-todo-md-yris
---

Make final implementation decisions for runtime, project structure, and package boundaries using research outputs.\n\n## Todo\n- [x] Select runtime and primary architecture direction (Go + server-first)\n- [x] Define module/package layout for core/storage/cli/sync\n- [x] Record rationale and rejected alternatives

## Decision Notes\n\n- Confirmed direction: build todo.open as server-first architecture with Go as implementation language.\n- Implication: CLI becomes one client/interface rather than the primary product boundary.

- Added architecture.md to document server-first component boundaries, deployment modes, request flow, and high-level Go package structure.\n- Rationale: server-first provides a stable API/domain boundary for multiple clients while preserving local-first behavior.\n- Rejected alternative: CLI-first as primary product boundary; retained CLI as first client implementation.

## Summary of Changes\n\n- Finalized runtime and architecture decision as Go + server-first for todo.open.\n- Added architecture.md describing the high-level vision, principles, components, deployment modes, and request flow.\n- Documented high-level Go package direction (internal/core, internal/store, internal/sync, internal/api, internal/app, and cmd entrypoints).\n- Captured rationale and rejected alternative (CLI-first as primary boundary).
