---
# go-todo-md-7h4f
title: 'Improve clarity: README marketing + adapter config restructuring'
status: completed
type: task
priority: normal
created_at: 2026-03-08T19:13:34Z
updated_at: 2026-03-08T19:18:53Z
---

Two tracks:
Track A: README rewrite — sharper tagline, concrete opening example, adapter table, cleaner structure inspired by qry.
Track B: Config restructuring — separate adapter config into .todoopen/config.toml (TOML, named sections, env var expansion, settings inlined); strip adapter fields from meta.json.

## Summary of Changes

Track A — README rewrite:
- New tagline: 'A local task server with an open API — agent-ready, CLI-friendly, your data.'
- Opens with a concrete command + JSON output so purpose is clear in 10 seconds
- Added adapter table (name, kind, what it does)
- Replaced mermaid flowchart with ASCII art — simpler, less intimidating
- Removed verbose prose, tightened every section

Track B — Adapter config restructuring:
- Added `github.com/BurntSushi/toml` dependency
- New `.todoopen/config.toml` format: named `[adapters.<name>]` sections, inline settings, ${VAR} env expansion
- `meta.json` now holds only workspace identity (workspace_version, schema_version, default_sort)
- New `internal/app/adapterfile.go`: AdapterFileConfig, AdapterEntry, AdapterGroupConfig types + loader
- Removed adapter fields from WorkspaceMeta (enabled_views, enabled_sync_adapters, adapter_plugins, ext.adapter_settings)
- `BuildAdapterRuntimeFromMeta` replaced by `BuildAdapterRuntimeFromConfig` (takes AdapterFileConfig directly)
- Updated server.go and cmd/todoopen/main.go to load config.toml
- Updated all tests and integration tests to new format
- Updated docs/adapters.md config schema section
