---
# go-todo-md-h6zc
title: External viewer contract + sample viewer for storage codecs
status: todo
type: task
created_at: 2026-03-05T23:25:23Z
updated_at: 2026-03-05T23:25:23Z
parent: go-todo-md-6ugh
---

Define and document a stable read-focused codec/viewer contract and provide at least one example external viewer implementation.\n\n## Goals\n\n- Make it straightforward for third parties to build read-only viewers for supported non-encrypted formats.\n- Define key-authorized viewing expectations for encrypted formats.\n- Validate format metadata/versioning is sufficient for external parsing.\n\n## Todo\n\n- [ ] Specify viewer-facing codec contract (interfaces + guarantees).\n- [ ] Document format metadata/versioning requirements.\n- [ ] Add one sample external viewer (minimal reference implementation).\n- [ ] Document encrypted-viewer flow and key requirements.\n- [ ] Add docs section describing compatibility and limitations.
