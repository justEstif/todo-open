---
status: accepted
---

# Human vs. Agent Interaction Model and Boundaries

todo-open is designed for both humans and agents to operate on the **same task list** with clear boundaries and safety guarantees. This document explains the mental model, responsibilities, and constraints for successful human-agent collaboration.

## 1. Mental Model

todo-open treats tasks as a shared resource where:

- **Humans own task creation, prioritization, and final judgment**
  - Humans define what work needs to be done
  - Humans set priorities and due dates
  - Humans have the final say on scope and cancellation

- **Agents execute, claim work, and report outcomes**
  - Agents autonomously pick up available work
  - Agents execute tasks within their capabilities
  - Agents report success, failure, or progress
  - Agents **do not** redefine scope or make business decisions

This creates a clear separation of concerns: humans determine **what** and **why**, agents determine **how** and **when**.

## 2. Responsibility Boundaries

| Operation | Human | Agent | Both |
|-----------|-------|-------|------|
| Create task | ✓ | ✓ (via PUT with ID) | |
| Set priority | ✓ | | |
| Claim task | | ✓ | |
| Heartbeat | | ✓ | |
| Complete task | ✓ (PATCH status=done) | ✓ (/complete) | |
| Cancel/delete | ✓ | | |
| Set trigger_ids | ✓ | ✓ | |

### Key Insights:

- **Task Creation**: Both can create, but humans typically create via `POST` (server-generated ID) while agents use `PUT` with specific IDs for reproducible workflows.
- **Priority Setting**: Only humans set priorities - this ensures business value alignment remains human-controlled.
- **Work Execution**: Only agents claim and work on tasks. This prevents race conditions between humans and agents.
- **Completion**: Both can mark tasks done, but through different endpoints reflecting their roles.
- **Scope Control**: Only humans can delete tasks - agents release failed work but never remove it from the system.

## 3. Safety Constraints

### ETag/If-Match on Mutations

**Agents MUST use ETag/If-Match on all mutations** (except claim/heartbeat/release/complete which have built-in concurrency control).

**Why**: Ensures optimistic concurrency and prevents silent overwrites of concurrent changes. Without ETags, two agents could read the same task version, make different changes, and the last writer would silently overwrite the first.

```bash
# WRONG: No concurrency protection
PATCH /v1/tasks/{id}
{"title": "updated by agent"}

# RIGHT: With ETag protection
GET /v1/tasks/{id}  # Returns ETag: "5"
PATCH /v1/tasks/{id}
If-Match: "5"
{"title": "updated by agent"}
```

### No Task Deletion by Agents

**Agents MUST NOT delete tasks.** Humans decide scope - agents can only release work back to the pool via `/release`.

**Why**: Deletion is a business decision about what work matters. Agents should not have the power to remove work from the system.

### Idempotency Keys

**Agents SHOULD use `X-Idempotency-Key` on all POST mutations.**

**Why**: Prevents duplicate work when network retries occur. Without idempotency keys, a temporary network error could cause an agent to claim the same task twice.

```bash
POST /v1/tasks/{id}/claim
X-Idempotency-Key: agent-worker-1-task_01HZYK-claim-attempt-1
{"agent_id": "worker-1", "lease_ttl_seconds": 300}
```

### Heartbeat Timing

**Agents SHOULD heartbeat at < 50% of their lease TTL.**

**Why**: Provides buffer time for network delays and prevents lease expiry due to temporary unavailability. For a 300s lease, heartbeat every 150s ensures the agent can miss 2-3 heartbeats before losing the lease.

### Manual File Editing

**Humans SHOULD NOT manually edit `tasks.jsonl` while the server is running.**

**Why**: The server maintains in-memory state and validation. Direct file edits can corrupt the data model, cause version conflicts, and lead to data loss. All mutations should go through the API.

## 4. Mixed Workflow Example

