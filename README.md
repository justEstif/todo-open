# todo.open

[![Go Version](https://img.shields.io/badge/go-1.26+-00ADD8?logo=go)](https://go.dev)
[![Status](https://img.shields.io/badge/status-MVP-blue)](https://github.com/justEstif/todo-open/blob/main/docs/mvp.md)
[![Storage](https://img.shields.io/badge/storage-JSONL-6f42c1)](https://github.com/justEstif/todo-open/blob/main/docs/schema.md)

**A local task server with an open API — agent-ready, CLI-friendly, your data.**

```sh
todoopen task create --title "Ship the release"
todoopen task list --json
```

```json
{
  "items": [
    { "id": "task_1", "title": "Ship the release", "status": "open", "priority": "normal" }
  ]
}
```

Tasks are stored as plain JSONL on your machine. A local HTTP server exposes them over a stable API. Any tool — CLI, web UI, agent, or script — talks to the same endpoint.

---

## Install

**mise (recommended):**

```sh
mise use -g go:github.com/justEstif/todo-open/cmd/todoopen@v0.1.2
mise reshim
todoopen --version
```

**Build from source:**

```sh
git clone https://github.com/justEstif/todo-open.git
cd todo-open
go build ./cmd/todoopen
./todoopen --version
```

---

## Usage

```sh
# Start the local server + web UI
todoopen web

# Manage tasks
todoopen task create --title "Write release notes"
todoopen task list
todoopen task list --json | jq '.items[] | select(.status == "open")'

# Check active adapters
todoopen adapters
```

| Flag                    | Description                          |
| ----------------------- | ------------------------------------ |
| `--addr 127.0.0.1:8080` | Custom bind address                  |
| `--no-open`             | Start server without opening browser |
| `--server http://...`   | Attach CLI to a running server       |

---

## Adapters

Adapters extend sync and view behavior without touching core task semantics. Install the ones you need:

| Adapter                  | Kind | What it does                        |
| ------------------------ | ---- | ----------------------------------- |
| `todoopen-plugin-sync-git` | sync | Push/pull `tasks.jsonl` to a git repo |
| *(build your own)*       | sync | Rsync, S3, custom backend           |
| *(build your own)*       | view | Markdown, TUI, custom renderer      |

Enable adapters in `.todoopen/config.toml`:

```toml
[views]
  enabled = ["json", "markdown"]

[sync]
  enabled = ["noop", "git"]

[adapters.git]
  bin  = "todoopen-plugin-sync-git"
  kind = "sync"

[adapters.git.config]
  remote = "${GIT_REMOTE}"
  branch = "tasks"

[adapters.markdown]
  bin  = "todoopen-plugin-view-markdown"
  kind = "view"
```

Use `${VAR}` syntax in adapter config values — todo.open expands them from the environment at runtime so secrets never live in the file.

See [docs/adapters.md](docs/adapters.md) for the full plugin protocol and [docs/schema.md](docs/schema.md) for task schema details.

---

## Your data, your tools

Tasks are JSONL. Pipe them anywhere:

```sh
# visidata
todoopen task list --json | vd -f json

# miller
todoopen task list --json | mlr --json filter '$status == "open"'
```

---

## How it works

```
┌─────────────────────────────────────────┐
│  Clients                                │
│  CLI · Web UI · TUI · Agents · Scripts  │
└───────────────┬─────────────────────────┘
                │ loopback HTTP
┌───────────────▼─────────────────────────┐
│  todo.open Server                       │
│  HTTP API → Core domain → JSONL store   │
│                         → Sync layer    │
└───────────────┬─────────────────────────┘
                │ Push / Pull adapter
┌───────────────▼─────────────────────────┐
│  Remote (your choice)                   │
│  git repo · S3 · rsync · custom         │
└─────────────────────────────────────────┘
```

- **Local-first** — runs on your machine, no cloud required
- **Open API** — all clients speak the same loopback HTTP contract
- **Plain JSONL** — tasks are readable, portable, and version-controllable
- **Pluggable sync** — push/pull via any adapter binary you point it at

---

## Roadmap

- [x] Local HTTP API + core domain
- [x] CLI client
- [x] Web UI
- [x] Pluggable sync and view adapter contracts
- [x] Git sync adapter (reference implementation)
- [ ] Agent task coordination / real-time sync
- [ ] TUI client
- [ ] Packaged binaries (`.deb`, `.apk`, `.exe`, `.dmg`)
- [ ] Desktop app

---

## Contributing

```sh
git clone https://github.com/justEstif/todo-open.git
cd todo-open
mise install
mise run build
mise run test
```

Common tasks: `mise run fmt` · `mise run vet` · `mise run test` · `mise run build`
