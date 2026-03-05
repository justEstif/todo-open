---
# go-todo-md-ankw
title: Make README more user-friendly
status: completed
type: task
priority: normal
created_at: 2026-03-05T22:07:47Z
updated_at: 2026-03-05T22:08:12Z
---

Restructure README for end users: quick install, quick start, common commands, and links to deeper docs while keeping contributor setup clear.

## Summary of Changes

- Reworked `README.md` to be more user-friendly and task-oriented.
- Added a **Quick install** section with:
  - mise global install via Go backend
  - mise project-local install
  - Go-from-source install
- Added a **Quick start** section emphasizing `todoopen web` (and `gui` alias) with key flags.
- Clarified server-direct usage as source-checkout command (`go run ./cmd/todoopen-server`).
- Added a concise **For contributors** section with common `mise run` checks.
- Kept/organized links to deeper project docs under a single Documentation section.
