---
# go-todo-md-3wrl
title: 'Improve clarity: README marketing + adapter config restructuring'
status: scrapped
type: task
priority: normal
created_at: 2026-03-08T19:13:31Z
updated_at: 2026-03-08T19:35:45Z
---

Two tracks:
Track A: README rewrite — sharper tagline, concrete opening example, adapter table, cleaner structure inspired by qry.
Track B: Config restructuring — separate adapter config into .todoopen/config.toml (TOML, named sections, env var expansion, settings inlined); strip adapter fields from meta.json.

- [ ] Add BurntSushi/toml dependency
- [ ] Define AdapterFileConfig TOML types in internal/app/adapterfile.go
- [ ] Strip adapter fields (enabled_views, enabled_sync_adapters, adapter_plugins, ext.adapter_settings) from WorkspaceMeta / workspace_meta.go
- [ ] Update BuildAdapterRuntimeFromMeta -> BuildAdapterRuntimeFromConfig in adapterconfig.go
- [ ] Update server.go to load config.toml alongside meta.json
- [ ] Update all tests (workspace_meta_test.go, adapterconfig_test.go)
- [ ] Rewrite README.md
- [ ] Update docs/adapters.md

## Reasons for Scrapping
Duplicate of go-todo-md-7h4f which was completed and committed.
