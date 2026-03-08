---
# go-todo-md-x7y0
title: 'Research: agent coordination primitives for todo.open'
status: completed
type: task
priority: normal
created_at: 2026-03-08T19:26:58Z
updated_at: 2026-03-08T19:35:45Z
---

Research what primitives todo.open needs for multi-agent parallel task coordination

## Summary
Research completed via subagents. Findings informed the design split:
- Dependency graph (trigger_ids, pending status, blocking) → core task model
- SSE event stream → server layer, zero schema impact
- Lease/claim/heartbeat → thin coordination layer, runtime-only, ext namespace or separate store
- Never pollute tasks.jsonl with ephemeral agent runtime state
