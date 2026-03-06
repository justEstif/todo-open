---
status: draft
---

# Web UI Manual QA Checklist (MVP)

Run server via:

```bash
go run ./cmd/todoopen web --no-open
```

Open `http://127.0.0.1:8080/`.

## Smoke checks

- [ ] Page loads without console errors
- [ ] CSS and JS load from `/static/app.css` and `/static/app.js`
- [ ] Task list section renders with initial status text

## Task CRUD checks

- [ ] Create task with valid title; task appears in list
- [ ] Create task with blank title is rejected (status/error shown)
- [ ] Select **Edit** for a task; form is populated
- [ ] Update title and save; list reflects new title
- [ ] Delete selected task; task disappears from list

## Error handling checks

- [ ] Stop server and verify UI surfaces request failure message
- [ ] Start server again and verify refresh recovers

## Mobile/responsive checks

- [ ] In browser mobile viewport (e.g., 390x844), layout is usable
- [ ] Input fields/buttons are comfortably tappable
- [ ] No horizontal scrolling on core views

## Logging checks

- [ ] Startup log appears on launch
- [ ] Access logs include `method`, `path`, `status`, `bytes`, `duration`
- [ ] Shutdown log appears on Ctrl+C
