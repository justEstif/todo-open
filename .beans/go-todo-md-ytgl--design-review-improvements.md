---
# go-todo-md-ytgl
title: Design review improvements
status: completed
type: task
priority: normal
created_at: 2026-03-09T00:11:13Z
updated_at: 2026-03-09T00:13:00Z
---

Implement 4 improvements from philosophy-of-software-design review:
- [x] Fix PatchTask double-update (merge status+title into single repo write)
- [x] Deduplicate client API boilerplate (extract generic request helper)
- [x] Resolve TaskEvent duplication between client/api and events packages
- [x] Eliminate view/sync registry pass-throughs

## Summary of Changes

1. **PatchTask single atomic write** — merged status and title changes into one `repo.Update` call with one version bump, eliminating intermediate states
2. **Client API dedup** — extracted `do()` helper method, reducing 7 methods from ~100 lines to ~30 lines of boilerplate
3. **TaskEvent type alias** — replaced duplicate `TaskEvent` struct in client/api with `type TaskEvent = events.Event`, single source of truth
4. **Registry type aliases** — replaced pass-through `view.Registry` and `sync.Registry` wrappers with `type Registry = adapterregistry.Registry[Adapter]`, keeping error re-exports for backward compatibility
