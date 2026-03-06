---
# go-todo-md-vkke
title: Add GitHub Actions release workflow on tag push
status: completed
type: task
priority: normal
created_at: 2026-03-06T03:07:11Z
updated_at: 2026-03-06T03:08:23Z
---

Set up .github/workflows/release.yml to run on push tags (v*), build Go binaries, and create GitHub release artifacts similar to hivemind style releases.\n\n## Todo\n- [x] Inspect existing workflows and project build commands\n- [x] Add release workflow triggered by tag push\n- [x] Validate workflow YAML syntax and references\n- [x] Summarize usage for tagging releases

## Summary of Changes\n\n- Added  triggered on tag pushes matching .\n- Workflow builds todoopen and todoopen-server for linux/darwin (amd64, arm64) and windows (amd64).\n- Each matrix build packages binaries into versioned archives and uploads them as artifacts.\n- A publish job gathers all archives and creates/updates the GitHub Release with generated release notes and attached assets.

- Explicitly configured workflow file path: .github/workflows/release.yml, triggered by push tags matching v*.
