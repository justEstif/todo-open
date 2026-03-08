---
# go-todo-md-0kz3
title: Fix concurrent claim race condition in store Update
status: completed
type: bug
priority: normal
created_at: 2026-03-08T23:20:52Z
updated_at: 2026-03-08T23:21:00Z
---

Two goroutines could both successfully claim the same task because GetByID and Update in ClaimTask were not atomic. Fixed by adding optimistic concurrency control (version check) to both memory and JSONL store Update methods.

## Summary of Changes

- **internal/store/memory/repo.go**: Added version check in  — rejects updates where  with 
- **internal/store/jsonl/repo.go**: Same version check added to the JSONL store's `Update` method

The root cause was a classic check-then-act race: `ClaimTask` in `core/coordination.go` calls `GetByID` (acquires RLock, releases), then `Update` (acquires Lock). Two goroutines could both read the task as "open" before either wrote. Now the store enforces optimistic concurrency via the Version field, so the second writer gets `ErrConflict` (mapped to HTTP 409).
