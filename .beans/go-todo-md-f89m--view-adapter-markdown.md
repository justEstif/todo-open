---
# go-todo-md-f89m
title: 'view adapter: markdown'
status: todo
type: feature
created_at: 2026-03-08T19:22:26Z
updated_at: 2026-03-08T19:22:26Z
---

Implement a markdown view adapter that renders tasks as a markdown checklist.

## Scope
- Lives at `internal/view/markdown`
- Renders tasks as a markdown checklist grouped by status or priority
- Useful for exporting to notes, GitHub issues, PRs, or documentation

## Output example
```markdown
## Open
- [ ] Ship the release (high)
- [ ] Write release notes (normal)

## Done
- [x] Refactor adapter config
```

## Config example
```toml
[adapters.markdown]
  bin  = "todoopen-plugin-view-markdown"
  kind = "view"

[adapters.markdown.config]
  group_by = "status"  # or "priority"
```
