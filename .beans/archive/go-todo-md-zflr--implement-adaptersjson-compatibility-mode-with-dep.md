---
# go-todo-md-zflr
title: Implement adapters.json compatibility mode with deprecation flow
status: scrapped
type: task
priority: normal
created_at: 2026-03-06T00:20:17Z
updated_at: 2026-03-06T00:26:11Z
parent: go-todo-md-dwib
---

Introduce temporary compatibility with .todoopen/adapters.json while guiding users toward metadata-based registration.

## Deliverables
- [ ] Compatibility read path implemented behind clear precedence rules
- [ ] Deprecation warnings emitted with migration guidance
- [ ] Removal criteria/version target defined
- [ ] Tests covering compatibility and precedence behavior

## Reasons for Scrapping

Direction changed: remove .todoopen/adapters.json runtime support directly with no deprecation/compatibility window. Compatibility-mode work is intentionally out of scope.