Let's walk through a realistic scenario: **Human creates task with dependencies, two agents race to pick up work, one succeeds, heartbeats, completes, dependency auto-resolves, second agent picks up the unblocked task.**

### Step 1: Human creates dependent tasks

```bash
# Human creates task A (prerequisite)
curl -X POST http://localhost:8080/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Set up database schema",
    "priority": "high"
  }'

# Response: task_A_id = "task_01HZYK6Q3P76Q76Z7X9R6ASJ7F"

# Human creates task B (depends on A)
curl -X POST http://localhost:8080/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Migrate data to new schema",
    "priority": "high",
    "trigger_ids": ["task_01HZYK6Q3P76Q76Z7X9R6ASJ7F"]
  }'

# Response: task_B_id = "task_01HZYK83Q9WJ5CG1ACJH5HMCJZ"
# Status: "pending" (waiting for A to complete)
```

### Step 2: Two agents seek work

**Agent 1:**
```bash
# Agent 1 checks for work
curl -X GET http://localhost:8080/v1/tasks/next

# Response: task_A (available to claim)
{
  "id": "task_01HZYK6Q3P76Q76Z7X9R6ASJ7F",
  "title": "Set up database schema",
  "status": "open",
  "priority": "high"
}
```

**Agent 2:**
```bash
# Agent 2 simultaneously checks for work (race condition)
curl -X GET http://localhost:8080/v1/tasks/next

# Response: task_A (same task, race begins)
```

### Step 3: Race to claim - Agent 1 wins

**Agent 1:**
```bash
# Agent 1 attempts claim
curl -X POST http://localhost:8080/v1/tasks/task_01HZYK6Q3P76Q76Z7X9R6ASJ7F/claim \
  -H "Content-Type: application/json" \
  -H "X-Idempotency-Key: agent-1-claim-001" \
  -d '{
    "agent_id": "agent-1",
    "lease_ttl_seconds": 300
  }'

# Response: 200 OK - Claim successful!
{
  "id": "task_01HZYK6Q3P76Q76Z7X9R6ASJ7F",
  "status": "in_progress",
  "ext": {
    "agent": {
      "id": "agent-1",
      "claimed_at": "2026-03-08T20:35:00Z",
      "lease_expires_at": "2026-03-08T20:40:00Z",
      "heartbeat_at": "2026-03-08T20:35:00Z",
      "lease_ttl_seconds": 300
    }
  }
}
```

**Agent 2:**
```bash
# Agent 2 attempts claim (loses race)
curl -X POST http://localhost:8080/v1/tasks/task_01HZYK6Q3P76Q76Z7X9R6ASJ7F/claim \
  -H "Content-Type: application/json" \
  -H "X-Idempotency-Key: agent-2-claim-001" \
  -d '{
    "agent_id": "agent-2",
    "lease_ttl_seconds": 300
  }'

# Response: 409 Conflict - Already claimed
{
  "error": {
    "code": "conflict",
    "message": "task already claimed by agent agent-1"
  }
}
```

### Step 4: Agent 1 works and heartbeats

```bash
# Agent 1 does work... then heartbeats after 75s (50% of 150s TTL)
curl -X POST http://localhost:8080/v1/tasks/task_01HZYK6Q3P76Q76Z7X9R6ASJ7F/heartbeat \
  -H "Content-Type: application/json" \
  -d '{
    "agent_id": "agent-1"
  }'

# Response: 200 OK - Lease extended
{
  "ext": {
    "agent": {
      "id": "agent-1",
      "lease_expires_at": "2026-03-08T20:42:30Z"  # Extended
    }
  }
}
```

### Step 5: Agent 1 completes task

```bash
# Agent 1 completes the work
curl -X POST http://localhost:8080/v1/tasks/task_01HZYK6Q3P76Q76Z7X9R6ASJ7F/complete

# Response: 200 OK - Task done
{
  "id": "task_01HZYK6Q3P76Q76Z7X9R6ASJ7F",
  "status": "done",
  "completed_at": "2026-03-08T20:37:00Z"
}
```

