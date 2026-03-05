---
status: accepted
---

# todo.open Server API Architecture

This document defines the high-level API model for todo.open’s **server-first** architecture.

- The server API is the canonical product boundary.
- CLI/web/mobile/TUI are all API clients.
- Core semantics (schema, lifecycle, conflict rules) are enforced server-side.

---

## 1) API Surfaces

## A. Task API (Core Domain)

Purpose: canonical task lifecycle and query surface.

Responsibilities:

- Create/read/update/archive tasks
- Validate against core + extension schema
- Enforce lifecycle transitions and invariants
- Expose query/filter/sort/pagination

Conceptual endpoints:

- `POST /v1/tasks`
- `GET /v1/tasks/{id}`
- `GET /v1/tasks`
- `PATCH /v1/tasks/{id}`
- `POST /v1/tasks/{id}/complete`
- `POST /v1/tasks/{id}/archive`
- `POST /v1/validate`

## B. Views API

Purpose: reusable, server-evaluated views for multiple clients.

Responsibilities:

- Register/list view definitions
- Materialize view results from canonical tasks
- Execute view-scoped actions

Conceptual endpoints:

- `GET /v1/views`
- `POST /v1/views`
- `GET /v1/views/{id}`
- `POST /v1/views/{id}/render`
- `POST /v1/views/{id}/actions/{actionId}`

## C. Sync API

Purpose: adapter-backed change exchange and conflict handling.

Responsibilities:

- Push/pull changes through configured adapters
- Manage sync checkpoints/tokens
- Surface conflicts and resolution paths

Conceptual endpoints:

- `GET /v1/sync/status`
- `POST /v1/sync/pull`
- `POST /v1/sync/push`
- `GET /v1/sync/conflicts`
- `POST /v1/sync/conflicts/{id}/resolve`

## D. Admin/Operations API (MVP-light)

Purpose: operational visibility and local deployment control.

Conceptual endpoints:

- `GET /healthz`
- `GET /readyz`
- `GET /v1/meta`
- `GET /v1/version`

---

## 2) Data and Contract Rules

- Canonical task schema: see `schema.md`
- Extension namespace: `ext.*`
- Server is source of truth for:
  - IDs and timestamp invariants
  - lifecycle transitions
  - conflict metadata shape
- API versioning uses explicit route prefix (`/v1/...`)

---

## 3) Client Model

All clients consume the same server contracts:

- `cmd/todoopen` (CLI client)
- web/mobile UI
- TUI and external integrations

Clients should not mutate JSONL files directly in server mode.

---

## 4) Deployment and Transport

MVP default:

- Local server over loopback HTTP
- Local filesystem JSONL backend

Later:

- Optional hosted/self-hosted remote mode
- authn/authz hardening and policy layers

---

## 5) Internal Mapping (Go)

High-level package responsibility mapping:

- `internal/api` → transport handlers, request/response mapping
- `internal/core` → domain rules and validation
- `internal/store` → persistence contracts and filesystem store
- `internal/sync` → sync orchestration and adapters
- `internal/app` → dependency wiring/composition root

---

## 6) Initial Non-Goals

- Real-time collaboration protocol
- Full enterprise permissions matrix
- Event-sourced/distributed consistency model

These can be layered after stable single-node server semantics are proven.
