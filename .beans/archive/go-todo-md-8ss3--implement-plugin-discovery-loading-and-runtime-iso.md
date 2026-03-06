---
# go-todo-md-8ss3
title: Implement plugin discovery, loading, and runtime isolation
status: completed
type: task
priority: normal
created_at: 2026-03-06T00:19:58Z
updated_at: 2026-03-06T00:32:30Z
parent: go-todo-md-dwib
---

Build runtime plugin loading and registration path with lifecycle management and failure isolation.

## Deliverables
- [x] Plugin discovery mechanism implemented
- [x] Loader + lifecycle (start/stop/timeout/retry) implemented
- [x] Runtime isolation/failure containment behaviors in place
- [x] Integration tests for load and failure scenarios

## Summary of Changes

- Added plugin runtime loader implementation in internal/plugin/runtime.go.
- Implemented command discovery:
  - supports direct path commands
  - supports PATH lookup via exec.LookPath
  - returns clear errors for missing/empty command
- Implemented process lifecycle management:
  - process start via exec.CommandContext
  - stdio pipe setup
  - startup handshake read with timeout
  - handshake validation against expected name/kind/protocol/capabilities
  - plugin close/kill behavior
- Implemented runtime isolation semantics:
  - handshake/read/validation failures kill plugin process and return host-side error
  - plugin process exit is surfaced as unhealthy state via Health()
- Added load/failure tests in internal/plugin/runtime_test.go covering:
  - successful load
  - missing command
  - handshake mismatch
  - handshake timeout
  - unhealthy state after plugin process exit
- Verified with go test ./internal/plugin and go test ./....
