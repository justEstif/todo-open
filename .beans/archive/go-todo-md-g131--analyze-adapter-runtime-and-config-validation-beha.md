---
# go-todo-md-g131
title: Analyze adapter runtime and config validation behavior
status: completed
type: task
priority: normal
created_at: 2026-03-05T23:49:11Z
updated_at: 2026-03-05T23:50:36Z
---

Investigate internal/app/adapterconfig.go, adapter types, API/CLI status surfaces, and tests. Answer startup behavior for missing enabled adapters, current defaults, and files impacted by moving new adapters from built-in registration to examples-only.

## Summary of Changes

Investigated adapter runtime/config behavior across app composition, registries, API/CLI status surfaces, and tests. Documented missing-enabled-adapter startup behavior, current defaults, and the file touch list required to move markdown/table/git adapters from built-in registration to examples-only usage.
