---
# go-todo-md-3x7p
title: Add CLI --version flag
status: completed
type: task
priority: normal
created_at: 2026-03-06T03:10:29Z
updated_at: 2026-03-06T03:12:08Z
---

Add a global CLI version flag to print the application version and exit.\n\n## Todo\n- [x] Inspect current CLI command setup\n- [x] Implement --version support with build-time version var fallback\n- [x] Add/adjust tests for version flag behavior\n- [x] Run focused tests\n- [x] Summarize changes

## Summary of Changes\n\n- Added global CLI version support in cmd/todoopen: --version, -v, and version now print todoopen <version> and exit 0.\n- Introduced build-time version variable "version" with default fallback "dev" for non-release builds.\n- Updated help text to document the new version flag.\n- Added tests covering help output and version output behavior (TestVersionCommand).\n- Updated release workflow to embed the tag into CLI builds via -ldflags -X main.version=${GITHUB_REF_NAME}.\n- Verified with go test ./cmd/todoopen.
