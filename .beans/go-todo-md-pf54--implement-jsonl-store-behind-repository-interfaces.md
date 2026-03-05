---
# go-todo-md-pf54
title: Implement JSONL store behind repository interfaces
status: todo
type: task
created_at: 2026-03-05T19:37:01Z
updated_at: 2026-03-05T19:37:01Z
parent: go-todo-md-y0ga
---

Build the initial persistence layer behind clean interfaces in internal/store and internal/core boundaries.\n\n## Todo\n- [ ] Define repository interfaces for task CRUD and list operations\n- [ ] Implement JSONL read/write/update with atomic file writes\n- [ ] Add metadata file handling and schema version checks\n- [ ] Add unit tests for happy path and corruption/recovery scenarios
