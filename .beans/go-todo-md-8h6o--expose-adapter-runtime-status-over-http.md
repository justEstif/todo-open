---
# go-todo-md-8h6o
title: Expose adapter runtime status over HTTP
status: todo
type: feature
created_at: 2026-03-05T22:14:45Z
updated_at: 2026-03-05T22:14:45Z
---

Add a user-facing API endpoint for adapter runtime status (prefer GET /v1/adapters; GET /v1/meta acceptable) that returns enabled/healthy state and startup/runtime diagnostics consistent with CLI adapters output. Include handler, router wiring, response schema docs, and API tests.
