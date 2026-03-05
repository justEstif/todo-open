---
# go-todo-md-dnrd
title: Implement first server CRUD API endpoints
status: completed
type: task
priority: normal
created_at: 2026-03-05T19:37:05Z
updated_at: 2026-03-05T20:29:42Z
parent: go-todo-md-y0ga
---

Deliver the first HTTP API slice for task CRUD aligned to api.md contracts.\n\n## Todo\n- [x] Implement handlers for create, get, list, update, delete task endpoints\n- [x] Add request validation and standard error response mapping\n- [x] Wire handlers to core services and repository interfaces\n- [x] Add integration tests for API happy path and key validation failures

## Summary of Changes
- Added  CRUD endpoints (create, list, get, update, delete) in .
- Implemented task handlers with request validation and standard JSON error mapping in .
- Expanded the core task service and domain model for CRUD + soft delete in .
- Added an in-memory repository implementation and updated repository contracts (, ).
- Added integration tests for happy path CRUD and key validation/not-found failures in .
- Verified with ?   	github.com/ebeyene/todo-open/cmd/todoopen	[no test files]
?   	github.com/ebeyene/todo-open/cmd/todoopen-server	[no test files]
ok  	github.com/ebeyene/todo-open/internal/api	(cached)
?   	github.com/ebeyene/todo-open/internal/api/handlers	[no test files]
?   	github.com/ebeyene/todo-open/internal/app	[no test files]
?   	github.com/ebeyene/todo-open/internal/client/api	[no test files]
?   	github.com/ebeyene/todo-open/internal/core	[no test files]
?   	github.com/ebeyene/todo-open/internal/store	[no test files]
?   	github.com/ebeyene/todo-open/internal/store/jsonl	[no test files]
?   	github.com/ebeyene/todo-open/internal/store/memory	[no test files].
