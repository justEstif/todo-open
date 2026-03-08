---
name: todoopen-tasks
description: Task management using todo.open, a local-first task server with an open REST + SSE API. Tasks live on disk as tasks.jsonl; any tool — CLI, agent, or script — talks to the same endpoint. Use this skill whenever you need to create, list, claim, complete, or manage tasks in a todo.open workspace. Triggers on "create a task", "list tasks", "claim the next task", "mark task done", "what should I work on next", "check task status", or any request to interact with a running todoopen server.
---

# todo.open Skill

`todoopen` is a local-first task server. Tasks are stored as `tasks.jsonl` on disk and exposed over a stable REST + SSE API. Any tool — CLI, agent, or `curl` — is a valid client.

## Check availability

```bash
todoopen --agent-info   # prints tool description, all endpoints, schema, and workflow as JSON
```

If `todoopen` is not found, install it:

```bash
# npm
npm install -g @justestif/todo-open

# mise
mise use -g go:github.com/justEstif/todo-open/cmd/todoopen@latest && mise reshim

# source
git clone https://github.com/justEstif/todo-open.git
cd todo-open && go build ./cmd/todoopen ./cmd/todoopen-server
```

## Start the server

```bash
todoopen web                           # start server + open web UI (default :8080)
todoopen web --no-open                 # start server without opening browser
todoopen tui                           # start server + terminal UI
todoopen web --server http://host:9090 # connect to an existing server instead
```

## Task CRUD (CLI)

```bash
# Create
todoopen task create --title "Implement login flow" --priority high

# List
todoopen task list

# Get a specific task
todoopen task get --id task_a1b2

# Update title
todoopen task update --id task_a1b2 --title "Implement OAuth login flow"

# Delete
todoopen task delete --id task_a1b2
```

## Agent workflow (REST)

The recommended loop — fetch, claim, heartbeat, complete:

```bash
# 1. Fetch the highest-priority open task
curl -s localhost:8080/v1/tasks/next

# 2. Claim it (lease = 5 min; idempotent with X-Idempotency-Key)
curl -s -X POST localhost:8080/v1/tasks/{id}/claim \
  -H 'X-Idempotency-Key: run-42' \
  -H 'Content-Type: application/json' \
  -d '{"assignee":"agent-claude"}'

# 3. Heartbeat every ~60 s to keep the lease alive
curl -s -X POST localhost:8080/v1/tasks/{id}/heartbeat

# 4a. Complete when done
curl -s -X POST localhost:8080/v1/tasks/{id}/complete

# 4b. Release without completing (error / abort)
curl -s -X POST localhost:8080/v1/tasks/{id}/release
```

## Watch live events (SSE)

```bash
curl -s localhost:8080/v1/tasks/events
# event: task.status_changed
# data: {"id":"task_a1b2","status":"in_progress","assignee":"agent-claude"}
```

## Discover the server contract

```bash
# Full self-describing payload — no server required
todoopen --agent-info | jq .

# Runtime capabilities from the server
curl -s localhost:8080/v1/capabilities | jq .agent
```

## Schema reference

**Statuses** (task lifecycle):
| Value | Meaning |
|---|---|
| `pending` | Waiting for `trigger_ids` dependencies to complete |
| `open` | Ready to be claimed |
| `in_progress` | Claimed; a lease is active |
| `done` | Completed successfully |
| `archived` | Removed from active queues; retained in history |

**Priorities** (claim ordering):
| Value | Meaning |
|---|---|
| `critical` | Urgent; address immediately |
| `high` | Prefer over normal tasks |
| `normal` | Default |
| `low` | Work when nothing else is queued |

## Adapter sync (optional)

```bash
# Check configured adapter status
todoopen adapters

# Example .todoopen/config.toml for S3 backup
[adapters.s3]
  bin = "todoopen-adapter-sync-s3"

[adapters.s3.config]
  bucket   = "${S3_BUCKET}"
  region   = "us-east-1"
  access_key = "${S3_ACCESS_KEY}"
  secret_key = "${S3_SECRET_KEY}"
```

## Tips

- **Agent self-config** — call `todoopen --agent-info` once at startup; parse `endpoints` and `workflow` fields to configure your client without hardcoding URLs
- **Idempotency** — always pass `X-Idempotency-Key` on `claim` and other mutating requests; replaying the same key returns the original response without side-effects
- **Lease expiry** — if no heartbeat arrives within 300 s the task reverts to `open` automatically; no manual cleanup needed
- **Human override** — humans can always release an agent's claim; the API never blocks human intervention
- **Plain files** — tasks.jsonl is readable in any editor; the server is optional for reading, required for writes
- **SSE** — subscribe to `/v1/tasks/events` for real-time updates instead of polling `/v1/tasks`
