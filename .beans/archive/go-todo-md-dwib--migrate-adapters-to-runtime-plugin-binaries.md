---
# go-todo-md-dwib
title: Migrate adapters to runtime plugin binaries
status: completed
type: epic
priority: high
created_at: 2026-03-06T00:14:10Z
updated_at: 2026-03-06T00:50:00Z
---

Shift adapter extensibility from compile-time registration to runtime installable plugins configured from workspace metadata (`.todoopen/meta.json`).

## Goals
- Let end users add custom view/sync adapters without rebuilding `todoopen`.
- Use workspace metadata as the source of truth for plugin registration and enablement.
- Remove .todoopen/adapters.json runtime config and use metadata-only adapter configuration.

## Scope
- Plugin contract design (handshake, capabilities, health, errors).
- Runtime plugin loader and registry integration.
- Metadata schema + validation updates for plugin registration.
- CLI/API status surfacing for plugin-backed adapters.
- Hard-cutover migration plan and docs for metadata-only runtime config.

## Phases
- [x] Define plugin protocol and process model
- [x] Implement plugin discovery/loading/runtime isolation
- [x] Add metadata keys and validation for registered plugins
- [x] Wire enabled adapters from metadata in server startup
- [x] Update todoopen adapters and GET /v1/adapters for plugin status
- [x] Remove adapters.json runtime loading path (no compatibility mode)
- [x] Publish metadata-only migration docs and examples
- [x] Publish migration docs and examples for plugin authors

## Success Criteria
- A prebuilt `todoopen` binary can use third-party adapter plugins installed separately.
- Adapter registration and enablement live in `.todoopen/meta.json`.
- Users can inspect plugin readiness/health through existing adapter status surfaces.
- Existing workspaces can migrate with a documented, low-risk path.

## Summary of Changes

- Implemented plugin protocol contract (handshake, capabilities, structured errors) in internal/plugin.
- Implemented plugin process discovery/loading and runtime isolation with startup handshake validation.
- Added workspace metadata schema + strict validation for adapter_plugins and metadata-based enablement.
- Wired server startup and todoopen adapters command to metadata-only adapter configuration from .todoopen/meta.json.
- Exposed plugin-backed adapter status in CLI/API with source field (builtin|plugin|unknown), readiness, health, and diagnostics.
- Removed adapters.json runtime loading path (hard cutover; no compatibility mode).
- Updated user-facing docs to metadata-only configuration model and plugin registration examples.
- Removed temporary planning/research/release-note docs created during migration cleanup per request.
