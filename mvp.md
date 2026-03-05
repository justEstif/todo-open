---
status: accepted
---

# todo.open MVP

## What

todo.open is a local-first task system built on JSONL that combines:

- A simple, portable task data format
- Schema validation and extension support
- Multiple clients (CLI, TUI, web/mobile)
- Composable querying with existing tools (e.g., Miller, VisiData)

The goal is to preserve the simplicity and durability of todo.txt while improving accessibility across devices and interfaces.

## Why

Many task systems force users into one UI, one sync provider, or one workflow. Plain-text systems are durable, but often difficult to use on mobile and hard to extend safely.

todo.open exists to provide:

- **Portable data first**: tasks remain user-owned and tool-readable
- **UI flexibility**: same data works in CLI, terminal UI, and mobile web
- **Structured interoperability**: JSONL + schema unlock robust tooling and automation
- **Extensible architecture**: views and sync are pluggable, not hardcoded

## Core Advantages

1. **Interoperable by default**
   - JSONL works with jq, Miller, VisiData, grep, and custom scripts.

2. **Accessible everywhere**
   - A web/PWA client can make the same task store usable on phones.

3. **Schema-backed confidence**
   - Validation reduces data drift while still allowing extension fields.

4. **Ecosystem hub model**
   - Third parties can build views and sync adapters without changing the core format.

5. **Local-first and future-proof**
   - Users can keep data in files, version control it, and migrate over time.

## MVP Scope

- Define a core task schema and metadata file
- Implement basic CRUD operations via CLI
- Add schema validation command
- Provide starter query recipes for Miller/VisiData
- Publish a first mobile-friendly web view (read/write)
- Define extension points for views and sync adapters

## Non-Goals (MVP)

- Complex collaboration workflows
- Enterprise permissions/roles
- Full event-sourcing from day one

## Success Criteria

- A user can create/manage tasks from CLI and mobile UI using the same JSONL files.
- Tasks validate against core schema and can include extension fields.
- Common reports can be generated with existing tools without vendor lock-in.
- A basic adapter contract exists for adding new sync backends.
