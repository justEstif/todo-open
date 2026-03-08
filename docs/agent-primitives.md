---
status: accepted
---

# Agent Primitives Contract

This document specifies the **stable, machine-oriented API surface** for software agents interacting with the todo.open server.

Agents are autonomous programs that pick up tasks, perform work, and report completion. This contract covers the full lifecycle, all relevant endpoints with schemas, concurrency semantics, event streaming, idempotency, and dependency graphs.

For the general HTTP API reference see `docs/api.md`. This document is designed to be self-contained for agent implementors.

---

## 1. Task Lifecycle for Agents

```
open ──► claim ──► in_progress ──► complete ──► done
          │            │
          │            ├──► heartbeat (extend lease, loop)
          │            │
          │            └──► release ──► open
          │
          └── (expired lease swept by server) ──► open
```

States visible to agents:

| Status       | Meaning                                              |
|--------------|------------------------------------------------------|
| `pending`    | Waiting for dependencies; not yet claimable          |
| `open`       | Available to claim                                   |
| `in_progress`| Claimed by an agent with an active lease             |
| `done`       | Work complete; may unblock pending tasks             |
| `archived`   | Terminal; excluded from work queues                  |

**Only `open` tasks are claimable.** An agent must not attempt to claim a task unless its status is `open`.

---

## 2. All Agent Endpoints

### 2.1 Work Queue

#### `GET /v1/tasks/next`

Returns the single highest-priority unclaimed `open` task.

Priority order (descending): `critical` > `high` > `normal` > `low`.

**Response 200 OK** — task object (see §3 for schema).

**Response 404 Not Found**
```json
{ "error": { "code": "not_found", "message": "task not found" } }
```

**Usage pattern:**
```
loop:
  GET /v1/tasks/next
  if 404 → sleep and retry
  if 200 → task_id = body.id; proceed to claim
```

---

### 2.2 Claim

#### `POST /v1/tasks/{id}/claim`

Atomically claim a task for exclusive work. Transitions `open → in_progress`.

**Request body:**
```json
{
  "agent_id": "agent-xyz",
  "lease_ttl_seconds": 300
}
```

| Field               | Type    | Required | Notes                                        |
|---------------------|---------|----------|----------------------------------------------|
| `agent_id`          | string  | yes      | Stable identifier for this agent instance    |
| `lease_ttl_seconds` | integer | no       | Default 300. How long until lease expires.   |

**Response 200 OK** — updated task object with `ext.agent` populated:
```json
{
  "id": "task_01...",
  "status": "in_progress",
  "ext": {
    "agent": {
      "id": "agent-xyz",
      "claimed_at": "2026-03-08T12:00:00Z",
      "lease_expires_at": "2026-03-08T12:05:00Z",
      "heartbeat_at": "2026-03-08T12:00:00Z",
      "lease_ttl_seconds": 300
    }
  }
}
```

**Response 409 Conflict** — task already claimed by another agent with an active lease, or task is not `open`.
```json
{ "error": { "code": "conflict", "message": "task already claimed by agent ..." } }
```

**Response 404 Not Found** — task does not exist.

---

### 2.3 Heartbeat

#### `POST /v1/tasks/{id}/heartbeat`

Extend the lease to prevent expiry while work is in progress. Must be called before `lease_expires_at` elapses.

**Request body:**
```json
{ "agent_id": "agent-xyz" }
```

**Response 200 OK** — updated task with extended `ext.agent.lease_expires_at`.

**Response 403 Forbidden** — `agent_id` does not match the current claim holder.
```json
{ "error": { "code": "forbidden", "message": "agent_id mismatch" } }
```

**Recommended heartbeat interval:** `lease_ttl_seconds / 2`. Example: TTL=300s → heartbeat every 150s.

---

### 2.4 Complete

#### `POST /v1/tasks/{id}/complete`

Mark a task as done. Transitions `in_progress → done` and clears `ext.agent`. Also evaluates any `pending` tasks whose `trigger_ids` are now all `done` and unblocks them to `open`.

**Request body:** none required.

**Response 200 OK** — completed task with `status: "done"` and `completed_at` set.

**Response 404 Not Found** — task does not exist.

**Note:** `/complete` does not check `agent_id`. Agents should call `/complete` only for tasks they hold the lease for. If the lease has expired and the task was swept back to `open`, `/complete` may still succeed but transitions from the current state.

---

### 2.5 Release

#### `POST /v1/tasks/{id}/release`

Voluntarily give up a claim. Transitions `in_progress → open`, clears `ext.agent`. Use when the agent cannot complete the task (e.g. unrecoverable error, shutdown).

**Request body:**
```json
{ "agent_id": "agent-xyz" }
```

**Response 200 OK** — updated task with `status: "open"`.

**Response 403 Forbidden** — `agent_id` does not match the claim holder.

---

### 2.6 List (for polling)

#### `GET /v1/tasks?status=open`

Returns all tasks matching the given status filter.

**Query parameters:**

| Param        | Values        | Notes                                    |
|--------------|---------------|------------------------------------------|
| `status`     | any TaskStatus| Filter by status                         |
| `is_blocked` | `true`        | Return only tasks with non-empty blocked_by |

