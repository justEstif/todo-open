---
# go-todo-md-7s9f
title: Agent coordination layer (claim, heartbeat, work queue)
status: todo
type: feature
priority: normal
created_at: 2026-03-08T19:36:50Z
updated_at: 2026-03-08T19:36:54Z
parent: go-todo-md-9yck
blocked_by:
    - go-todo-md-9lyj
    - go-todo-md-k4yb
---

Thin coordination layer for agents claiming and executing tasks. Runtime-only — no ephemeral state pollutes tasks.jsonl.

## Design principle
Lease/claim state (`agent_id`, `lease_expires_at`, `heartbeat_at`) is ephemeral machine bookkeeping — it does NOT belong in the core task schema. It lives in `ext` or a separate in-memory/coordination store and evaporates when the agent is done. The task file records outcomes (`status: done`), not runtime mechanics.

## API surface

### Work queue
`GET /v1/tasks/next` — returns highest-priority unclaimed `open` task. The agent work queue.

### Atomic claim
`POST /v1/tasks/{id}/claim` — atomically transitions `open → in_progress`, sets `ext.agent.id` and `ext.agent.lease_expires_at`. Returns 409 if already claimed. The only safe entry point for agents — prevents two agents racing to the same task.

### Heartbeat
`POST /v1/tasks/{id}/heartbeat` — extends lease TTL (default 5m). Agents call this periodically while working.

### Release
`POST /v1/tasks/{id}/release` — agent explicitly gives a task back to `open`. Clears lease fields.

## Sweeper
Background goroutine checks for tasks where `ext.agent.lease_expires_at < now`, transitions them back to `open`, clears lease fields. Prevents tasks getting stuck if an agent crashes.

## ETag enforcement
All mutations require `If-Match` header (current `version` as ETag). Prevents lost-update races. Foundational — implement before claim/heartbeat.

## Blocked by
- go-todo-md-9lyj (dependency graph, for `pending` status and trigger evaluation on complete)
- go-todo-md-k4yb (SSE, so agents can watch for `task.status_changed: open` instead of polling)
