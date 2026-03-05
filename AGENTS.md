# AGENTS.md

Guidance for coding agents working in `todo-open`.

## Project Snapshot

- Language: Go (`go 1.26` in `go.mod`)
- Module: `github.com/ebeyene/todo-open`
- Architecture: server-first, local-first task system
- Entrypoints:
  - `cmd/todoopen-server` (HTTP server)
  - `cmd/todoopen` (CLI client)
- Core package layout:
  - `internal/api` transport handlers/router
  - `internal/app` composition root
  - `internal/core` domain models/contracts
  - `internal/store` storage contracts + impls
  - `internal/client/api` server client used by CLI

## Source of Truth Docs

Read these before significant design or behavioral changes:

- `@docs/architecture.md`
- `@docs/api.md`
- `@docs/schema.md`
- `@docs/testing.md`
- `@docs/coding-standards.md`

## Environment and Tooling

- Required Go toolchain: 1.26
- Task runner: `mise` (partial tasks currently defined)
- Modernization checks in `mise.toml`:
  - `mise run modernize-check` -> `go fix -diff ./...`
  - `mise run modernize` -> `go fix ./...`
- Lint tool expected by standards: `golangci-lint`
  - It may not be installed in all local environments; install before CI-quality runs.

## Build, Lint, Test Commands

Run commands from repository root.

### Build

- Build all packages: `go build ./...`
- Build server binary package: `go build ./cmd/todoopen-server`
- Build CLI binary package: `go build ./cmd/todoopen`

### Format

- Apply formatting: `gofmt -w .`
- Check formatting only: `gofmt -l .`

### Vet and Lint

- Static checks: `go vet ./...`
- Lint (when installed): `golangci-lint run`
- Modernization report: `go fix -diff ./...`
- Apply modernization: `go fix ./...`

### Tests

- Run full test suite: `go test ./...`
- Run with race + coverage (CI-style): `go test ./... -race -coverprofile=coverage.out`
- Verbose tests: `go test -v ./...`

### Run a Single Test (important)

- Single test by name in one package:
  - `go test ./internal/core -run '^TestTaskCreate$'`
- Single subtest:
  - `go test ./internal/core -run '^TestTaskCreate$/invalid_title$'`
- Single package only:
  - `go test ./internal/store/jsonl`
- Re-run without test cache:
  - `go test ./internal/core -run '^TestTaskCreate$' -count=1`

Note: the current codebase may have few/no `_test.go` files yet; use these commands as the standard pattern as tests are added.

## Local Run Commands

- Start server (default `:8080`): `go run ./cmd/todoopen-server`
- Override server bind address:
  - `TODOOPEN_SERVER_ADDR=127.0.0.1:8081 go run ./cmd/todoopen-server`
- Run CLI against default server:
  - `go run ./cmd/todoopen`
- Run CLI against custom server:
  - `go run ./cmd/todoopen --server http://127.0.0.1:8081`

## Coding Style and Conventions

Canonical style and quality guidance lives in:

- `@docs/coding-standards.md`
- `@docs/testing.md`
- `@docs/schema.md`
- `@docs/api.md`

When those docs conflict with examples in code, follow the docs.

Repository-specific reminders for agents:

- Keep package boundaries aligned with `@docs/architecture.md`.
- Keep handlers thin in `internal/api`; put domain invariants in `internal/core`.
- Use versioned domain routes under `/v1/...`; keep `/healthz` and `/readyz` lightweight.
- Keep task extensions under `ext` and preserve RFC3339 UTC timestamp behavior.

## Agent Workflow Expectations

- Keep changes small, scoped, and behavior-preserving unless feature work requires otherwise.
- Do not suppress lints without a documented reason.
- Prefer separate commits for large modernization (`go fix`) changes.
- Update docs when commands, contracts, or package responsibilities change.

## Cursor and Copilot Rules

- `.cursorrules`: not found
- `.cursor/rules/`: not found
- `.github/copilot-instructions.md`: not found

If any of these files are added later, treat them as higher-priority agent instructions and update this document accordingly.
