---
# go-todo-md-bn4j
title: Create npm package for todoopen CLI
status: completed
type: task
priority: normal
created_at: 2026-03-08T21:27:36Z
updated_at: 2026-03-08T21:35:53Z
---

Package the todoopen CLI as an npm package (similar to qry's npm/ setup) so users can install via `npm install -g @justestif/todo-open`. Should wrap the Go binary download for the target platform, follow the same pattern as /home/estifanos/Documents/projects/qry/npm/.

## Summary of Changes

- `npm/todo-open/package.json` — `@justestif/todo-open` package exposing `todoopen` and `todoopen-server` bin entries
- `npm/todo-open/install.js` — postinstall script downloads the correct release archive from GitHub, handles redirects, extracts both binaries
- `npm/todo-open/bin.js` / `bin-server.js` — thin Node shims that exec the downloaded binaries with full stdio passthrough
- `.github/workflows/release.yml` — added `publish-npm` job: stamps version from tag, publishes with provenance; added `id-token: write` permission; checksums generated in publish job

Remaining: add `NPM_TOKEN` secret to the GitHub repo before the first tag push.
