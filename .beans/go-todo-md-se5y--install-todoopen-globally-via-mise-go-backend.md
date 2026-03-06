---
# go-todo-md-se5y
title: Install todoopen globally via mise go backend
status: completed
type: task
priority: normal
created_at: 2026-03-06T01:08:28Z
updated_at: 2026-03-06T01:09:08Z
---

Run mise global install for todoopen CLI using go backend syntax.

## Deliverables
- [x] Execute: mise use -g go:github.com/justEstif/todo-open/cmd/todoopen@latest
- [x] Confirm command completed successfully

## Summary of Changes

- Ran: mise use -g go:github.com/justEstif/todo-open/cmd/todoopen@latest
- Mise reported the tool installed successfully.
- Verified binary resolution with /Users/ebeyene/.local/share/mise/installs/go-github-com-just-estif-todo-open-cmd-todoopen/latest/bin/todoopen and execution via todoopen - server-first local task client

Usage:
  todoopen --help
  todoopen [--server URL]                # health check
  todoopen web [--addr ADDR] [--no-open] # launch web app
  todoopen validate [flags]
  todoopen task <create|list|get|update|delete> [flags]
  todoopen adapters [--workspace PATH] [--json].
- Note: plain  is not available in this non-interactive shell PATH unless mise activation/shims are loaded in the user shell profile.
