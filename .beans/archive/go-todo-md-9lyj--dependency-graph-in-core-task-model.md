---
# go-todo-md-9lyj
title: Dependency graph in core task model
status: completed
type: feature
priority: normal
created_at: 2026-03-08T19:36:02Z
updated_at: 2026-03-08T19:58:23Z
parent: go-todo-md-9yck
---

Add task dependency primitives to the core model. These are genuine task-data concerns — meaningful to humans and agents alike.

## What belongs in core
- `trigger_ids: []string` — task stays `pending` until all listed tasks reach `done`. Maps to GitHub Actions `needs:`.
- `pending` status — new status meaning 'not ready yet, waiting on dependencies'. State machine: `pending → open → in_progress → done`.
- Blocking/blocked-by relationships — bidirectional edges on the task, stored in the JSONL record. DFS cycle detection on write.

## What does NOT belong here
Lease/claim/heartbeat state. That's runtime coordination, not task data.

## Schema additions
- `trigger_ids` field in core.Task
- `pending` added to TaskStatus enum
- Server evaluates trigger_ids on any task completion: if all triggers are `done`, transition dependents `pending → open`

## API additions
- Filter support: `GET /v1/tasks?status=pending`, `GET /v1/tasks?blocked=true`
- `POST /v1/tasks/{id}/complete` evaluates and opens newly unblocked tasks

## Notes
- tasks.jsonl stays clean and human-readable — these are task relationships, not machine bookkeeping
- Aligns with beans blocking/blocked-by model but adds trigger semantics for fan-in pipelines

## Summary of Changes

- Added  to core task status enum ()
- Added , ,  fields to  struct
- Updated  interface:  now accepts ; added 
- Implemented  in : sets status=done, then calls  which auto-transitions pending tasks to open when all their trigger_ids are done
- Added  DFS helper (available for future use)
- Updated  to accept  as a valid status and register the three new fields in 
- Handler:  reads  and  query params; added  handler
- Router: registered 
- Tests:  (pending→open evaluation, filter),  integration test for complete endpoint
- Updated  and  to reflect all changes
