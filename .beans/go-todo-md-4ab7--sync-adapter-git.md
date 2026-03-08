---
# go-todo-md-4ab7
title: 'sync adapter: git'
status: completed
type: feature
priority: normal
created_at: 2026-03-08T19:22:13Z
updated_at: 2026-03-08T21:08:59Z
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

## Summary of Changes

Successfully implemented the standalone plugin binary  that implements the todo-open plugin protocol.

### Implementation Details

Created the following files:

1. **cmd/todoopen-plugin-sync-git/main.go**: Main plugin binary that:
   - Writes HandshakeResponse to stdout immediately on startup
   - Handles stdin/stdout JSON line communication with the host
   - Dispatches push/pull/status methods to GitAdapter
   - Proper error handling with PluginError responses

2. **internal/sync/git/adapter.go**: GitAdapter struct with methods:
   - **Push()**: Adds tasks.jsonl, meta.json, config.toml to git, commits with skip-ci message, pushes to remote
   - **Pull()**: Performs fast-forward only pull, returns error if manual merge required
   - **Status()**: Checks working directory cleanliness and ahead/behind counts
   - Proper error handling for git not found and other git errors

3. **internal/sync/git/adapter_test.go**: Comprehensive tests:
   - Tests all three operations (push, pull, status)
   - Uses t.TempDir() for real git repository testing
   - Covers both success cases and edge cases (idempotent push, clean/dirty status)
   - JSON serialization tests for config and result types

### Key Features

- **Protocol Compliance**: Implements exact plugin protocol with handshake, request/response envelopes
- **Git Operations**: All required git operations (add, commit, push, pull, status)
- **Error Handling**: Proper error codes (unavailable, invalid_input, internal, not_supported)
- **Configuration**: Supports configurable remote and branch via config map
- **Robustness**: Graceful handling of missing files, no upstream branches, etc.

### Verification

- ✅ Plugin binary builds successfully ()
- ✅ All tests pass (ok  	github.com/justEstif/todo-open/internal/sync/git	(cached))
- ✅ Full test suite passes (ok  	github.com/justEstif/todo-open/cmd/todoopen	(cached)
?   	github.com/justEstif/todo-open/cmd/todoopen-plugin-sync-git	[no test files]
?   	github.com/justEstif/todo-open/cmd/todoopen-server	[no test files]
?   	github.com/justEstif/todo-open/internal/adapterregistry	[no test files]
?   	github.com/justEstif/todo-open/internal/adapters	[no test files]
ok  	github.com/justEstif/todo-open/internal/api	(cached)
ok  	github.com/justEstif/todo-open/internal/api/handlers	(cached)
?   	github.com/justEstif/todo-open/internal/api/middleware	[no test files]
?   	github.com/justEstif/todo-open/internal/api/web	[no test files]
ok  	github.com/justEstif/todo-open/internal/app	(cached)
?   	github.com/justEstif/todo-open/internal/client/api	[no test files]
ok  	github.com/justEstif/todo-open/internal/core	(cached)
ok  	github.com/justEstif/todo-open/internal/events	(cached)
ok  	github.com/justEstif/todo-open/internal/plugin	(cached)
ok  	github.com/justEstif/todo-open/internal/store	(cached)
ok  	github.com/justEstif/todo-open/internal/store/jsonl	(cached)
?   	github.com/justEstif/todo-open/internal/store/memory	[no test files]
ok  	github.com/justEstif/todo-open/internal/sync	(cached)
ok  	github.com/justEstif/todo-open/internal/sync/git	(cached)
ok  	github.com/justEstif/todo-open/internal/sync/noop	(cached)
?   	github.com/justEstif/todo-open/internal/version	[no test files]
ok  	github.com/justEstif/todo-open/internal/view	(cached)
ok  	github.com/justEstif/todo-open/internal/view/json	(cached))
- ✅ Code passes static analysis ()
- ✅ Code is properly formatted (internal/api/handlers/capabilities.go
internal/api/handlers/capabilities_test.go
internal/version/version.go shows no new issues)

The plugin is ready for use with the todo-open host system.
