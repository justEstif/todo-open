---
# go-todo-md-1kt4
title: Decide MVP sync scope and quality gates
status: completed
type: task
priority: normal
created_at: 2026-03-05T18:16:08Z
updated_at: 2026-03-05T18:58:39Z
parent: go-todo-md-yris
---

Finalize what sync capabilities ship in MVP and what test/release gates are required.\n\n## Todo\n- [x] Choose MVP sync capability set and explicit non-goals\n- [x] Set conflict-handling policy for MVP\n- [x] Define mandatory test and CI gates before coding milestones

## Summary of Changes\n\n- Confirmed decisions are captured in sync.md: adapter-first MVP sync with one file-exchange adapter, explicit in-scope and out-of-scope items, and server-authoritative conflict handling.\n- Documented conflict policy: deterministic field-level LWW with tie-breakers plus post-merge invariants and conflict artifact persistence.\n- Defined quality gates in sync.md for determinism, no silent data loss, conflict artifact generation, offline operation, and post-merge schema/lifecycle validity.\n- Verified alignment with api.md and architecture.md sync endpoints and server-first boundary.
