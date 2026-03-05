---
# go-todo-md-ul4h
title: Replace ebeyene references with justEstif
status: completed
type: task
priority: normal
created_at: 2026-03-05T20:48:55Z
updated_at: 2026-03-05T20:49:50Z
parent: go-todo-md-y0ga
---

Update repository/code/docs references from ebeyene to justEstif.

## Todo
- [x] Find all occurrences of "ebeyene" in tracked files
- [x] Update references to "justEstif" where appropriate
- [x] Run tests/checks to ensure no breakage
- [x] Add summary and complete bean

## Summary of Changes
- Updated module path from `github.com/ebeyene/todo-open` to `github.com/justEstif/todo-open` in `go.mod`.
- Updated Go import paths across cmd/internal packages to the new module owner.
- Updated `AGENTS.md` module reference and other remaining `ebeyene` text occurrences in tracked files (including beans context text).
- Verified with `go test ./...`.
