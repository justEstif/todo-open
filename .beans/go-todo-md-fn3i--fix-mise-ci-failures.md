---
# go-todo-md-fn3i
title: Fix mise ci failures
status: completed
type: task
priority: normal
created_at: 2026-03-08T22:30:30Z
updated_at: 2026-03-08T22:30:53Z
---

Fix: mod-tidy-check, fmt (tui/model.go, tui/render.go), modernize-check (client/api/client.go, tui/model.go)

## Summary of Changes\n\n- : promoted  and  to direct deps, removed unused  indirect dep\n- : formatted  and \n- : modernized + →  in ; removed hand-rolled  in ; replaced manual min/max patterns with builtin  in 
