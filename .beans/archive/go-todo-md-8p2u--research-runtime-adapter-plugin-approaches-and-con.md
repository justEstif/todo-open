---
# go-todo-md-8p2u
title: Research runtime adapter plugin approaches and constraints
status: completed
type: task
priority: high
created_at: 2026-03-06T00:19:42Z
updated_at: 2026-03-06T00:21:58Z
parent: go-todo-md-dwib
---

Investigate viable plugin architectures for todoopen runtime adapters (process model, protocol, isolation, security, distribution). Capture tradeoffs and recommend a direction aligned with current server-first/local-first architecture.

## Deliverables
- [x] Survey at least 2-3 implementation approaches
- [x] Evaluate handshake/capability/error-reporting options
- [x] Identify operational/security risks and mitigations
- [x] Recommend a preferred approach with rationale

## Summary of Changes

- Added research doc at docs/adapter-plugin-research.md with comparison of three runtime plugin approaches:
  - in-process Go plugin shared objects
  - out-of-process stdio JSON-RPC
  - out-of-process gRPC/socket
- Evaluated handshake, capabilities, and structured error model options for adapter plugin protocol.
- Documented key operational/security risks and mitigations (supply chain, hung/crashing plugins, resource abuse, protocol drift, migration ambiguity).
- Recommended out-of-process stdio JSON-RPC as the MVP path due to best balance of isolation, simplicity, and local-first deployability.
- Added planning inputs and open questions for next bean (go-todo-md-f11d).
