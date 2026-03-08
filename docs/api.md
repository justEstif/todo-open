---
status: accepted
---

# todo.open Server API Architecture

This document separates the **implemented HTTP API** from the **planned architecture** so roadmap concepts do not read as shipped behavior.

- The server API is the canonical product boundary.
- CLI/web/mobile/TUI are all API clients.
- Core semantics are enforced server-side for shipped routes, with additional lifecycle/conflict semantics defined for planned routes.

---

## 1) Implemented API Surfaces (Current)

### A. Task API (Core Domain)

Purpose: canonical task lifecycle and query surface.

Responsibilities currently implemented:

- Create/read/update/delete tasks
- Validate request payload shape and core task constraints
- Enforce core task invariants
- Expose list/get primitives

Implemented endpoints:

- `POST /v1/tasks`
- `GET /v1/tasks` — supports query params: `status=<status>`, `is_blocked=true`
- `GET /v1/tasks/{id}`
- `PATCH /v1/tasks/{id}`
- `DELETE /v1/tasks/{id}`
- `POST /v1/tasks/{id}/complete` — sets status=done and evaluates pending tasks whose trigger_ids are now all done
- `GET /v1/tasks/events` — Server-Sent Events stream (`Content-Type: text/event-stream`). Emits `task.created`, `task.updated`, `task.deleted`, and `task.status_changed` events in real time after each successful mutation. Each SSE frame uses `id: <task_id>@<version>` for client-side deduplication on reconnect via `Last-Event-ID`.

### B. Admin/Operations API

Purpose: lightweight operational visibility for local deployment.

Implemented endpoints:

- `GET /healthz`
- `GET /v1/adapters` (runtime adapter status)

---

## 2) Planned/Conceptual API Surfaces

The following surfaces describe intended architecture and extension direction. They are **not fully implemented as HTTP routes** in the current server.

### A. Views API

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

### B. Sync API

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

### C. Additional task convenience routes

Conceptual endpoints:

- `POST /v1/tasks/{id}/archive`
- `POST /v1/validate`
- `GET /readyz`
- `GET /v1/meta`
- `GET /v1/version`

---

## 3) Data and Contract Rules

- Canonical task schema: see `schema.md`
- Extension namespace: `ext.*`
- Server is source of truth for:
  - IDs and timestamp invariants
  - lifecycle transitions
  - conflict metadata shape
- API versioning uses explicit route prefix (`/v1/...`)

---

## 4) Client Model

All clients consume the same server contracts:

- `cmd/todoopen` (CLI client; shipped)
- web UI (shipped)
- mobile/TUI and external integrations (planned/partial)

Clients should not mutate JSONL files directly in server mode.

---

## 5) Deployment and Transport

MVP default:

- Local server over loopback HTTP
- Local filesystem JSONL backend

Later:

- Optional hosted/self-hosted remote mode
- authn/authz hardening and policy layers

---

## 6) Internal Mapping (Go)

High-level package responsibility mapping:

- `internal/api` → transport handlers, request/response mapping
- `internal/core` → domain rules and validation
- `internal/store` → persistence contracts and filesystem store
- `internal/sync` → sync orchestration and adapters
- `internal/app` → dependency wiring/composition root

---

## 7) Initial Non-Goals

- Real-time collaboration protocol
- Full enterprise permissions matrix
- Event-sourced/distributed consistency model

These can be layered after stable single-node server semantics are proven.
