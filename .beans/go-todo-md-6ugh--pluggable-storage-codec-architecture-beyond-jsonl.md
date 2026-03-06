---
# go-todo-md-6ugh
title: Pluggable storage codec architecture (beyond JSONL)
status: todo
type: feature
priority: normal
created_at: 2026-03-05T23:07:30Z
updated_at: 2026-03-05T23:25:00Z
---

Introduce a format-agnostic storage codec architecture so todo.open is not limited to JSONL. Preserve canonical core task schema/invariants while allowing configurable codecs/parsers and future custom adapters.

\n## Scope additions (privacy + external viewers)\n\nEnsure this feature explicitly supports:\n\n- User-controllable storage codec selection (not hardcoded JSONL only).\n- A stable codec contract that enables third-party read/write adapters.\n- Documented format metadata/versioning so external tools can reliably parse files.\n- Privacy-oriented codec composition (e.g., serialize -> compress -> encrypt) where encryption/key handling is explicit and not implied by file extension alone.\n- Support for external read-only viewer tooling for non-encrypted formats, and key-authorized viewing flows for encrypted formats.\n\n## Acceptance criteria additions\n\n- CLI/config supports selecting codec pipeline components exposed to users.\n- At least one documented example of building an external viewer against the codec/format contract.\n- Security docs clearly distinguish compression/obfuscation from actual encryption/privacy guarantees.\n
