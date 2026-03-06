---
# go-todo-md-bce7
title: Build adapter set from user choices
status: completed
type: task
priority: normal
created_at: 2026-03-05T23:40:45Z
updated_at: 2026-03-05T23:59:45Z
---

Move adapter implementations from built-in runtime registration to examples and improve extension guidance docs.

## Todo
- [x] Remove markdown/table/git from built-in adapter registration
- [x] Move adapter implementations into examples/adapters (no adapter tests)
- [x] Update adapter docs to mark only json/noop as built-in
- [x] Add README extension wiring example linking to examples folder
- [x] Run formatting and verification tests
- [x] Update bean summary and mark completed

## Summary of Changes

- Reverted built-in runtime registration to json view + noop sync only in app composition and restored registry tests to those defaults.
- Moved the markdown/table view adapters and git sync adapter to compileable examples under examples/adapters, and removed their internal built-in packages and adapter-specific tests.
- Added extension wiring guidance and examples documentation in README and docs/adapters, including a concrete include-extension code snippet and direct links to examples/adapters/.
- Verified with gofmt, targeted registry tests, compile-check of example packages, and full go test ./... run.
