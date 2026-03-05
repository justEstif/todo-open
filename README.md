# todo.open

todo.open is a **server-first, local-first task system** built in Go.

The project keeps task data portable (JSONL), enforces clear domain rules on the server, and supports multiple clients (CLI first, with web/mobile/TUI planned).

## What we are building

- A canonical **server API** for task lifecycle, validation, and sync workflows
- A local-first persistence model using JSONL with schema/version contracts
- A clean architecture that keeps domain logic independent from transport/client UI
- A path to add sync adapters and additional clients without changing core semantics

## Project documentation

All strategy and planning docs live in [`docs/`](docs):

- [Architecture](docs/architecture.md)
- [API architecture](docs/api.md)
- [MVP scope](docs/mvp.md)
- [Canonical schema](docs/schema.md)
- [Sync decision](docs/sync.md)
- [Testing and release strategy](docs/testing.md)
- [Coding standards](docs/coding-standards.md)

## Current status

Planning and architecture decisions are complete. The next phase is coding execution (project bootstrap, storage implementation, API endpoints, and CI baseline).
