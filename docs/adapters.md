---
status: accepted
---

# Adapter Contracts (Views and Sync)

This document defines the MVP extension workflow for runtime adapters.

## Goals

- Keep caller-side usage simple.
- Isolate adapter implementation details behind small, stable contracts.
- Allow adapters to be added without changing core task semantics.

## View Adapters

Package: `internal/view`

Contract:

- `Name() string`: stable identifier used for registry lookup.
- `RenderTasks(context.Context, []core.Task) ([]byte, error)`: render task data for a view target.

Registry behavior:

- Register adapters by unique non-empty name.
- Duplicate names are rejected.
- Callers fetch adapters by name and use them through the interface only.

Reference implementation:

- `internal/view/json.Adapter` renders `{"items": [...]}` JSON.

## Sync Adapters

Package: `internal/sync`

Contract:

- `Name() string`: stable identifier used for registry lookup.
- `Push(context.Context, []core.Task) error`: send local tasks/changes to a remote.
- `Pull(context.Context) ([]core.Task, error)`: fetch tasks/changes from a remote.

Registry behavior:

- Register adapters by unique non-empty name.
- Duplicate names are rejected.
- Callers fetch adapters by name and use them through the interface only.

Reference implementation:

- `internal/sync/noop.Adapter` provides a no-op baseline implementation for wiring/tests.

## Extension Workflow

1. Create a package implementing the relevant interface (`internal/view` or `internal/sync`).
2. Give the adapter a stable `Name()` value.
3. Register the adapter during app composition.
4. Resolve by name where needed; avoid concrete-type coupling outside adapter package.
5. Add focused unit tests for contract conformance.

This keeps extension complexity inside adapters and minimizes change amplification for the rest of the codebase.
