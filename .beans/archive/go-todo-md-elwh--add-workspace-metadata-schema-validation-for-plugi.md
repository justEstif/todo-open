---
# go-todo-md-elwh
title: Add workspace metadata schema + validation for plugin registration
status: completed
type: task
priority: normal
created_at: 2026-03-06T00:20:02Z
updated_at: 2026-03-06T00:35:46Z
parent: go-todo-md-dwib
---

Extend .todoopen/meta.json schema and validation to support plugin registration and adapter enablement as source of truth.

## Deliverables
- [x] Metadata keys/schema added for plugin registration
- [x] Strict validation with actionable errors implemented
- [x] Backward-compatible defaults documented
- [x] Tests for valid/invalid metadata cases

## Summary of Changes

- Added metadata schema types and loader in internal/app/workspace_meta.go.
- Added plugin registration schema support via adapter_plugins entries (name, kind, command, optional args).
- Implemented strict metadata decoding with unknown-field rejection and actionable validation errors.
- Added metadata defaults for omitted values: enabled views [json] and enabled sync adapters [noop].
- Added tests in internal/app/workspace_meta_test.go for:
  - default behavior when metadata file is missing
  - valid metadata with plugin registrations
  - invalid metadata cases (unknown fields, invalid kind, duplicates, missing command)
- Updated docs/schema.md workspace metadata section to document adapter plugin keys and defaults.
- Verified with go test ./internal/app and go test ./....
