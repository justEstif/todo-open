# todo.open

![Go Version](https://img.shields.io/badge/go-1.26+-00ADD8?logo=go)
![Status](https://img.shields.io/badge/status-MVP-blue)
![Storage](https://img.shields.io/badge/storage-JSONL-6f42c1)

A **server-first, local-first** task tool built in Go.

- ✅ Portable local data (JSONL)
- ✅ Clear server-side domain rules
- ✅ CLI + web flow today, more clients later

## TL;DR — why this exists

todo.open exists to avoid the usual tradeoff: **either convenient apps with locked-in data**, or **portable plain-text systems with weak UX**.

The project aims to give you both:
- local, user-owned task data (`tasks.jsonl`)
- one canonical server contract for all clients
- extensibility (views/sync adapters) without changing core task semantics

For the full rationale and MVP scope, see **[docs/mvp.md](docs/mvp.md)**.

---

## 30-second demo

```bash
# 1) Start local web app (starts server + opens browser)
todoopen web

# 2) In another terminal, create a task via CLI
todoopen task create -title "Ship README polish"

# 3) List tasks
todoopen task list
```

If your server runs on a custom address, add `-server http://127.0.0.1:8080` to task commands.

---

## Install

### Option A: mise (recommended)

Install globally with the mise Go backend:

```bash
mise use -g go:github.com/justEstif/todo-open/cmd/todoopen@latest
todoopen --help
```

Project-local install (no global binary):

```bash
mise use go:github.com/justEstif/todo-open/cmd/todoopen@latest
mise x -- todoopen --help
```

### Option B: build from source with Go

Prerequisite: **Go 1.26+**

```bash
git clone https://github.com/justEstif/todo-open.git
cd todo-open
go build ./cmd/todoopen
./todoopen --help
```

---

## Quick usage

```bash
todoopen --help
todoopen web
todoopen task create -title "My first task"
todoopen task list
```

Useful web flags:

- `--addr 127.0.0.1:8080` set local bind address
- `--no-open` start without opening browser
- `--server http://127.0.0.1:8080` attach to existing server

Alias:

```bash
todoopen gui
```

Run server directly from source checkout:

```bash
go run ./cmd/todoopen-server
```

Then open `http://127.0.0.1:8080/`.

---

## Contributing

Use pinned toolchain/tasks via `mise.toml`:

```bash
git clone https://github.com/justEstif/todo-open.git
cd todo-open
mise install
mise run build
mise run test
```

Common checks:

```bash
mise run fmt
mise run vet
mise run test
mise run build
```

Adapter runtime inspection:

```bash
todoopen adapters
todoopen adapters --json
```

---

## Docs

Deep-dive docs in [`docs/`](docs):

- [Architecture](docs/architecture.md)
- [API architecture](docs/api.md)
- [MVP scope](docs/mvp.md)
- [Canonical schema](docs/schema.md)
- [Sync decision](docs/sync.md)
- [Testing and release strategy](docs/testing.md)
- [Coding standards](docs/coding-standards.md)
