---
# go-todo-md-8h6o
title: Expose adapter runtime status over HTTP
status: completed
type: feature
priority: high
created_at: 2026-03-05T22:14:45Z
updated_at: 2026-03-05T23:21:41Z
parent: go-todo-md-ua64
---

Add a user-facing API endpoint for adapter runtime status (prefer GET /v1/adapters; GET /v1/meta acceptable) that returns enabled/healthy state and startup/runtime diagnostics consistent with CLI adapters output. Include handler, router wiring, response schema docs, and API tests.

## Summary of Changes

- Added `GET /v1/adapters` HTTP endpoint with a dedicated handler and router wiring.
- Exposed adapter runtime payload over HTTP (`config`, per-adapter `status`, `ready`, and optional `errors`) consistent with CLI adapter runtime reporting.
- Added API integration coverage for the new endpoint in `internal/api/adapters_integration_test.go`.
- Updated API and adapters documentation to include the new endpoint and response schema details.
