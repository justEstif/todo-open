---
# go-todo-md-r6mw
title: Remove adapters.json runtime path and publish metadata-only migration docs
status: completed
type: task
priority: normal
created_at: 2026-03-06T00:20:22Z
updated_at: 2026-03-06T00:47:50Z
parent: go-todo-md-dwib
---

Finalize migration by removing adapters.json runtime loading path and publishing plugin author/user migration docs and examples for metadata-only configuration.

## Deliverables
- [x] adapters.json runtime loading path removed (no fallback mode)
- [x] Migration documentation added (user + plugin author)
- [x] Example workspace/plugin config published
- [x] Release notes/changelog entry drafted with hard-cutover callout

## Summary of Changes

- Removed legacy adapters.json runtime loading code from application path by deleting LoadAdapterConfig/DefaultAdapterConfig usage and related validation helpers.
- Removed adapters.json-focused unit tests from internal/app/adapterconfig_test.go.
- Kept runtime adapter evaluation metadata-only via .todoopen/meta.json and plugin handshake checks.
- Updated docs/adapters.md to reflect metadata as current runtime source-of-truth and updated schema/examples.
- Updated README.md sync extension guidance to use .todoopen/meta.json plugin registration and ext.adapter_settings.
- Updated docs/adapter-plugin-migration.md current-state section to metadata-only runtime config.
- Added release-notes draft at docs/release-notes.md with explicit hard-cutover breaking-change callout and migration mapping.
- Verified with go test ./internal/app ./cmd/todoopen and go test ./....
