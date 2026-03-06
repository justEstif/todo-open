---
# go-todo-md-zq1g
title: Align README/docs with implemented behavior
status: completed
type: task
priority: high
created_at: 2026-03-05T23:13:03Z
updated_at: 2026-03-05T23:41:46Z
parent: go-todo-md-ua64
---

Audit and update README + docs/api.md to clearly separate implemented vs conceptual endpoints/features; avoid roadmap phrasing that reads as shipped.

## Checklist

- [x] Audit README for conceptual vs implemented wording
- [x] Audit docs/api.md for conceptual vs implemented wording
- [x] Update docs to clearly mark current behavior
- [x] Run targeted checks (format/tests if needed)
- [x] Add summary of changes

## Summary of Changes

- Updated `README.md` to add a clear implemented-vs-planned API status section with explicit shipped endpoints.
- Clarified README sync/view sections so adapter contracts are presented as extension guidance, not shipped sync/view HTTP features.
- Reworked `docs/api.md` to split implemented API surfaces from planned/conceptual routes and aligned endpoint lists with the current router.
- Ran `go test ./...` to confirm the repository remains green after documentation updates.
