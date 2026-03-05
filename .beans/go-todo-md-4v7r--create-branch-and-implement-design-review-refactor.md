---
# go-todo-md-4v7r
title: Create branch and implement design review refactors
status: completed
type: task
priority: normal
created_at: 2026-03-05T21:33:59Z
updated_at: 2026-03-05T21:36:34Z
---

## Objective
Create a dedicated git branch and apply high-impact refactors from the recent design review.

## Checklist
- [x] Create and switch to a new git branch
- [x] Move delete/archive lifecycle invariants into core service
- [x] Reduce JSONL repository change amplification by centralizing read-modify-write flow
- [x] Run go test ./... and report results

## Summary of Changes
- Created and switched to branch .
- Moved delete/archive lifecycle invariants into  by loading the task, applying domain state transition (, archived status, updated_at, version++), and persisting via .
- Simplified repository contract by removing  from  and removing  implementations from memory/jsonl stores.
- Reduced JSONL repository change amplification by introducing  and  helpers and refactoring Create/GetByID/List/Update to use them.
- Updated JSONL repository tests to match the new repository contract.
- Verified with ok  	github.com/justEstif/todo-open/cmd/todoopen	(cached)
?   	github.com/justEstif/todo-open/cmd/todoopen-server	[no test files]
ok  	github.com/justEstif/todo-open/internal/api	(cached)
?   	github.com/justEstif/todo-open/internal/api/handlers	[no test files]
?   	github.com/justEstif/todo-open/internal/api/web	[no test files]
?   	github.com/justEstif/todo-open/internal/app	[no test files]
?   	github.com/justEstif/todo-open/internal/client/api	[no test files]
ok  	github.com/justEstif/todo-open/internal/core	(cached)
?   	github.com/justEstif/todo-open/internal/store	[no test files]
ok  	github.com/justEstif/todo-open/internal/store/jsonl	(cached)
?   	github.com/justEstif/todo-open/internal/store/memory	[no test files] (all passing).

## Summary of Changes (corrected)
- Created and switched to branch `refactor/design-depth-pass-1`.
- Moved delete/archive lifecycle invariants into `core.Service.DeleteTask` by loading the task, applying domain state transition (`deleted_at`, archived status, `updated_at`, `version++`), and persisting via `repo.Update`.
- Simplified repository contract by removing `Delete` from `core.TaskRepository` and removing `Delete` implementations from memory/jsonl stores.
- Reduced JSONL repository change amplification by introducing `withTasksRead` and `withTasksMutation` helpers and refactoring Create/GetByID/List/Update to use them.
- Updated JSONL repository tests to match the new repository contract.
- Verified with `go test ./...` (all passing).
