---
# go-todo-md-twix
title: Research TypeScript/Node implementation choice for todo.open
status: completed
type: task
priority: high
created_at: 2026-03-05T18:28:22Z
updated_at: 2026-03-05T18:29:38Z
parent: go-todo-md-yris
---

Evaluate TypeScript/Node for local-first JSONL CLI + future sync.\n\n- [x] Assess portability and distribution\n- [x] Assess ecosystem and developer velocity\n- [x] Compare commander vs oclif for MVP\n- [x] Assess JSONL/file I/O ergonomics\n- [x] Assess testing/tooling/release ergonomics\n- [x] Compare tradeoffs vs Go\n- [x] Provide recommendation and when TS is wrong choice

## Summary of Changes\n\nResearched TypeScript/Node as an implementation choice for todo.open (local-first JSONL CLI with future sync), with focused comparison against Go. Evaluated portability/distribution options (npm/pnpm/npx/global install vs single-executable SEA), ecosystem and development velocity, JSONL/file I/O ergonomics in Node fs/streams, and testing/release tooling. Compared Commander and oclif for MVP and recommended Commander for lower complexity and faster iteration, while noting oclif strengths for plugin-driven enterprise-style CLIs.
