---
# go-todo-md-5mgo
title: Explain parent_id migration approach
status: completed
type: task
priority: normal
created_at: 2026-03-05T18:27:19Z
updated_at: 2026-03-05T18:27:38Z
---

Explain safe migration strategies for introducing or evolving parent_id hierarchy semantics in task schema.\n\n## Todo\n- [x] Outline migration risks and compatibility concerns\n- [x] Propose phased migration strategy\n- [x] Share concrete examples and guardrails

## Summary of Changes\n\n- Documented key parent_id migration risks: orphan references, cycles, ordering assumptions, and compatibility with older clients.\n- Proposed phased rollout: optional read-only support, validation warnings, write support behind version gate, then stricter enforcement.\n- Added guardrails and examples for backfill, integrity checks, and rollback strategy.
