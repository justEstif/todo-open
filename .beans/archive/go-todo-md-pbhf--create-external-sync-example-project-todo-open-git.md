---
# go-todo-md-pbhf
title: Create external sync example project todo-open-git-sync
status: completed
type: task
priority: normal
created_at: 2026-03-06T01:03:03Z
updated_at: 2026-03-06T01:04:06Z
---

Scaffold a standalone sync plugin example one level above the current repo.

## Deliverables
- [x] Create ../todo-open-git-sync project structure
- [x] Add standalone Go module + sync plugin example binary
- [x] Add README with build/run/config instructions
- [x] Run go build in the new project

## Summary of Changes

- Created ../todo-open-git-sync as a standalone Go module (go 1.26).
- Added plugin binary entrypoint at cmd/todoopen-plugin-sync-git/main.go.
- Implemented handshake and minimal request handling for required sync capabilities: pull, push, status.
- Added README with build instructions and .todoopen/meta.json wiring example.
- Ran gofmt and go build ./... successfully in the new project.
