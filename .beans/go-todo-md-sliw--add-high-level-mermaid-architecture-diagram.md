---
# go-todo-md-sliw
title: Add high-level mermaid architecture diagram
status: completed
type: task
priority: normal
created_at: 2026-03-05T21:13:29Z
updated_at: 2026-03-05T21:13:57Z
---

Update docs/architecture.md with a mermaid diagram showing todo-open server, clients, sync flow, and core layers.

## Summary of Changes

- Added a new mermaid flowchart section to docs/architecture.md under **Top-Level System Components**.
- Diagram now shows high-level interactions among clients, server layers (internal/api, internal/core, internal/store, internal/sync), and a remote sync peer/provider.
- Included a durable JSONL relationship annotation between storage and core.
