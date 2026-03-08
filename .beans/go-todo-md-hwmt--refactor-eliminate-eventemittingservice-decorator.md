---
# go-todo-md-hwmt
title: 'Refactor: eliminate EventEmittingService decorator and consolidate ETag enforcement'
status: completed
type: task
priority: normal
created_at: 2026-03-08T20:39:37Z
updated_at: 2026-03-08T20:44:36Z
---

Three targeted design improvements from philosophy-of-software-design review:

1. Expose old status from service layer (PatchStatus/CompleteTask return transition info) — eliminates read-before-write in EventEmittingService
2. Replace EventEmittingService decorator with a MutationHook on core.Service — removes 14-method pass-through class
3. Consolidate ETag enforcement into UpdateTask service method (accept ifMatch *int like UpsertTask) — removes duplicated logic in handler

Also: wire detectCycle or delete it, type Ext as map[string]any, make nowFn injectable in IdempotencyStore.



## Summary of Changes

Successfully implemented all three design improvements:

### Change 1: MutationHook replaces EventEmittingService
- **Added MutationHook and MutationEvent types** to internal/core/service.go
- **Added hook field and OnMutation setter** to Service struct  
- **Updated all mutation methods** (CreateTask, UpdateTask, UpsertTask, PatchStatus, DeleteTask, CompleteTask, ClaimTask, ReleaseTask) to call hook after successful mutations
- **Deleted internal/events/service.go** - eliminated the 14-method pass-through decorator entirely
- **Updated internal/app/server.go** to use MutationHook instead of EventEmittingService wrapper
- **Deleted internal/events/service_test.go** - no longer needed
- **Eliminated read-before-write pattern** in PatchStatus and CompleteTask by capturing old status before update

### Change 2: Consolidated ETag enforcement into UpdateTask  
- **Updated UpdateTask signature** to accept  parameter
- **Added ETag validation** inside UpdateTask method (returns ErrConflict on mismatch)
- **Updated TaskService interface** to match new signature
- **Updated internal/api/handlers/tasks.go Update handler** to parse If-Match and pass to service
- **Removed inline GetTask + version comparison** from handler - now centralized in service

### Change 3: Minor fixes
- **Simplified toMap function** in internal/core/coordination.go - removed JSON round-trip fallback, now simple type assertion
- **Deleted detectCycle function** from internal/core/service.go - was never called
- **Added TODO comment** in CreateTask to validate trigger_ids cycle on write
- **Updated NewIdempotencyStore** in internal/api/middleware/idempotency.go to accept optional nowFn parameter
- **Updated callers** in internal/api/router.go and internal/api/idempotency_test.go

### Verification
- ✅ go test ./... -race - passes
- ✅ go vet ./... - passes  
- ✅ gofmt -l . - clean
- ✅ All functionality preserved with cleaner architecture


## Summary of Changes

Successfully implemented all three design improvements:

### Change 1: MutationHook replaces EventEmittingService
- **Added MutationHook and MutationEvent types** to internal/core/service.go
- **Added hook field and OnMutation setter** to Service struct  
- **Updated all mutation methods** (CreateTask, UpdateTask, UpsertTask, PatchStatus, DeleteTask, CompleteTask, ClaimTask, ReleaseTask) to call hook after successful mutations
- **Deleted internal/events/service.go** - eliminated the 14-method pass-through decorator entirely
- **Updated internal/app/server.go** to use MutationHook instead of EventEmittingService wrapper
- **Deleted internal/events/service_test.go** - no longer needed
- **Eliminated read-before-write pattern** in PatchStatus and CompleteTask by capturing old status before update

### Change 2: Consolidated ETag enforcement into UpdateTask  
- **Updated UpdateTask signature** to accept ifMatch parameter
- **Added ETag validation** inside UpdateTask method (returns ErrConflict on mismatch)
- **Updated TaskService interface** to match new signature
- **Updated internal/api/handlers/tasks.go Update handler** to parse If-Match and pass to service
- **Removed inline GetTask + version comparison** from handler - now centralized in service

### Change 3: Minor fixes
- **Simplified toMap function** in internal/core/coordination.go - removed JSON round-trip fallback, now simple type assertion
- **Deleted detectCycle function** from internal/core/service.go - was never called
- **Added TODO comment** in CreateTask to validate trigger_ids cycle on write
- **Updated NewIdempotencyStore** in internal/api/middleware/idempotency.go to accept optional nowFn parameter
- **Updated callers** in internal/api/router.go and internal/api/idempotency_test.go

### Verification
- ✅ go test ./... -race - passes
- ✅ go vet ./... - passes  
- ✅ gofmt -l . - clean
- ✅ All functionality preserved with cleaner architecture
