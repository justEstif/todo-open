---
# go-todo-md-qnxq
title: Remove in-repo examples folder
status: completed
type: task
priority: normal
created_at: 2026-03-06T00:52:44Z
updated_at: 2026-03-06T00:53:07Z
---

Delete local examples/ directory and clean up references since examples should live in a separate repository.

## Deliverables
- [x] Delete examples/ directory
- [x] Remove/update references to examples/ paths in docs/code
- [x] Run test/build sanity checks after removal

## Summary of Changes

- Deleted the in-repo examples/ directory.
- Updated README.md references that pointed to examples/adapters paths.
- Updated docs/adapters.md example section to reference a separate examples repository.
- Verified repository health with go test ./....
