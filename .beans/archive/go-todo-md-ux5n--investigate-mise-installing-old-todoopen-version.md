---
# go-todo-md-ux5n
title: Investigate mise installing old todoopen version
status: completed
type: bug
priority: normal
created_at: 2026-03-06T03:21:00Z
updated_at: 2026-03-06T03:22:22Z
---

Diagnose why mise install still yields an old todoopen version after tagging v0.1.1 and adding CLI version embedding.

## Todo
- [x] Review current install instructions/config for mise
- [x] Check how version is embedded for non-release installs
- [x] Reproduce local install/version behavior
- [x] Propose and implement fix
- [x] Validate fix and summarize

## Summary of Changes

- Root cause 1: `mise use -g go:...@latest` resolves via Go module proxy and may lag immediately after new pushes/tags.
- Root cause 2: Go-module installs do not pass release ldflags, so CLI fell back to a default version string.
- Implemented fix in `cmd/todoopen/main.go`: version output now prefers injected build var, then falls back to `debug.ReadBuildInfo().Main.Version`, then `dev`.
- Verified with `go test ./cmd/todoopen` (pass).
- Practical guidance: install a pinned release tag (`@v0.1.1` or newer) for predictable upgrades.