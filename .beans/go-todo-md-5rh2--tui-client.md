---
# go-todo-md-5rh2
title: TUI client
status: todo
type: feature
created_at: 2026-03-08T19:22:36Z
updated_at: 2026-03-08T19:22:36Z
---

Build a terminal UI client for todo.open using Bubble Tea.

The server-first architecture means the TUI is just another HTTP client — it talks to the local server over loopback, the same as the CLI and web UI.

## Scope
- Lives at `cmd/todoopen-tui` or as `todoopen tui` subcommand
- Uses [Bubble Tea](https://github.com/charmbracelet/bubbletea) for the TUI framework
- Consumes `internal/client/api` HTTP client (same as CLI)
- Core views: task list, task detail, create/edit task, filter/sort

## MVP interactions
- Navigate task list (arrow keys / j/k)
- Open task detail pane
- Create task inline
- Toggle status (open → done)
- Filter by status, priority

## Nice to have
- Real-time refresh (poll or server-sent events)
- Fuzzy search
- Inline edit all fields

## Notes
- Keep TUI stateless — all mutations go through the server API, never direct file access
- Styling via Lip Gloss (standard Bubble Tea companion library)
