---
# go-todo-md-0547
title: Research sync and conflict-resolution approaches
status: completed
type: task
priority: normal
created_at: 2026-03-05T18:16:01Z
updated_at: 2026-03-05T18:40:45Z
parent: go-todo-md-yris
---

Investigate MVP-appropriate sync designs for local-first JSONL data.\n\n## Todo\n- [x] Evaluate sync architecture options (file exchange, simple remote API, plugin adapter)\n- [x] Document conflict-resolution strategies and complexity\n- [x] Recommend MVP sync approach and defer list

## Summary of Changes

- Added sync-research.md comparing three sync architecture options: file exchange, simple remote API, and adapter-first plugin model.
- Documented conflict resolution alternatives and complexity tradeoffs (record LWW, field-level LWW, CRDT/OT).
- Recommended MVP approach: adapter-first architecture with one file-exchange adapter, field-level LWW merge, conflict logging, and explicit defer list.
