---
# go-todo-md-g8j3
title: Remove task ID from web UI
status: completed
type: task
priority: normal
created_at: 2026-03-06T01:50:18Z
updated_at: 2026-03-06T01:50:40Z
---

Hide internal task IDs in the HTML task list.\n\n## Todo\n- [x] Update task meta rendering to omit IDs\n- [x] Run quick API/web tests\n- [x] Add summary and complete bean

## Summary of Changes\n\n- Updated internal/api/web/static/app.js to remove task IDs from the task metadata line in the HTML list.\n- Kept status and updated-time metadata visible for context.\n- Verified with go test ./internal/api.
