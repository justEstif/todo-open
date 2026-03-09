---
# go-todo-md-kyek
title: Suppress server logging when running embedded server in TUI mode
status: completed
type: bug
priority: normal
created_at: 2026-03-09T00:28:15Z
updated_at: 2026-03-09T00:28:37Z
---

The server's HTTP request logging middleware writes to log.Default() (stderr), which corrupts the bubbletea TUI display. Need to discard or redirect log output when starting the embedded server from the tui subcommand.

## Summary of Changes\n\nAdded `log.SetOutput(io.Discard)` in `runTui` before starting the embedded server, so the HTTP request logging middleware no longer writes to stderr and corrupts the bubbletea TUI display.
