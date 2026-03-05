---
# go-todo-md-qpqb
title: Refine README with mise Go backend install command
status: completed
type: task
priority: normal
created_at: 2026-03-05T22:07:03Z
updated_at: 2026-03-05T22:07:18Z
---

Update README install section to show users how to install todoopen via mise Go backend (dev-tools backends/go) in addition to local build/run steps.

## Summary of Changes

- Updated `README.md` install instructions to explicitly show installing `todoopen` via the **mise Go backend**.
- Added both global and project-local mise install examples:
  - `mise use -g go:github.com/justEstif/todo-open/cmd/todoopen@latest`
  - `mise use go:github.com/justEstif/todo-open/cmd/todoopen@latest`
- Kept contributor-focused setup instructions (`mise install`, `mise run build`) for local development from source.
