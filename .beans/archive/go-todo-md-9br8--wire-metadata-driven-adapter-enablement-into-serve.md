---
# go-todo-md-9br8
title: Wire metadata-driven adapter enablement into server startup
status: completed
type: task
priority: normal
created_at: 2026-03-06T00:20:05Z
updated_at: 2026-03-06T00:41:49Z
parent: go-todo-md-dwib
---

Update app startup/composition to build enabled adapters from workspace metadata and plugin registry.

## Deliverables
- [x] Server startup reads adapter registration/enablement from metadata
- [x] Registry wiring updated for plugin-backed adapters
- [x] Startup diagnostics for missing/unhealthy plugins
- [x] Integration tests for startup behavior

## Summary of Changes

- Updated server startup to load adapter enablement/registration from .todoopen/meta.json (LoadWorkspaceMeta) instead of adapters.json.
- Added BuildAdapterRuntimeFromMeta in internal/app/adapterconfig.go to combine:
  - built-in adapter registry state
  - metadata-declared plugin registrations
  - plugin startup handshake checks
- Added startup diagnostics for plugin-backed adapters:
  - handshake failures are surfaced in runtime errors
  - missing registrations remain startup errors
  - successful plugin handshake marks adapter healthy
- Added helper logic to upsert adapter statuses and clear stale 'not registered' errors once plugin handshake succeeds.
- Added integration-style coverage in internal/app/adapterconfig_test.go:
  - valid plugin handshake path marks runtime ready
  - handshake mismatch path marks runtime not ready
- Added server startup test in internal/app/server_integration_test.go that validates startup fails when an enabled plugin handshake fails.
- Verified with go test ./internal/app and go test ./....
