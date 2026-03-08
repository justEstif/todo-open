---
status: accepted
---

# Adapter Contracts (Views and Sync)

This document defines the current adapter contracts and runtime behavior.

## Runtime Configuration Model

- Built-in adapters are registered in Go at startup.
- Plugin-backed adapters are discovered via `.todoopen/config.toml`.
- Runtime enablement source of truth is `.todoopen/config.toml`.
- Workspace identity metadata (version, schema) lives separately in `.todoopen/meta.json`.

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

## Adapter Configuration Schema

Adapter configuration lives in `.todoopen/config.toml`, separate from workspace
identity metadata (`meta.json`). If `config.toml` is absent, defaults apply
(`json` view + `noop` sync enabled).

```toml
[views]
  enabled = ["json", "markdown"]

[sync]
  enabled = ["noop", "git"]

[adapters.git]
  bin  = "todoopen-plugin-sync-git"
  kind = "sync"

[adapters.git.config]
  remote = "${GIT_REMOTE}"
  branch = "tasks"

[adapters.markdown]
  bin  = "todoopen-plugin-view-markdown"
  kind = "view"
```

Use `${VAR}` syntax in `[adapters.<name>.config]` values — todo.open expands them
from the environment at runtime so secrets never live in the file.

Validation rules:

- adapter names in `views.enabled` and `sync.enabled` must be non-empty and unique
- any non-builtin enabled adapter must have a matching `[adapters.<name>]` entry
- `bin` is required for every declared adapter entry
- `kind` must be `view` or `sync`

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
4. For external plugins, declare the binary in `.todoopen/config.toml` and add the name to `views.enabled` or `sync.enabled`.
5. Keep adapter-specific settings under `[adapters.<name>.config]`.
6. Add focused unit tests for contract conformance.

This keeps extension complexity inside adapters and minimizes change amplification for the rest of the codebase.
