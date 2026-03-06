---
# go-todo-md-krxn
title: Fix web filter alignment
status: completed
type: bug
priority: normal
created_at: 2026-03-06T01:48:18Z
updated_at: 2026-03-06T01:48:58Z
---

Resolve alignment issue in Tasks filter controls on desktop layout.\n\n## Todo\n- [x] Group filter label+control pairs in HTML\n- [x] Adjust filter grid CSS for consistent alignment\n- [x] Run relevant web/API tests\n- [x] Add summary and complete bean

## Summary of Changes\n\n- Fixed filter alignment by wrapping each filter label/control pair in a dedicated field container in internal/api/web/static/index.html.\n- Updated filter layout styles in internal/api/web/static/app.css to align fields consistently across desktop and mobile breakpoints.\n- Verified behavior with go test ./internal/api.
