---
status: accepted
---

# todo.open Sync Decision (MVP)

## Decision

For MVP, todo.open will use an **adapter-first sync architecture** with **one concrete file-exchange adapter**.

- Architecture: plugin adapter contract (`init`, `pull`, `push`, `status`)
- Initial adapter: local/file-exchange (Git/shared-folder/export-import friendly)
- Conflict policy: **field-level LWW** with deterministic tie-breakers

This balances speed, local-first durability, and future extensibility.

---

## Why this approach

1. Keeps MVP implementation small and understandable.
2. Preserves offline-first behavior (no hosted backend required).
3. Matches server-first architecture by keeping domain conflict rules canonical.
4. Leaves a clean path for future WebDAV/S3/HTTP sync adapters.

---

## Server-first alignment

- In server mode, merge/conflict decisions are server-authoritative.
- In file-exchange workflows, the local server applies the same canonical merge policy.
- Clients (CLI/web/TUI) do not define merge semantics.

---

## MVP sync contract

## Adapter interface (conceptual)

- `init(config) -> AdapterHandle`
- `pull(sinceToken?) -> ChangeSet`
- `push(changeSet, baseToken?) -> PushResult`
- `status() -> SyncStatus`

## API surface (conceptual)

Sync HTTP endpoints are defined in `api.md` (Sync API section) as the canonical transport contract.
This document focuses on sync behavior/policy decisions rather than endpoint duplication.

---

## Conflict policy (MVP)

Primary merge strategy:

1. Field-level last-write-wins
2. Primary comparator: `updated_at`
3. Tie-breaker: lexical compare on `device_id`/`source_id`

Conflict handling:

- Persist conflict records under `.todoopen/meta/conflicts/`
- Preserve losing values in conflict artifacts (not task payload)
- Expose conflict visibility via CLI/API

Special rules:

- `status=done` invariants with `completed_at` must be revalidated post-merge
- `parent_id` integrity checks run after merge (no self-parent/cycle/orphan)
- `deleted_at` wins only when newer than competing field updates

---

## MVP scope

In scope:

- Adapter contract + one file-exchange adapter
- Pull/push/checkpoint flow
- Deterministic field-level LWW merge
- Conflict artifact persistence + inspection commands

Out of scope (defer):

- Hosted multi-tenant sync service
- Real-time subscriptions/presence
- Real-time agent coordination/orchestration sync
- CRDT/OT merge engines
- Team authz/collaboration semantics in core model

---

## Future opportunity: AI agent orchestration

Real-time sync is intentionally out of MVP scope, but it is a meaningful future unlock for agent-heavy workflows.

Potential benefits once introduced:

- Faster shared state propagation across parallel agents
- Reduced coordination lag for planning/replanning loops
- Better support for near-real-time task claiming, handoff, and conflict visibility

This should be reconsidered after MVP server semantics and deterministic sync correctness are stable.

---

## Quality gates (MVP)

1. Deterministic merge outcomes for same inputs
2. No silent data loss for independent field edits
3. Conflict artifact generated for every true collision
4. Full offline operation without remote dependency
5. Post-merge schema and lifecycle invariants always pass

---

## Follow-up implementation tasks

1. Define exact `ChangeSet`, `PushResult`, and token schemas.
2. Implement sync conflict artifact schema and retention policy.
3. Add `todo sync status` and `todo sync conflicts` CLI commands.
4. Add sync test suite for determinism, replay safety, and conflict invariants.
