---
# go-todo-md-f11d
title: Plan migration and phased rollout for runtime adapter plugins
status: completed
type: task
priority: high
created_at: 2026-03-06T00:19:49Z
updated_at: 2026-03-06T00:26:56Z
parent: go-todo-md-dwib
---

Translate research into an executable plan with milestones, sequencing, compatibility strategy, and acceptance criteria.

## Deliverables
- [x] Define phase-by-phase implementation sequence
- [x] Specify compatibility/deprecation path for .todoopen/adapters.json
- [x] Define test strategy (unit/integration/e2e) and rollout gates
- [x] Produce implementation-ready task breakdown

## Summary of Changes

- Added planning doc at docs/adapter-plugin-plan.md.
- Defined 4 milestone sequence: protocol/runtime foundation, metadata source-of-truth, status+compatibility mode, and final removal/hardening.
- Updated plan to hard cutover: remove adapters.json runtime loading path with no compatibility/deprecation stages.
- Defined unit/integration/e2e test strategy plus rollout gates tied to CI and regression safety.
- Mapped implementation-ready execution to existing child beans and proposed execution order.
