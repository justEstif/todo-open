---
# go-todo-md-uxmv
title: Expose plugin-backed adapter status in CLI and API
status: completed
type: task
priority: normal
created_at: 2026-03-06T00:20:09Z
updated_at: 2026-03-06T00:45:15Z
parent: go-todo-md-dwib
---

Update todoopen adapters and GET /v1/adapters to report plugin-backed adapter readiness, health, and relevant diagnostics.

## Deliverables
- [x] API response model includes plugin runtime status fields
- [x] CLI output surfaces plugin readiness/health clearly
- [x] Error envelopes/messages stay consistent with API conventions
- [x] Tests for status output and degraded plugin states

## Summary of Changes

- Extended adapter status model to include source field in API payload (builtin, plugin, unknown).
- Updated runtime status construction so built-in adapters are labeled source=builtin, plugin-backed adapters source=plugin, unresolved entries source=unknown.
- Updated todoopen adapters CLI command to use workspace metadata (.todoopen/meta.json) and plugin-aware runtime evaluation instead of legacy adapters.json config.
- Updated CLI output formatting to include source plus readiness/health diagnostics.
- Added CLI tests for adapter status reporting:
  - builtin source appears in adapters output
  - plugin source appears when plugin handshake succeeds
- Kept API error behavior consistent by returning existing runtime ready/errors envelope through GET /v1/adapters.
- Updated docs/adapters.md examples to reflect metadata-driven adapters command usage.
- Verified with go test ./cmd/todoopen ./internal/api ./internal/app and go test ./....
