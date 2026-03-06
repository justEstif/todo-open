---
# go-todo-md-7sfk
title: Explain internal/sync/noop package purpose
status: completed
type: task
priority: normal
created_at: 2026-03-05T23:31:50Z
updated_at: 2026-03-05T23:32:17Z
---

Answer what @internal/sync/noop/ is for.

## Todo
- [x] Inspect @internal/sync/noop implementation and references
- [x] Reply with concise purpose explanation

## Summary of Changes

Reviewed the noop sync adapter and app wiring, then explained that it is a built-in no-op/default sync backend used to keep adapter runtime healthy while external sync integrations are not yet implemented.