### Step 6: Dependency auto-resolves

**Server action:** The server automatically evaluates all `pending` tasks whose `trigger_ids` are now all `done`. Task B's `trigger_ids` contained task A, so:

```bash
# Server emits SSE event (or agents can poll)
# Event: task.status_changed for task B
# New status: "open" (was "pending")
```

### Step 7: Agent 2 picks up unblocked task

```bash
# Agent 2 checks for work again
curl -X GET http://localhost:8080/v1/tasks/next

# Response: Now task B is available!
{
  "id": "task_01HZYK83Q9WJ5CG1ACJH5HMCJZ",
  "title": "Migrate data to new schema", 
  "status": "open",
  "priority": "high"
}

# Agent 2 claims and completes task B...
```

## 5. What Lives Where

### `tasks.jsonl`: Task Data + Outcomes
- **Purpose**: Canonical task records, human-readable, versionable
- **Contents**: One JSON task object per line
- **Format**: See [schema.md](./schema.md) for full specification
- **Access**: Server-managed, read/write through API only
- **Human role**: Define task structure, priorities, dependencies
- **Agent role**: Read via API, never edit directly

### `ext.agent.*`: Ephemeral Lease State
- **Purpose**: Machine bookkeeping for agent coordination
- **Contents**: `agent_id`, `claimed_at`, `lease_expires_at`, `heartbeat_at`, `lease_ttl_seconds`
- **Lifetime**: Cleared on task completion, release, or lease expiry
- **Access**: Server-managed via agent coordination endpoints
- **Human role**: Ignore (implementation detail)
- **Agent role**: Manage via `/claim`, `/heartbeat`, `/release`

### `config.toml`: Adapter Configuration
- **Purpose**: Runtime configuration for adapters and plugins
- **Contents**: Plugin registrations, adapter settings
- **Format**: TOML configuration file
- **Access**: Edit by humans, read by server at startup
- **Human role**: Configure adapters and views
- **Agent role**: None (configuration is human-controlled)

### `meta.json`: Workspace Version Metadata
- **Location**: `.todoopen/meta.json`
- **Purpose**: Workspace-level configuration and schema versioning
- **Contents**: `workspace_version`, `schema_version`, enabled adapters/views
- **Format**: JSON configuration file
- **Access**: Server-managed, human-editable for workspace settings
- **Human role**: Set workspace preferences
- **Agent role**: None

## 6. Cross-References

- **[API Documentation](./api.md)**: Full HTTP API reference including all endpoints used by both humans and agents
- **[Agent Primitives](./agent-primitives.md)**: Detailed agent contract with all endpoints, error codes, and best practices
- **[Schema Documentation](./schema.md)**: Canonical task record format, status transitions, and validation rules
- **[Architecture Overview](./architecture.md)**: High-level system design and component responsibilities

## 7. Best Practices

### For Humans
1. **Use the CLI or web UI** for task creation and management
2. **Set clear priorities** to guide agent work selection
3. **Use dependencies** (`trigger_ids`) for complex workflows
4. **Monitor agent activity** through task status changes
5. **Never edit `tasks.jsonl` directly** while the server runs

### For Agent Developers
1. **Always use ETag/If-Match** on PATCH operations
2. **Include idempotency keys** on all POST operations
3. **Heartbeat regularly** (at 50% of lease TTL)
4. **Handle all error codes** (see agent-primitives.md §5)
5. **Use SSE events** for real-time updates when possible
6. **Gracefully release** tasks on shutdown/errors

### For Mixed Environments
1. **Coordinate access patterns** - humans during business hours, agents overnight
2. **Use priority tiers** to ensure critical human tasks aren't blocked by agent work
3. **Monitor conflict rates** and adjust agent behavior if too many races occur
4. **Establish review processes** for agent-completed work before human sign-off