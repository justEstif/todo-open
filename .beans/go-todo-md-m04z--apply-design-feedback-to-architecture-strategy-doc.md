---
# go-todo-md-m04z
title: Apply design feedback to architecture strategy docs
status: completed
type: task
priority: normal
created_at: 2026-03-05T19:11:50Z
updated_at: 2026-03-05T19:13:34Z
parent: go-todo-md-yris
---

Implement prioritized design feedback across architecture/api/sync/mvp/schema docs: remove contract duplication, align MVP wording to server-first, and formalize storage evolution policy.\n\n## Todo\n- [x] Make sync contract single-source-of-truth via api.md reference\n- [x] Update MVP scope wording to server-first API + CLI as first client\n- [x] Resolve architecture storage notes with explicit evolution trigger policy\n- [x] Add strict/compat default guidance for server mode\n- [x] Share updated docs and changes

## Summary of Changes\n\n- Updated sync.md to avoid duplicating sync endpoint contracts; it now references api.md as the canonical HTTP contract source.\n- Updated mvp.md scope to align with server-first direction: CRUD via server API with CLI as first client.\n- Updated architecture.md to remove inline uncertainty notes, clarify loopback meaning, and add an explicit storage evolution policy with threshold-based triggers.\n- Updated schema.md validation section to define default mode policy: strict in server mode, compat only for explicit import/migration flows.
