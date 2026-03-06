---
# go-todo-md-soxg
title: Evaluate configurable parser/data type beyond JSONL
status: completed
type: task
priority: normal
created_at: 2026-03-05T23:06:36Z
updated_at: 2026-03-05T23:06:48Z
---

Review existing adapter/data model beans and assess implications of allowing user-provided parser with configurable data type (not JSONL-only). Propose architecture impact, API/config changes, validation, and migration strategy.

## Summary of Changes

Reviewed current beans/docs and assessed implications of moving from JSONL-only assumptions to a configurable parser/data-type model. Key recommendation: keep canonical domain schema stable while introducing pluggable codec/store adapters behind interfaces, with strict validation and migration/import boundaries. Provided concrete impact areas (architecture, API/config, validation, testing, migration, safety constraints).