**Response 200 OK:**
```json
{ "items": [ /* task objects */ ] }
```

---

### 2.7 Get Single Task

#### `GET /v1/tasks/{id}`

**Response 200 OK** — task object with `ETag: "<version>"` header.

**Response 404 Not Found** — task does not exist.

---

### 2.8 Update Task (PATCH)

#### `PATCH /v1/tasks/{id}`

Update task fields (title, status, etc.).

**Request body:**
```json
{ "title": "new title" }
```

Agents may use this to update task metadata. For status transitions use the dedicated endpoints when possible.

Supports `If-Match` header for optimistic concurrency (see §4).

**Response 200 OK** — updated task.

**Response 409 Conflict** — `If-Match` version mismatch.

---

### 2.9 SSE Event Stream

#### `GET /v1/tasks/events`

See §6 for full SSE usage pattern.

---

## 3. Task Object Schema

```json
{
  "id":           "task_01HZYK...",
  "title":        "Do something",
  "status":       "open",
  "priority":     "normal",
  "created_at":   "2026-03-08T12:00:00Z",
  "updated_at":   "2026-03-08T12:00:00Z",
  "started_at":   null,
  "completed_at": null,
  "version":      1,
  "description":  "",
  "project":      "",
  "parent_id":    "",
  "assignee":     "",
  "tags":         [],
  "trigger_ids":  [],
  "blocking":     [],
  "blocked_by":   [],
  "ext":          {}
}
```

All timestamps are RFC3339 in UTC.

---

## 4. ETag Semantics and If-Match Usage

Every mutating response includes an `ETag` header:

```
ETag: "5"
```

The value is the integer `version` of the task, wrapped in quotes.

### Optimistic concurrency with If-Match

To update a task only if it hasn't changed since you last read it:

```
GET /v1/tasks/{id}
→  ETag: "5"

PATCH /v1/tasks/{id}
→  If-Match: "5"
   Body: { ... }
```

If the stored version is still `5`, the update proceeds and the response carries `ETag: "6"`.

If the version has changed (someone else updated it), the server returns:

```
HTTP 409 Conflict
{ "error": { "code": "version_conflict", "message": "ETag mismatch; resource was modified" } }
```

### When to use If-Match

- **Always** use `If-Match` on `PATCH` when you read the task before writing, to avoid overwriting concurrent changes.
- The claim/heartbeat/release/complete endpoints manage concurrency internally — do not add `If-Match` to those.

---

## 5. Error Codes Agents Must Handle

| HTTP Status | Error Code          | Meaning and action                                                 |
|-------------|---------------------|--------------------------------------------------------------------|
| 400         | `invalid_json`      | Malformed request body. Fix and retry.                             |
| 400         | `validation_error`  | Schema/field violation. Fix and retry.                             |
| 403         | `forbidden`         | `agent_id` does not own the claim. Do not retry with same agent_id.|
| 404         | `not_found`         | Task does not exist or queue is empty. Retry after delay.          |
| 409         | `conflict`          | Task already claimed or version mismatch. Re-read and re-evaluate. |
| 409         | `version_conflict`  | `If-Match` mismatch. Re-read (GET) and retry with new ETag.        |
| 412         | (reserved)          | Not currently used; treat as conflict and re-read.                 |
| 500         | `internal_error`    | Server fault. Retry with exponential backoff.                      |

**409 on claim**: The task was claimed by another agent, or its status changed. Fetch `/v1/tasks/next` again — a different task may be available.

**403 on heartbeat/release**: Your agent_id does not match the claim. This can happen if the lease expired and the server swept the task back to `open`, and a different agent claimed it. Abandon this work unit.

---

## 6. SSE Event Stream Usage

### Connecting

```
GET /v1/tasks/events
Accept: text/event-stream
```

The server sends events continuously. Keep the connection open.

### Event format

```
id: task_01HZYK@5
event: task.updated
data: {"id":"task_01HZYK...","status":"in_progress",...}

```

Each SSE frame:
- `id`: `<task_id>@<version>` — use for deduplication and reconnect
- `event`: one of `task.created`, `task.updated`, `task.deleted`, `task.status_changed`
- `data`: full task JSON object

### Reconnecting via Last-Event-ID

On disconnect, reconnect with:

```
GET /v1/tasks/events
Last-Event-ID: task_01HZYK@5
```

The server will resume from that point (events after that version for that task). This prevents duplicate processing of events already seen.

### Agent usage pattern

```
connect to /v1/tasks/events
loop:
  receive event
  if event.event == "task.status_changed":
    if event.data.status == "open" and want_work:
      attempt claim of event.data.id
  if event.event == "task.created":
    if event.data.status == "open":
      attempt claim if this agent should pick it up
  on disconnect: reconnect with Last-Event-ID from last seen frame
```

Agents may combine SSE listening with periodic polling of `/v1/tasks/next` as a fallback.

---

## 7. Idempotency Guidance

### Safe to retry (idempotent)

