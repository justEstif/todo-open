---
# go-todo-md-9ly5
title: Redesign web UI without Simple.css
status: completed
type: feature
priority: normal
created_at: 2026-03-06T01:42:43Z
updated_at: 2026-03-06T01:47:45Z
---

Replace Simple.css with a custom minimal-productivity design system while keeping Go-served HTML/CSS/JS and existing APIs only.\n\n## Todo\n- [x] Redesign internal/api/web/static/index.html structure and semantics\n- [x] Replace stylesheet with custom internal/api/web/static/app.css and remove Simple.css usage\n- [x] Refactor internal/api/web/static/app.js for improved UX using existing APIs\n- [x] Update web integration tests for new asset references and markers\n- [x] Run focused and full Go tests\n- [x] Add summary and complete bean

## Summary of Changes\n\n- Replaced the old Simple.css-based page shell with a custom minimal-productivity layout, including skip link, runtime health messaging, filters, and a clearer selected-task panel.\n- Added a new first-party stylesheet at internal/api/web/static/app.css and removed internal/api/web/static/simple.css from the static surface.\n- Refactored internal/api/web/static/app.js for local filtering/sorting, optimistic task updates, inline field errors, runtime health checks via existing endpoints, and delete confirmation safeguards.\n- Updated web integration coverage to validate app.css and new DOM markers, and refreshed manual QA docs for the new asset path.
