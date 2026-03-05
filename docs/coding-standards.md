---
status: accepted
---

# todo.open Go Coding Standards (2026)

This document defines the coding standards for Go code in todo.open.

## Official baseline

All code should follow these official Go sources:

1. Effective Go
   - https://go.dev/doc/effective_go
2. Go Code Review Comments
   - https://github.com/golang/go/wiki/CodeReviewComments

These are the primary style and idiom references for this repository.

## Why modernization is explicit in this repo

The Go team has documented that coding assistants/LLMs often generate older Go idioms from historical training data. In Go 1.26, `go fix` was rewritten with modernizer analyzers to help update code to newer, clearer language/library patterns.

References:

- Using go fix to modernize Go code
  - https://go.dev/blog/gofix
- Go 1.26 release notes/blog
  - https://go.dev/blog/go1.26
- modernize analyzers package
  - https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/modernize

## Required tooling and checks

For all contributions:

1. `gofmt` formatting is required
2. `go vet ./...` must pass
3. `golangci-lint run` must pass
4. `go test ./...` must pass
5. Modernization check should be run: `go fix -diff ./...`

Notes:
- `go fix -diff` is a check/report mode for review.
- `go fix ./...` can be used in dedicated modernization commits.

## Practical guidance for humans and coding agents

- Prefer current idioms over legacy patterns when equivalent.
- Do not suppress lints without clear justification.
- Keep edits small and behavior-preserving.
- When upgrading Go versions, run `go fix` and review changes in a separate commit where practical.

## Local workflow

Use `mise` tasks where available:

- `mise run modernize-check` (runs `go fix -diff ./...`)
- `mise run modernize` (applies `go fix ./...`)

CI should include modernization checks as part of repo quality gates.
