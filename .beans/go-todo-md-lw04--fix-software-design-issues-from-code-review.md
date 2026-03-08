---
# go-todo-md-lw04
title: Fix software design issues from code review
status: completed
type: task
priority: normal
created_at: 2026-03-08T22:02:15Z
updated_at: 2026-03-08T22:06:21Z
---

Address all findings from philosophy-of-software-design review:

- [x] Remove store.TaskRepository pass-through file
- [x] Move ensureWorkspace to construction time in jsonl.TaskRepo
- [x] Fix MutationEvent → events.Event manual copy (single conversion point)
- [x] Fix Update handler implicit dual-dispatch
- [x] Eliminate getAgentExt JSON round-trip with typed ext wrapper
- [x] Fix misleading origShutdown comment/dead code in app/server.go
- [x] Extract shared applyTitleUpdate helper to reduce UpdateTask/UpsertTask duplication

## Summary of Changes

All six findings from the design review addressed:

1. **Deleted**  — was a zero-value pass-through type alias
2. **** now returns  and runs  once at construction; removed per-operation workspace checks
3. **** added — single conversion point from  to ;  now calls 
4. **** added with  struct — unified partial-update path replaces the handler's implicit dispatch between  and ; consistent  semantics
5. **** changed from  to  — // simplified to a single / with no double round-trip
6. **Removed** dead  variable and misleading comment in 
7. **Extracted**  helper — shared by  and 
