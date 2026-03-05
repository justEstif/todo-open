---
status: accepted
---

# todo.open Testing and Release Strategy (MVP)

This document defines the practical quality strategy for todo.open (Go, server-first, local-first JSONL).

## 1) Testing scope (MVP)

### Unit tests (priority P0)

- `internal/core`
  - schema validation and invariants
  - lifecycle transitions
  - merge policy helpers (LWW + tie-breakers)
- `internal/store`
  - JSONL read/write/update behavior
  - metadata persistence
  - atomic write/recovery behavior
- `internal/sync`
  - push/pull/status orchestration with fake adapters
  - token/checkpoint handling
  - conflict artifact generation flow
- `internal/api`
  - request validation and error mapping
  - endpoint-to-service wiring
  - response shape/status codes

### Integration tests (priority P0/P1)

P0:

1. API + Core + Store CRUD persistence correctness
2. Store durability under interrupted write simulation
3. Sync merge determinism and conflict artifact generation

P1:

1. API sync endpoints happy/failure path coverage
2. query/filter behavior defaults (e.g., soft-delete filtering)
3. validation endpoint strict/compat behavior

### Defer until post-MVP

- Real remote adapter end-to-end suites (WebDAV/S3/hosted)
- High-scale perf/load and compaction benchmarking
- Real-time multi-client behavior and subscription tests
- Full authz/security matrix

---

## 2) Fixtures and golden files (minimal approach)

Use canonical `testdata/` layouts per package for fixture inputs.

Example structure:

```text
internal/store/testdata/
  fixtures/
    tasks/minimal/{tasks.jsonl,meta.json}
    tasks/with-completed/{tasks.jsonl,meta.json}
    conflicts/concurrent-edit/{local.tasks.jsonl,remote.tasks.jsonl,expected.merged.jsonl,expected.conflicts.jsonl}
  goldens/
    api/get-tasks.minimal.golden.json
    cli/list.minimal.golden.txt
```

Guidance:

- Use fixtures broadly (`testdata`) for API/core/store/sync test inputs.
- Use golden files selectively for stable, user-visible contract outputs (API JSON or CLI text).
- Prefer explicit assertions over snapshots when behavior is small and easy to assert directly.

### Golden update policy

- Golden files are updated only with intentional behavior changes.
- Prefer separate commit for golden updates in PRs.
- CI must fail if tests mutate files unexpectedly.

### Optional normalization (only if needed)

Do not add heavy normalization by default.
Apply only when tests become flaky due to unstable values (timestamps, random IDs, temp paths, ordering).

---

## 3) CI quality gates (MVP)

Required checks on each PR (typically executed through `mise run ci` in CI):

1. `go mod tidy` cleanliness check
2. formatting check (`gofmt -l`)
3. `go vet ./...`
4. `golangci-lint run`
5. `go test ./... -race -coverprofile=coverage.out`
6. build binaries:
   - `go build ./cmd/todoopen-server`
   - `go build ./cmd/todoopen`
7. sync determinism smoke tests

Initial coverage target: ~60% overall, raised over time.

---

## 4) Release and versioning workflow

### Release pipeline (tag-driven)

On `vX.Y.Z` tag:

1. re-run full CI
2. build cross-platform binaries (server + CLI)
3. generate checksums (`SHA256SUMS`)
4. publish release artifacts and notes

### Versioning policy (pre-1.0)

- Start at `0.1.0`
- Use `0.MINOR.PATCH`
  - PATCH: fixes/internal improvements, no intended contract break
  - MINOR: features and possible contract changes (pre-1.0)
- Document all breaking changes explicitly in release notes.

---

## 5) Local commands mirroring CI

Recommended via `mise` tasks (defined in `mise.toml`):

- `mise run mod-tidy-check`
- `mise run fmt`
- `mise run vet`
- `mise run lint`
- `mise run test`
- `mise run test-race`
- `mise run build`
- `mise run modernize-check`
- `mise run ci` (full parity with required CI checks)

CI should call the same `mise run ...` tasks so local and PR workflows stay aligned.
