---
# go-todo-md-zksz
title: Improve mise install path for tagged builds
status: completed
type: task
priority: normal
created_at: 2026-03-06T03:02:14Z
updated_at: 2026-03-06T03:03:11Z
---

User reported mise distribution needs a better build-script discovery path and support for install tags (go install --tags). Investigate current install/build docs and tooling, then implement a robust tagged install flow.

## Todo
- [x] Review current mise install flow and constraints
- [x] Add tagged install guidance using go install semantics
- [x] Validate command locally
- [x] Summarize changes

## Summary of Changes
- Updated  install section with tagged-install instructions for mise users.
- Documented mise ~/.config/mise/config.toml tools: go:github.com/justEstif/todo-open/cmd/todoopen@latest as the supported equivalent of .
- Added multi-tag example using comma-separated tags.
- Validated locally with: todoopen - server-first local task client

Usage:
  todoopen --help
  todoopen [--server URL]                # health check
  todoopen web [--addr ADDR] [--no-open] # launch web app
  todoopen validate [flags]
  todoopen task <create|list|get|update|delete> [flags]
  todoopen adapters [--workspace PATH] [--json].
