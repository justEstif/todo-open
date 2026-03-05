---
# go-todo-md-48q3
title: Define canonical JSONL task/data schema
status: completed
type: task
priority: normal
created_at: 2026-03-05T18:14:49Z
updated_at: 2026-03-05T18:16:24Z
parent: go-todo-md-yris
---

Specify the canonical JSONL record schema for todo.open tasks, including required fields, lifecycle/state transitions, timestamps, and validation rules.\n\n## Todo\n- [x] Review existing MVP/API docs for schema constraints\n- [x] Draft schema document with fields and examples\n- [x] Define status transition and timestamp rules\n- [x] Share output and next implementation hooks

## Summary of Changes\n\n- Added schema.md defining the canonical todo.open JSONL task schema (required/optional fields, ext namespace, and sample records).\n- Defined lifecycle transition rules (open, in_progress, done, archived) and timestamp/version invariants.\n- Added workspace metadata contract (.todoopen/meta.json) and clear implementation hooks for Core API methods.
