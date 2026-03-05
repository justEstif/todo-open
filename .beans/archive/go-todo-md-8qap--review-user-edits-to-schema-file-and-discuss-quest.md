---
# go-todo-md-8qap
title: Review user edits to schema file and discuss questions
status: completed
type: task
priority: normal
created_at: 2026-03-05T18:23:39Z
updated_at: 2026-03-05T18:24:36Z
---

Review recent user changes to schema-related file(s), summarize implications, and discuss open questions.\n\n## Todo\n- [x] Inspect changed file content and diff\n- [x] Identify notable changes and tradeoffs\n- [x] Discuss user questions and recommendations

## Summary of Changes\n\n- Reviewed  and identified user-added questions around  purpose and authentication placement.\n- Provided recommendations: keep  for local task hierarchy (not team coupling) and keep auth out of core schema, scoped to sync adapters.\n- Noted formatting issue: use Markdown comments instead of double-slash lines in doc content.

## Corrected Summary of Changes\n\n- Reviewed schema.md and identified user questions about parent_id and authentication scope.\n- Recommended keeping parent_id for local hierarchical subtasks, independent of team features.\n- Recommended keeping authentication outside core schema and implementing auth in sync adapters.\n- Recommended replacing double-slash inline notes with Markdown-style comments or a Questions section.
