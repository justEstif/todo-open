---
# go-todo-md-rxnd
title: Remove unnecessary adapter-plugin-migration doc
status: completed
type: task
priority: normal
created_at: 2026-03-06T00:50:34Z
updated_at: 2026-03-06T00:51:14Z
parent: go-todo-md-dwib
---

Delete docs/adapter-plugin-migration.md and clean up references.

## Deliverables
- [x] Delete docs/adapter-plugin-migration.md
- [x] Remove or update references to deleted migration doc
- [x] Confirm docs/tests still pass basic checks

## Summary of Changes

- Deleted docs/adapter-plugin-migration.md.
- Updated remaining references in README.md, docs/adapters.md, and docs/schema.md to point to current canonical docs instead.
- Verified no remaining references to adapter-plugin-migration.md.
- Ran go test ./... successfully.
