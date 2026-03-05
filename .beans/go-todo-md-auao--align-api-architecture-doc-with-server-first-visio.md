---
# go-todo-md-auao
title: Align API architecture doc with server-first vision
status: completed
type: task
priority: normal
created_at: 2026-03-05T18:41:25Z
updated_at: 2026-03-05T18:42:12Z
parent: go-todo-md-yris
---

Update api.md to reflect Go + server-first architecture and ensure API boundaries align with architecture.md.\n\n## Todo\n- [x] Review api.md and architecture.md for mismatches\n- [x] Update api.md to server-first API model\n- [x] Validate terminology consistency across docs

## Summary of Changes\n\n- Rewrote api.md as a server-first API architecture document aligned with architecture.md.\n- Defined explicit API surfaces: Task API, Views API, Sync API, and MVP-light Admin/Operations API with conceptual endpoints.\n- Added contract rules, client model, deployment model, and Go internal package mapping consistent with server-first + local-first direction.\n- Clarified that CLI is an API client and direct JSONL mutation by clients is out-of-bound in server mode.
