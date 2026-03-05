---
status: accepted
---

# todo.open Canonical Data Schema (v1)

This document defines the canonical JSONL task record for todo.open MVP.

## 1) Storage layout

- `tasks.jsonl`: one JSON task record per line
- `.todoopen/meta.json`: workspace metadata/config

Each task line in `tasks.jsonl` is a complete current-state object (not an event log).

---

## 2) Canonical task record

### Required fields

- `id` (string): globally unique task id (UUIDv7 preferred)
- `title` (string): short task text, non-empty
- `status` (enum): `"open" | "in_progress" | "done" | "archived"`
- `created_at` (string, RFC3339 UTC timestamp)
- `updated_at` (string, RFC3339 UTC timestamp)

### Optional core fields

- `description` (string)
- `project` (string)
- `tags` (string[]; unique values)
- `priority` (enum): `"low" | "normal" | "high" | "critical"`
- `due_at` (string, RFC3339 UTC timestamp)
- `started_at` (string, RFC3339 UTC timestamp)
- `completed_at` (string, RFC3339 UTC timestamp)
- `deleted_at` (string, RFC3339 UTC timestamp; soft-delete marker)
- `parent_id` (string; parent task id for local hierarchical subtasks)
- `assignee` (string)
- `estimate_minutes` (integer >= 0)
- `sort_order` (number)
- `version` (integer >= 1; increments on each mutation)

### Extension namespace

- `ext` (object): plugin/custom fields only
- All non-core custom fields MUST be nested under `ext`.

Example:

```json
{
  "id": "task_01HZYK6Q3P76Q76Z7X9R6ASJ7F",
  "title": "Prepare weekly review",
  "status": "open",
  "created_at": "2026-03-05T18:00:00Z",
  "updated_at": "2026-03-05T18:00:00Z",
  "tags": ["planning", "weekly"],
  "priority": "normal",
  "ext": {
    "kanban": { "column": "backlog" }
  }
}
```

---

## 3) Status lifecycle and transition rules

Allowed transitions:

- `open -> in_progress | done | archived`
- `in_progress -> open | done | archived`
- `done -> open | archived`
- `archived -> open`

Transition side effects:

- Entering `in_progress` sets `started_at` if absent.
- Entering `done` sets `completed_at` to transition timestamp.
- Leaving `done` clears `completed_at`.
- Entering `archived` keeps `completed_at` unchanged.

Soft delete behavior:

- Setting `deleted_at` marks record deleted without removing line history.
- Deleted records are excluded from default list views.

---

## 4) Timestamp and version invariants

- All timestamps MUST be RFC3339 in UTC (e.g. `2026-03-05T18:00:00Z`).
- `created_at` is immutable after create.
- `updated_at` MUST change on every mutation.
- If present: `started_at >= created_at`, `completed_at >= created_at`, `due_at` unconstrained.
- If `status == "done"`, `completed_at` is required.
- If `status != "done"`, `completed_at` SHOULD be absent.
- `version` starts at `1` and increments by 1 for each accepted update.

---

## 5) Validation profile (MVP)

MVP validation should enforce:

1. Required field presence and types
2. Enum constraints (`status`, `priority`)
3. Timestamp format + UTC normalization
4. Status transition validity against prior state
5. Core field namespace protection (unknown fields must be under `ext`)

Recommended validation modes:

- `strict`: fail on any unknown top-level non-core field
- `compat`: warn and preserve unknown top-level fields (for migration)

---

## 6) Workspace metadata contract (`.todoopen/meta.json`)

```json
{
  "workspace_version": 1,
  "schema_version": "todo.open.task.v1",
  "default_sort": ["status", "priority", "updated_at"],
  "enabled_views": [],
  "enabled_sync_adapters": []
}
```

Required keys:

- `workspace_version` (integer)
- `schema_version` (string)

---

## 7) JSONL examples

`tasks.jsonl` (one task per line):

```json
{"id":"task_01HZYK6Q3P76Q76Z7X9R6ASJ7F","title":"Prepare weekly review","status":"open","created_at":"2026-03-05T18:00:00Z","updated_at":"2026-03-05T18:00:00Z","tags":["planning"],"priority":"normal"}
{"id":"task_01HZYK83Q9WJ5CG1ACJH5HMCJZ","title":"Ship MVP schema docs","status":"done","created_at":"2026-03-05T16:00:00Z","updated_at":"2026-03-05T18:10:00Z","completed_at":"2026-03-05T18:10:00Z","priority":"high"}
```

---

## 8) Implementation hooks

- Core API `createTask` must initialize: `id`, `status=open`, `created_at`, `updated_at`, `version=1`.
- Core API `updateTask` must enforce transition + invariants and bump `version`.
- `completeTask` is a specialized transition helper to `status=done` + `completed_at`.
- `validateTask`/`validateStore` should load `schema_version` from `.todoopen/meta.json`.
