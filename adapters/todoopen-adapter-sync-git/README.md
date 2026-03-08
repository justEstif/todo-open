# todoopen-adapter-sync-git

Sync adapter for [todo.open](https://justestif.github.io/todo-open) that pushes and pulls your `tasks.jsonl` workspace to/from a git repository.

Every sync commits the workspace files (`tasks.jsonl`, `meta.json`, `config.toml`) and pushes to your configured remote. Pull does a fast-forward-only fetch — it never silently merges.

---

## Install

```sh
# npm
npm install -g @justestif/todoopen-adapter-sync-git

# mise
mise use -g go:github.com/justEstif/todo-open/adapters/todoopen-adapter-sync-git@latest
mise reshim

# build from source
git clone https://github.com/justEstif/todo-open.git
cd todo-open
go build -o todoopen-adapter-sync-git ./adapters/todoopen-adapter-sync-git
```

**Requires:** `git` in your `PATH`.

---

## Setup

### 1. Initialize a git repo inside your workspace

Your todo.open workspace lives at `.todoopen/` inside your project (or a standalone directory). That folder needs to be a git repo:

```sh
cd .todoopen
git init
git remote add origin git@github.com:youruser/your-tasks-repo.git
```

Or use an existing repo — the adapter only touches `tasks.jsonl`, `meta.json`, and `config.toml`.

### 2. Make the first commit

```sh
cd .todoopen
git add tasks.jsonl meta.json config.toml
git commit -m "init: todo-open workspace"
git push -u origin main
```

### 3. Configure the adapter

Add to `.todoopen/config.toml`:

```toml
[adapters.git]
  bin = "todoopen-adapter-sync-git"

[adapters.git.config]
  remote = "origin"   # optional, default: origin
  branch = "main"     # optional, default: current HEAD branch
```

Use `${VAR}` to reference environment variables — they are expanded at runtime and never written to disk:

```toml
[adapters.git.config]
  remote = "${GIT_REMOTE}"
```

### 4. Verify

```sh
todoopen adapters
```

You should see the `git` adapter listed as `healthy`.

---

## GitHub Setup

If you want tasks synced to a private GitHub repo, the simplest setup is SSH:

```sh
# generate a key if you don't have one
ssh-keygen -t ed25519 -C "todo-open sync"

# add the public key to GitHub
# https://github.com/settings/ssh/new
cat ~/.ssh/id_ed25519.pub

# verify the connection
ssh -T git@github.com
```

Then set the remote to the SSH URL:

```sh
cd .todoopen
git remote add origin git@github.com:youruser/your-tasks-repo.git
```

For HTTPS with a token (e.g. in CI or on a machine without SSH):

```sh
git remote add origin https://github.com/youruser/your-tasks-repo.git
# git will prompt for credentials, or use a credential helper:
git config credential.helper store
```

---

## Usage

Sync is triggered through the todo.open CLI or server — you don't call the adapter binary directly.

```sh
# push local changes to remote
todoopen sync push

# pull remote changes
todoopen sync pull

# check sync status (ahead/behind/clean)
todoopen sync status
```

Or via the HTTP API if the server is running:

```sh
curl -s -X POST http://127.0.0.1:8080/v1/sync/push
curl -s -X POST http://127.0.0.1:8080/v1/sync/pull
curl -s        http://127.0.0.1:8080/v1/sync/status
```

---

## Behaviour

| Operation | What happens |
|---|---|
| `push` | Stages `tasks.jsonl`, `meta.json`, `config.toml` → commits with message `chore: sync todo-open workspace [skip ci]` → pushes to remote. If nothing changed, it exits cleanly with `"nothing to commit"`. |
| `pull` | Runs `git pull --ff-only <remote> <branch>`. Fails clearly if a merge would be required — no silent conflict creation. |
| `status` | Reports whether the workspace files are clean, and how many commits ahead/behind the remote tracking branch. |

The `[skip ci]` tag in the commit message prevents CI pipelines from triggering on sync commits — remove it from the source if you don't want that.

---

## Multi-machine workflow

The intended workflow for using tasks across machines:

```sh
# start of session — pull before making changes
todoopen sync pull

# do your work...
todoopen task create --title "..."
todoopen task done task_abc

# end of session — push
todoopen sync push
```

If you forget to pull first and get a diverged branch, resolve it manually:

```sh
cd .todoopen
git fetch origin
git rebase origin/main   # or merge — your choice
```

Then push again with `todoopen sync push`.

---

## Protocol

This adapter implements the `todoopen.plugin.v1` protocol over stdin/stdout JSON. Each line on stdin is a request envelope; each line on stdout is a response envelope. The adapter binary is started fresh per sync operation.

See [docs/adapters.md](../../docs/adapters.md) for the full protocol spec and how to build your own adapter.
