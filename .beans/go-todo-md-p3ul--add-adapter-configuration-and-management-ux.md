---
# go-todo-md-p3ul
title: Add adapter configuration and management UX
status: completed
type: feature
priority: normal
created_at: 2026-03-05T22:03:10Z
updated_at: 2026-03-05T22:12:33Z
---

Design and implement user-facing configuration and runtime management for view/sync adapters, including config schema, loading/validation behavior, discovery/listing UX, and docs/examples for custom adapters.\n\n## Todo\n- [x] Define adapter configuration schema (enabled adapters + per-adapter settings)\n- [x] Implement config loading and validation in app composition\n- [x] Add user-facing adapter discovery/status surface (CLI and/or API)\n- [x] Define failure/health semantics for adapter initialization\n- [x] Document end-to-end adapter setup and extension workflow

## Summary of Changes\n- Added adapter runtime config schema and loader in  with defaults, JSON parsing, and validation.\n- Added adapter runtime status model to evaluate enabled/healthy adapters against registered adapters.\n- Wired config loading + runtime validation into server composition () with fail-fast startup semantics for invalid adapter setups.\n- Added  CLI command for user-facing adapter discovery/status output (table or JSON).\n- Documented adapter configuration schema, startup semantics, and CLI workflow in  and .\n- Added tests for config loading/validation/runtime status in .

## Summary of Changes (corrected)\n- Added adapter runtime config schema and loader in internal/app/adapterconfig.go with defaults, JSON parsing, and validation.\n- Added adapter runtime status model to evaluate enabled/healthy adapters against registered adapters.\n- Wired config loading + runtime validation into server composition in internal/app/server.go with fail-fast startup semantics for invalid adapter setups.\n- Added todoopen adapters CLI command for user-facing adapter discovery/status output (table or JSON).\n- Documented adapter configuration schema, startup semantics, and CLI workflow in docs/adapters.md and README.md.\n- Added tests for config loading/validation/runtime status in internal/app/adapterconfig_test.go.
