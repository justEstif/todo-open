---
# go-todo-md-uc87
title: Fix mise version reporting and update README
status: completed
type: task
priority: normal
created_at: 2026-03-06T04:21:09Z
updated_at: 2026-03-06T04:22:15Z
---

Commit the CLI version fallback fix, update README install guidance for mise, push to main, and cut a new release tag.\n\n## Todo\n- [x] Review pending changes including README updates needed\n- [x] Update README install/version guidance for mise\n- [x] Commit code, docs, and bean files\n- [x] Push main to origin\n- [x] Create and push new release tag\n- [x] Summarize changes

## Summary of Changes\n\n- Updated CLI version reporting to prefer injected release version, then Go build info version, then fallback to dev.\n- Updated README mise install guidance to recommend pinned tags and verify with todoopen --version.\n- Committed and pushed changes to main.\n- Created and pushed release tag v0.1.2 to trigger a new release build.
