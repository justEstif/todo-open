---
# go-todo-md-jr4l
title: Refresh AGENTS.md guidance
status: completed
type: task
priority: normal
created_at: 2026-03-05T23:27:52Z
updated_at: 2026-03-05T23:30:45Z
---

Update /Users/ebeyene/Documents/personal/todo-open/AGENTS.md after analyzing repository commands and coding standards. Include build/lint/test commands (with single-test guidance), style rules, and references to Cursor/Copilot rules if present.\n\n## Todo\n- [x] Inspect current AGENTS.md and project docs for source-of-truth instructions\n- [x] Verify whether Cursor/Copilot instruction files exist\n- [x] Draft improved AGENTS.md (~150 lines) with commands and coding standards\n- [x] Validate accuracy against repo configuration files\n- [x] Add summary of changes and complete bean

## Summary of Changes\n\n- Reworked  into a comprehensive agent guide focused on this repository's Go 1.26 server-first architecture.\n- Updated build/lint/test command coverage using  and  tasks from , with a dedicated single-test section and concrete regex examples.\n- Expanded coding standards guidance for package boundaries, formatting/imports, types, naming, error handling, HTTP conventions, and test practices based on docs and existing code.\n- Verified Cursor/Copilot instruction file locations and documented current status (, ,  all not present).

## Summary of Changes (Corrected)

- Reworked AGENTS.md into a comprehensive agent guide focused on this repository and its Go 1.26 server-first architecture.
- Updated build/lint/test command coverage using go commands and mise tasks from mise.toml, with a dedicated single-test section and regex examples.
- Expanded coding standards guidance for package boundaries, formatting and imports, types, naming, error handling, HTTP conventions, and testing practices.
- Verified Cursor/Copilot instruction file locations and documented that .cursorrules, .cursor/rules/, and .github/copilot-instructions.md are not present.
