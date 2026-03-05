---
# go-todo-md-83b7
title: Bootstrap Go module and app entrypoints
status: completed
type: task
priority: normal
created_at: 2026-03-05T19:36:54Z
updated_at: 2026-03-05T20:13:18Z
parent: go-todo-md-y0ga
---

Create initial Go project scaffolding aligned to architecture docs.

## Todo
- [x] Add module layout under cmd and internal packages
- [x] Add cmd/todoopen-server main with config + HTTP bootstrap
- [x] Add cmd/todoopen CLI main wired as first API client
- [x] Ensure go test and go build pass for skeleton

## Scope additions
- [x] Add mise.toml with Go tool configured
- [x] Research and document Go coding standard for this repo

## Summary of Changes
- Bootstrapped Go module/layout with server and CLI entrypoints under cmd/ and internal/ packages.
- Added and configured mise.toml with Go tool plus modernization tasks (modernize-check, modernize).
- Researched modern Go guidance and documented repository coding standards in docs/coding-standards.md.
- Linked coding standards doc from README project documentation list.
- Closed zellij research session: go-standards-research.
