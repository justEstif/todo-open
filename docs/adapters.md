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

## Adapter Configuration Schema

Runtime adapter config is loaded from `.todoopen/adapters.json` (or `TODOOPEN_ADAPTER_CONFIG`).
If the file is missing, defaults are applied (`json` view + `noop` sync enabled).

```json
{
  "enabled_views": ["json"],
  "enabled_sync_adapters": ["noop"],
  "view_settings": {
    "json": {
      "indent": 2
    }
  },
  "sync_settings": {
    "noop": {
      "dry_run": true
    }
  }
}
```

Validation rules:

- adapter names in `enabled_views` and `enabled_sync_adapters` must be non-empty
- names in each enabled list must be unique
- any enabled adapter must also be registered at startup

## Startup and Health Semantics

- **Invalid config JSON or duplicate/empty names**: server startup fails fast.
- **Enabled adapter is not registered**: startup fails fast with adapter-specific errors.
- **Missing config file**: treated as defaults (server starts).

Use the CLI to inspect runtime status:

```bash
todoopen adapters
todoopen adapters --config .todoopen/adapters.json --json
```

This command prints per-adapter `enabled`/`healthy` status and exits non-zero when runtime is not ready.

## Extension Workflow

1. Create a package implementing the relevant interface (`internal/view` or `internal/sync`).
2. Give the adapter a stable `Name()` value.
3. Register the adapter during app composition.
4. Add adapter-specific settings under `view_settings.<name>` or `sync_settings.<name>`.
5. Resolve by name where needed; avoid concrete-type coupling outside adapter package.
6. Add focused unit tests for contract conformance.

This keeps extension complexity inside adapters and minimizes change amplification for the rest of the codebase.