| Operation                   | Safe? | Notes                                                         |
|-----------------------------|-------|---------------------------------------------------------------|
| `GET /v1/tasks/next`        | ✅    | Read-only.                                                    |
| `GET /v1/tasks/{id}`        | ✅    | Read-only.                                                    |
| `GET /v1/tasks`             | ✅    | Read-only.                                                    |
| `POST /v1/tasks/{id}/heartbeat` | ✅ | Re-extending an already-extended lease is harmless.          |
| `PATCH /v1/tasks/{id}` with `If-Match` | ✅ | Safe: 409 on conflict, not a silent overwrite.       |
| `POST /v1/tasks/{id}/complete` | ⚠️  | Calling twice: first call succeeds; second returns 400 (already done). |
| `POST /v1/tasks/{id}/claim`   | ⚠️  | Second call by same agent on already-claimed task returns 409. |
| `POST /v1/tasks/{id}/release` | ⚠️  | Second call returns 403 (no claim to release).               |
| `POST /v1/tasks` (create)     | ❌   | Each call creates a new task. Use `X-Idempotency-Key` (see §8). |
| `PATCH /v1/tasks/{id}` without `If-Match` | ❌ | Blind update; each call applies changes.          |

### Retry strategy

Use **exponential backoff with jitter** for 500 errors. For 409 (conflict) on claim, re-read `/v1/tasks/next` immediately without backoff (another task may be available). For 403, log and abort (do not retry with same agent_id).

---

## 8. Idempotency Key for Agent POST Endpoints

Agents may include `X-Idempotency-Key: <unique-key>` on `POST` endpoints (claim, complete, heartbeat, release, create) to deduplicate retried requests.

**Behavior:**

- If the server has seen this key within the last 5 minutes, it returns the **cached response** from the original request (including status code and body).
- If the key is new, the request is processed normally and the response is cached.
- The cache is **in-memory only** — it does not survive server restarts. Agents must not rely on idempotency keys surviving a server restart.

**Key guidelines:**

- Use a unique key per logical operation attempt, e.g. `<agent_id>-<task_id>-<attempt_id>`.
- Keys are scoped server-side; reuse across different operations is safe (but confusing — avoid it).
- Example: `X-Idempotency-Key: agent-worker-1-task_01HZYK-claim-attempt-1`

---

## 9. Dependency Graph: trigger_ids

### Creating a dependent task

To create task B that should only become `open` after task A is `done`:

```json
POST /v1/tasks
{
  "title": "Task B (depends on A)",
  "trigger_ids": ["task_A_id"]
}
```

Task B is created with `status: "pending"`.

### What signals readiness

When `POST /v1/tasks/{A_id}/complete` is called:

1. Task A transitions to `done`.
2. The server evaluates all `pending` tasks whose `trigger_ids` are a subset of `done` tasks.
3. If all of B's `trigger_ids` are now `done`, B transitions `pending → open`.
4. A `task.status_changed` SSE event is emitted for B.

Agents listening to the SSE stream will receive the `task.status_changed` event for B and can immediately attempt to claim it.

### Polling alternative

If not using SSE, poll `GET /v1/tasks?status=open` periodically after completing tasks that have dependents.

### Multiple dependencies

```json
{
  "title": "Task C (depends on A and B)",
  "trigger_ids": ["task_A_id", "task_B_id"]
}
```

Task C remains `pending` until **both** A and B are `done`.

### blocking / blocked_by fields

These are informational edge annotations:

- `blocking`: task IDs that this task is blocking (outgoing edges).
- `blocked_by`: task IDs blocking this task (incoming edges).

These are distinct from `trigger_ids`. `trigger_ids` is the authoritative dependency list for automatic `pending → open` transitions. `blocking`/`blocked_by` are annotations for human display and filtering (`GET /v1/tasks?is_blocked=true`).

---

## 10. Full Agent Loop Example

```
agent_id = "worker-1"
lease_ttl = 300  # seconds

loop:
  # 1. Find work
  resp = GET /v1/tasks/next
  if resp.status == 404:
    sleep(10)
    continue
  task_id = resp.body.id

  # 2. Claim it
  resp = POST /v1/tasks/{task_id}/claim
          body: { agent_id, lease_ttl_seconds: lease_ttl }
  if resp.status == 409:
    continue  # someone else got it; try next
  if resp.status != 200:
    backoff_and_retry()

  # 3. Work loop with heartbeat
  heartbeat_interval = lease_ttl / 2
  last_heartbeat = now()
  while not work_done():
    do_work_chunk()
    if now() - last_heartbeat > heartbeat_interval:
      POST /v1/tasks/{task_id}/heartbeat  body: { agent_id }
      last_heartbeat = now()

  # 4. Complete
  POST /v1/tasks/{task_id}/complete

  # On error: release
  on_error:
    POST /v1/tasks/{task_id}/release  body: { agent_id }
```

---

## 11. Lease Sweeper

A background goroutine on the server runs every **30 seconds**. It scans all `in_progress` tasks and transitions any whose `ext.agent.lease_expires_at` is in the past back to `open`, clearing `ext.agent`. Each expiry is logged server-side.

**Agents should not rely on the sweeper as a normal code path.** Always call `/release` on graceful shutdown and maintain heartbeats during work. The sweeper is a safety net for crashed agents.
