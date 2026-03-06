---
# go-todo-md-lgs6
title: Add zero-friction quickstart + persistence smoke test
status: completed
type: task
priority: high
created_at: 2026-03-05T23:13:03Z
updated_at: 2026-03-06T01:14:15Z
parent: go-todo-md-ua64
---

Add copy-paste quickstart and a restart persistence proof flow in README (plus optional script) so HN users can verify core value in under 2 minutes.

## Summary of Changes

- Added a copy-paste quickstart persistence section to README.md ("2-minute persistence smoke test").
- Added executable helper script: scripts/persistence-smoke.sh.
- Script flow: start server, create task, restart server, verify task persists.
- Verified script locally: PASS.
