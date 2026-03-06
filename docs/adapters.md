---
status: accepted
---

# Adapter Contracts (Views and Sync)

This document defines the current adapter contracts and runtime behavior.

## Runtime Configuration Model

- Built-in adapters are registered in Go at startup.
- Plugin-backed adapters are discovered via `.todoopen/meta.json`.
- Runtime enablement source of truth is `.todoopen/meta.json`.

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

Built-in implementations:

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

Built-in implementations:

- `internal/sync/noop.Adapter` provides a no-op baseline implementation for wiring/tests.

## Example Extensions

Adapter examples are maintained in a separate examples repository.
Use those examples as templates for project-specific adapters.

## Current Adapter Configuration Schema

Runtime adapter config is loaded from workspace metadata in `.todoopen/meta.json`.
If metadata is missing, defaults are applied (`json` view + `noop` sync enabled).

```json
{
  "workspace_version": 1,
  "schema_version": "todo.open.task.v1",
  "enabled_views": ["json", "markdown"],
  "enabled_sync_adapters": ["noop", "git"],
  "adapter_plugins": [
    {"name": "markdown", "kind": "view", "command": "todoopen-plugin-view-markdown"},
    {"name": "git", "kind": "sync", "command": "todoopen-plugin-sync-git"}
  ],
  "ext": {
    "adapter_settings": {
      "git": {"remote": "origin", "branch": "tasks"}
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
todoopen adapters --workspace /path/to/workspace --json
```

Or query the server endpoint:

```bash
curl -s http://127.0.0.1:8080/v1/adapters
```

`GET /v1/adapters` returns:

- `config`: resolved adapter config (`enabled_views`, `enabled_sync_adapters`, optional settings)
- `status`: per-adapter entries with `kind`, `name`, `source` (`builtin|plugin|unknown`), `enabled`, `healthy`, optional `message`
- `ready`: overall startup/runtime readiness
- `errors`: adapter initialization errors when not ready

The CLI command and HTTP endpoint expose the same runtime readiness model.

## Extension Workflow

1. Implement the relevant adapter contract (`internal/view` or `internal/sync`).
2. Use a stable adapter `Name()`.
3. For built-ins, register during app composition.
4. For external plugins, register and enable through `.todoopen/meta.json`.
5. Keep adapter-specific settings under metadata extension fields.
6. Add focused unit tests for contract conformance.

This keeps extension complexity inside adapters and minimizes change amplification for the rest of the codebase.
