---
# go-todo-md-tk7c
title: Define human-first task workflows and UX invariants
status: completed
type: feature
priority: normal
created_at: 2026-03-06T02:24:12Z
updated_at: 2026-03-08T21:22:03Z
parent: go-todo-md-9yck
---

Codify core human task workflows, guardrails, and UX invariants so agent primitives do not degrade human usability.



## Summary of Changes

Created comprehensive human-ux-invariants.md document that codifies how the system must behave from a human user's perspective. The document includes:

### Core sections written:
1. **Core principle** - Humans set intent, agents execute; humans never lose control
2. **Workflow inventory** - 10 detailed canonical human workflows covering:
   - Task creation (with/without dependencies)
   - Task listing and inspection
   - Task editing and completion
   - Task cancellation/deletion
   - Prioritization and reordering
   - Live progress monitoring via SSE
   - Conflict resolution scenarios
   - Completed work review

3. **UX invariants** - 10 hard rules CLI/TUI must never violate:
   - Destructive operations require confirmation
   - Agent claim state always visible
   - Conflict errors show what changed
   - Humans can override agent claims
   - Stable list output ordering
   - Informative empty states
   - Silent success prohibited
   - Local timestamps for humans
   - Agent lease expiry visibility
   - Priority/deletion as human-only operations

4. **Conflict resolution UX** - Detailed walkthrough with exact error message wording and options

5. **CLI output contract** - Standardized formats for:
   - Task list rows (columns, truncation, color hints)
   - Single task detail views (field ordering, display rules)
   - Success/error/confirmation message patterns

6. **TUI requirements checklist** - 18 required capabilities for TUI implementation

7. **Cross-references** - Links to related documentation

### Updated README.md:
Added comprehensive Documentation section with links to all project docs including the new human-ux-invariants.md.

This document now serves as the north star for CLI and TUI implementation, ensuring consistent, human-centric behavior across all clients.
