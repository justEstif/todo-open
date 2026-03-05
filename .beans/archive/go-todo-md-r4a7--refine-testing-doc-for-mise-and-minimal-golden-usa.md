---
# go-todo-md-r4a7
title: Refine testing doc for mise and minimal golden usage
status: completed
type: task
priority: normal
created_at: 2026-03-05T19:28:56Z
updated_at: 2026-03-05T19:30:01Z
parent: go-todo-md-yris
---

Update testing.md to use mise-based commands and simplify fixture/golden guidance to a minimal practical approach.\n\n## Todo\n- [x] Replace make-based command examples with mise run workflow\n- [x] Simplify golden/normalization guidance to minimal usage\n- [x] Share updated testing doc summary

## Summary of Changes\n\n- Updated testing.md command guidance from make/just to mise (mise run ...) for local and CI parity.\n- Simplified fixture/golden guidance to a minimal approach: fixtures broadly, goldens only for stable user-visible API/CLI contracts.\n- Reframed normalization as optional and only for actual flakiness, not default complexity.
