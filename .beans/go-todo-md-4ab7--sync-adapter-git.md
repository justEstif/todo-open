---
# go-todo-md-4ab7
title: 'sync adapter: git'
status: todo
type: feature
created_at: 2026-03-08T19:22:13Z
updated_at: 2026-03-08T19:22:13Z
---

Implement a git sync adapter that pushes/pulls tasks.jsonl to a git repo branch.

This is the highest-value first sync adapter — developers get version history for free and need zero extra infrastructure.

## Scope
- Lives at `internal/sync/git` (built-in) or as a plugin binary `cmd/todoopen-plugin-sync-git`
- Config: remote, branch, commit message template
- Push: write tasks.jsonl to branch, commit, push
- Pull: fetch branch, read tasks.jsonl, return tasks
- Handle first-push (empty remote) gracefully
- Handle merge conflicts with a clear error (no silent data loss)

## Config example
```toml
[adapters.git]
  bin  = "todoopen-plugin-sync-git"
  kind = "sync"

[adapters.git.config]
  remote = "${GIT_REMOTE}"
  branch = "tasks"
```

## Decision needed
- Built-in (compiled into the binary) vs plugin binary (separate executable)
- Built-in is simpler for users; plugin keeps the core binary lean and lets users skip it
