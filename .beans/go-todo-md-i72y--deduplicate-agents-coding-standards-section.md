---
# go-todo-md-i72y
title: Deduplicate AGENTS coding standards section
status: in-progress
type: task
priority: normal
created_at: 2026-03-05T20:19:59Z
updated_at: 2026-03-05T20:20:58Z
---

Reduce duplication in AGENTS.md by replacing the long coding standards section with a pointer to the canonical coding standards document and only repo-specific notes.\n\n## Tasks\n- [x] Edit AGENTS.md to replace duplicated coding standards details with concise references\n- [x] Ensure links use @docs path style\n- [x] Verify document still covers required commands and agent guidance\n\n## Summary of Changes\n- Replaced the long duplicated coding standards section in AGENTS.md with concise pointers to canonical docs.\n- Linked style and quality guidance to @docs/coding-standards.md, @docs/testing.md, @docs/schema.md, and @docs/api.md.\n- Kept only repo-specific implementation reminders in AGENTS.md (package boundaries, thin handlers, /v1 routes, ext + UTC behavior).
