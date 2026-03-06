---
# go-todo-md-glc5
title: Switch sync example install docs to mise go backend install
status: completed
type: task
priority: normal
created_at: 2026-03-06T01:07:24Z
updated_at: 2026-03-06T01:07:48Z
---

Update ../todo-open-git-sync/README.md to install the plugin via mise Go backend syntax (go:module/cmd@version) instead of local build steps.

## Deliverables
- [x] Replace install section with mise go backend install command
- [x] Keep registration instructions intact
- [x] Verify README examples are consistent

## Summary of Changes

- Updated ../todo-open-git-sync/README.md install section to use mise Go backend syntax.
- Added command: mise use -g go:github.com/justEstif/todo-open-git-sync/cmd/todoopen-plugin-sync-git@latest.
- Kept workspace registration and verify sections aligned with command name todoopen-plugin-sync-git.
