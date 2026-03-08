# todo.open

**A local-first task server with an open API. Your tasks live on disk. Any tool — CLI, agent, or script — talks to the same endpoint.**

```sh
# human creates a task
todoopen task create --title "Refactor auth module" --priority high
# Created task_a1b2: Refactor auth module

# agent picks it up
curl -s localhost:8080/v1/tasks/next
curl -s -X POST localhost:8080/v1/tasks/task_a1b2/claim -H 'X-Idempotency-Key: run-42'

# human checks progress
todoopen task list
# ID          TITLE                  STATUS       HELD BY
# task_a1b2   Refactor auth module   in_progress  agent-claude (expires 4m)
```

Tasks are stored as plain `tasks.jsonl` — one JSON object per line. A local HTTP server exposes them over a stable REST + SSE API. Agents get a lease system and idempotency keys; humans always stay in control.

→ **[todo-open.pages.dev](https://justestif.github.io/todo-open)** for full docs and install instructions.

---

## Install

```sh
# npm
npm install -g @justestif/todo-open

# mise
mise use -g go:github.com/justEstif/todo-open/cmd/todoopen@latest && mise reshim

# source
git clone https://github.com/justEstif/todo-open.git
cd todo-open && go build ./cmd/todoopen ./cmd/todoopen-server
```

Start the server and open the web UI:

```sh
todoopen web
# Server listening on http://127.0.0.1:8080
```

Or use the terminal UI (starts the server automatically, just like `web`):

```sh
todoopen tui
# Server listening · TUI connected

# point at an existing server
todoopen tui --server http://127.0.0.1:8080
```

---

## How it works

```
.todoopen/tasks.jsonl      ← plain text, yours forever
         │
   todoopen server         ← local HTTP :8080
   (REST + SSE)
         │
   ┌─────┼──────────────┬──────────┐
  CLI   web UI        TUI      AI agents
       (live SSE)  (live SSE)  /v1/tasks/next
```

- **Plain files** — `tasks.jsonl`, readable in any editor, no database
- **Open API** — full REST + Server-Sent Events; `curl` is a valid client
- **Agent-safe** — lease system, heartbeats, idempotency keys; human can always override
- **Live UI** — web UI and terminal UI both receive real-time task updates over SSE
- **Terminal UI** — `todoopen-tui` is a Bubble Tea client: list, create, edit, navigate deps, all from the terminal

---

## Agent API

```sh
# get the full self-describing agent contract (no server needed)
todoopen --agent-info
todoopen --agent-info --server http://host:9090   # custom server URL
todoopen -A                                        # short form

# discover the contract from a running server
curl -s localhost:8080/v1/capabilities | jq .agent

# next unclaimed task → claim → heartbeat → complete
curl -s localhost:8080/v1/tasks/next
curl -s -X POST localhost:8080/v1/tasks/{id}/claim -H 'X-Idempotency-Key: run-42'
curl -s -X POST localhost:8080/v1/tasks/{id}/heartbeat
curl -s -X POST localhost:8080/v1/tasks/{id}/complete

# watch live events
curl -s localhost:8080/v1/tasks/events
```

---

## Adapters

Adapters are separate binaries that extend sync and view behavior. Install only what you need.

| Adapter | Kind | What it does |
|---|---|---|
| `todoopen-adapter-sync-git` | sync | Push/pull `tasks.jsonl` to a git repo |
| `todoopen-adapter-sync-s3` | sync | Sync workspace to S3 |
| `todoopen-adapter-view-markdown` | view | Render tasks as `TASKS.md` |
| *build your own* | sync/view | Any language, any backend |

Configure adapters in `.todoopen/config.toml`:

```toml
[adapters.git]
  bin = "todoopen-adapter-sync-git"

[adapters.git.config]
  remote = "${GIT_REMOTE}"   # env vars expanded at runtime
  branch = "main"
```

See [docs/adapters.md](docs/adapters.md) to build your own.

---

## Documentation

| Doc | What it covers |
|---|---|
| [docs/api.md](docs/api.md) | Full REST + SSE API reference |
| [docs/adapters.md](docs/adapters.md) | Adapter protocol and how to build your own |
| [docs/schema.md](docs/schema.md) | Task schema, JSONL format, field definitions |
| [docs/agent-primitives.md](docs/agent-primitives.md) | Agent coordination contract |
| [docs/human-ux-invariants.md](docs/human-ux-invariants.md) | Human-first UX rules and CLI contract |
| [docs/architecture.md](docs/architecture.md) | Internal design and package layout |

---

MIT License
