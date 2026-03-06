---
# go-todo-md-ykrg
title: Implement plugin protocol and process runtime contract
status: completed
type: task
priority: normal
created_at: 2026-03-06T00:19:54Z
updated_at: 2026-03-06T00:28:52Z
parent: go-todo-md-dwib
---

Implement the selected plugin protocol and process model, including handshake, capability negotiation, and standardized error semantics.

## Deliverables
- [x] Protocol types/interfaces in code
- [x] Handshake + capability negotiation implemented
- [x] Health/error contract implemented
- [x] Unit tests for contract semantics

## Summary of Changes

- Added new protocol contract package: internal/plugin.
- Implemented protocol/runtime contract types in internal/plugin/protocol.go:
  - adapter kind/capability enums
  - handshake request/response structs
  - health state model
  - transport-neutral request/response envelopes
  - structured plugin error model with stable error codes
- Implemented handshake and capability validation via ValidateHandshake, including required capability checks by adapter kind.
- Added unit tests in internal/plugin/protocol_test.go covering:
  - successful view and sync handshakes
  - version/name/kind mismatch failures
  - missing required capability failure
  - structured plugin error formatting.
- Verified with go test ./internal/plugin and go test ./....
