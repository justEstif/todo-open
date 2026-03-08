---
# go-todo-md-20x6
title: Add npm package for S3 adapter
status: completed
type: task
priority: normal
created_at: 2026-03-08T22:31:31Z
updated_at: 2026-03-08T22:32:02Z
---

Create npm/todoopen-adapter-sync-s3 package mirroring the git adapter npm package.

## Summary of Changes\n\nCreated npm/todoopen-adapter-sync-s3/ with three files mirroring the git adapter package:\n- package.json — package metadata, bin entry, postinstall hook\n- bin.js — thin JS wrapper that delegates to the native binary\n- install.js — downloads and extracts the correct release binary on postinstall
