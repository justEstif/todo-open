# AGENTS.md

Guidance for coding agents working in `todo-open`.

## Project Snapshot

- Language: Go (`go 1.26` in `go.mod`)
- Module: `github.com/justEstif/todo-open`
- Architecture: server-first, local-first task system
- Entrypoints:
  - `cmd/todoopen-server` (HTTP server)
  - `cmd/todoopen` (CLI client)
- Core package layout:
  - `internal/api` transport handlers/router
  - `internal/app` composition root
  - `internal/core` domain models/contracts
  - `internal/store` persistence contracts + implementations
  - `internal/client/api` HTTP client used by CLI

## Source-of-Truth Docs

Read these before significant design or behavior changes:

- `@docs/architecture.md`
- `@docs/api.md`
- `@docs/schema.md`
- `@docs/testing.md`
- `@docs/coding-standards.md`

When examples in code conflict with docs, follow docs.

## Environment and Tooling

- Required Go toolchain: `1.26`
- Task runner: `mise` (`mise run ci` for CI parity)
- Lint task: `mise run lint`
- Modernization tasks:
  - `mise run modernize-check` (runs `go fix -diff ./...`)
  - `mise run modernize` (runs `go fix ./...`)

## Build, Lint, and Test Commands

Run all commands from repository root.

### Build

- Build all packages: `go build ./...`
- Build server binary package: `go build ./cmd/todoopen-server`
- Build CLI binary package: `go build ./cmd/todoopen`
- Build both binaries through task runner: `mise run build`

### Format and Module Hygiene

- Check formatting only: `gofmt -l .`
- Apply formatting: `gofmt -w .`
- Check module tidiness: `go mod tidy -diff`
- Module tidy check task: `mise run mod-tidy-check`

### Vet and Lint

- Static checks: `go vet ./...`
- Lint: `mise run lint`
- Direct lint (if installed): `golangci-lint run`
- Modernization report: `go fix -diff ./...`

### Tests

- Full suite: `go test ./...`
- Verbose: `go test -v ./...`
- Race + coverage (CI style): `go test ./... -race -coverprofile=coverage.out`
- Task runner equivalents: `mise run test`, `mise run test-race`

### Run a Single Test (Important)

- One test in one package:
  - `go test ./internal/core -run '^TestValidateTaskJSONLValidRecord$'`
- One subtest:
  - `go test ./internal/core -run '^TestValidateTaskJSONL$/strict_mode$'`
- Multiple tests by regex:
  - `go test ./internal/store/jsonl -run 'TestTaskRepo(CRUDAndMetaBootstrap|RejectsCorruptJSONL)$'`
- Entire single package:
  - `go test ./internal/store/jsonl`
- Disable cache when iterating:
  - `go test ./internal/core -run '^TestValidateTaskJSONLValidRecord$' -count=1`

## Local Run Commands

- Start server: `go run ./cmd/todoopen-server`
- Common server overrides:
  - `TODOOPEN_SERVER_ADDR=127.0.0.1:8081 go run ./cmd/todoopen-server`
  - `TODOOPEN_STORE=memory go run ./cmd/todoopen-server`
  - `TODOOPEN_WORKSPACE_ROOT=/path/to/workspace go run ./cmd/todoopen-server`
- Run CLI: `go run ./cmd/todoopen`

## Coding Style and Conventions

Baseline references:

- Effective Go
- Go Code Review Comments
- `@docs/coding-standards.md`

### Package and Boundary Rules

- Keep transport concerns in `internal/api`.
- Keep domain invariants and validation in `internal/core`.
- Keep persistence details in `internal/store`.
- Keep wiring in `internal/app`; avoid domain logic there.
- CLI should consume API/client contracts; do not bypass server semantics in CLI features.

### Formatting and Imports

- Always run `gofmt`; do not hand-format alignment.
- Prefer standard import grouping: stdlib first, then blank line, then internal module imports.
- Avoid unused imports and aliases unless aliasing resolves a real name conflict.
- Keep files ASCII unless existing file or feature requires Unicode.

### Types and Data Modeling

- Use typed string enums for constrained domains (e.g., `TaskStatus`, `TaskPriority`).
- Keep JSON field names snake_case in tags to match schema/API.
- Use pointers for optional timestamps/fields where absence is meaningful.
- Keep extension data under `ext` only; do not add custom top-level task fields.
- Preserve RFC3339 UTC behavior for all persisted/API timestamps.

### Naming

- Exported identifiers: `PascalCase`; unexported: `camelCase`.
- Prefer concise, domain-specific names (`TaskRepo`, `ValidationIssue`, `NewServer`).
- Boolean names should read as predicates (`Ready`, `Healthy`, `enabled`).
- Keep receiver names short and consistent (`s`, `r`, `c`, `h`).

### Error Handling

- Return errors instead of panicking (except fatal process startup in `main`).
- Define sentinel errors for domain categories (e.g., `ErrInvalidInput`, `ErrNotFound`).
- Wrap errors with `%w` when adding context.
- Map domain errors to HTTP status codes in handlers; keep raw internal errors out of responses.
- Ignore errors only when explicitly safe and intentional (for example, best-effort response encoding).

### HTTP and API Conventions

- Keep `/healthz` lightweight.
- Keep domain routes versioned under `/v1/...`.
- Use strict JSON decoding for request bodies (`DisallowUnknownFields`).
- Return consistent JSON error envelopes from handlers.

### Testing Expectations

- Prefer table-driven tests when it improves clarity; keep simple tests direct.
- Use `t.Parallel()` for independent tests.
- Use `t.TempDir()` for filesystem-backed repository tests.
- Validate user-visible behavior with explicit assertions.
- Add/adjust tests in the same change when behavior changes.

## Agent Workflow Expectations

- Keep changes small, scoped, and behavior-preserving unless feature work requires broader changes.
- Do not suppress lint warnings without a documented reason.
- Prefer separate commits for broad modernization (`go fix`) changes.
- Update docs when contracts, commands, or package responsibilities change.

## Cursor and Copilot Rules

- `.cursorrules`: not found
- `.cursor/rules/`: not found
- `.github/copilot-instructions.md`: not found

If these files are added later, treat them as higher-priority repository instructions and update this file.
