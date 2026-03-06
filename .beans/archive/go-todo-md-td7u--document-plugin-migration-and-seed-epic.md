---
# go-todo-md-td7u
title: Document plugin migration and seed epic
status: completed
type: task
priority: normal
created_at: 2026-03-06T00:12:54Z
updated_at: 2026-03-06T00:15:02Z
---

Update docs to reflect runtime installable adapter plugins and metadata-based registration.

## Todo
- [x] Update docs to describe adapter plugin migration direction
- [x] Create an epic bean describing migration scope and phases
- [x] Add summary of changes and close this bean

## Summary of Changes

- Updated adapter docs to clearly separate current compile-time behavior from target runtime plugin model.
- Added `docs/adapter-plugin-migration.md` with migration goals, proposed `.todoopen/meta.json` shape, phased rollout, and non-goals.
- Linked schema and README docs to the migration direction.
- Created epic `go-todo-md-dwib` to track implementation phases for runtime plugin migration.
