const runtimeEl = document.getElementById('runtime-status');
const statusEl = document.getElementById('status');
const selectedMetaEl = document.getElementById('selected-meta');
const taskListEl = document.getElementById('task-list');
const createForm = document.getElementById('create-form');
const createTitle = document.getElementById('create-title');
const createSubmit = document.getElementById('create-submit');
const createErrorEl = document.getElementById('create-error');
const filterForm = document.getElementById('filter-form');
const queryInput = document.getElementById('query-input');
const statusFilter = document.getElementById('status-filter');
const sortSelect = document.getElementById('sort-select');
const editForm = document.getElementById('edit-form');
const editSubmit = document.getElementById('edit-submit');
const editID = document.getElementById('edit-id');
const editTitle = document.getElementById('edit-title');
const editErrorEl = document.getElementById('edit-error');
const deleteBtn = document.getElementById('delete-btn');

const dateTime = new Intl.DateTimeFormat(undefined, {
  dateStyle: 'medium',
  timeStyle: 'short',
});

const state = {
  tasks: [],
  selectedID: '',
};

function setStatus(message) {
  statusEl.textContent = message;
}

function clearFieldErrors() {
  createErrorEl.textContent = '';
  editErrorEl.textContent = '';
}

function setFieldError(kind, message) {
  if (kind === 'create') {
    createErrorEl.textContent = message;
  }
  if (kind === 'edit') {
    editErrorEl.textContent = message;
  }
}

function setRuntimeState(message, mode) {
  runtimeEl.textContent = message;
  runtimeEl.classList.remove('is-ok', 'is-warn', 'is-error');
  if (mode) {
    runtimeEl.classList.add(mode);
  }
}

function selectedTask() {
  return state.tasks.find((task) => task.id === state.selectedID);
}

function statusLabel(value) {
  const labels = {
    open: 'Open',
    in_progress: 'In Progress',
    done: 'Done',
    archived: 'Archived',
  };
  return labels[value] || 'Open';
}

function formatUpdatedAt(task) {
  if (!task.updated_at) {
    return 'Updated Recently';
  }
  const parsed = new Date(task.updated_at);
  if (Number.isNaN(parsed.getTime())) {
    return 'Updated Recently';
  }
  return `Updated ${dateTime.format(parsed)}`;
}

function filteredTasks() {
  const query = queryInput.value.trim().toLowerCase();
  const status = statusFilter.value;
  const sort = sortSelect.value;

  let items = state.tasks.filter((task) => {
    if (status !== 'all' && task.status !== status) {
      return false;
    }
    if (!query) {
      return true;
    }
    return task.title.toLowerCase().includes(query);
  });

  items = items.slice().sort((a, b) => {
    if (sort === 'title_asc') {
      return a.title.localeCompare(b.title);
    }
    const updatedA = Date.parse(a.updated_at || '') || 0;
    const updatedB = Date.parse(b.updated_at || '') || 0;
    if (sort === 'updated_asc') {
      return updatedA - updatedB;
    }
    return updatedB - updatedA;
  });

  return items;
}

function selectTask(task) {
  state.selectedID = task.id;
  editID.value = task.id;
  editTitle.value = task.title;
  isDirty = false;
  selectedMetaEl.textContent = `${statusLabel(task.status)} - ${formatUpdatedAt(task)}`;
  editTitle.focus();
  renderTasks();
}

function clearSelection() {
  state.selectedID = '';
  editID.value = '';
  editTitle.value = '';
  isDirty = false;
  selectedMetaEl.textContent = 'Choose a task from the list to edit details.';
  renderTasks();
}

function renderTasks() {
  const items = filteredTasks();
  taskListEl.innerHTML = '';

  if (items.length === 0) {
    if (state.tasks.length === 0) {
      setStatus('No tasks yet. Add your first task above.');
    } else {
      setStatus('No tasks match this filter.');
    }
    return;
  }

  setStatus(`${items.length} task(s) visible`);

  for (const task of items) {
    const li = document.createElement('li');
    li.className = 'task-item';
    if (task.id === state.selectedID) {
      li.classList.add('is-selected');
    }

    const top = document.createElement('div');
    top.className = 'row';

    const title = document.createElement('p');
    title.className = 'task-title grow';
    title.textContent = task.title;

    const selectBtn = document.createElement('button');
    selectBtn.className = 'task-select';
    selectBtn.type = 'button';
    selectBtn.textContent = 'Edit';
    selectBtn.addEventListener('click', () => selectTask(task));

    top.append(title, selectBtn);

    const meta = document.createElement('p');
    meta.className = 'task-meta';
    meta.textContent = `${statusLabel(task.status)} - ${formatUpdatedAt(task)}`;

    li.append(top, meta);
    taskListEl.appendChild(li);
  }
}

function setBusy(button, busyText, fn) {
  const originalText = button.textContent;
  button.disabled = true;
  button.textContent = busyText;
  return fn().finally(() => {
    button.disabled = false;
    button.textContent = originalText;
  });
}

function mapAPIError(err) {
  const message = (err && err.message) || 'Request failed. Try again.';
  if (message.includes('title')) {
    return 'Provide a task title and try again.';
  }
  if (message.includes('not found')) {
    return 'This task no longer exists. Refresh and try another task.';
  }
  return message;
}

let isDirty = false;

editTitle.addEventListener('input', () => {
  isDirty = Boolean(selectedTask() && editTitle.value.trim() !== selectedTask().title);
});

window.addEventListener('beforeunload', (event) => {
  if (!isDirty) {
    return;
  }
  event.preventDefault();
  event.returnValue = '';
});

async function api(path, options = {}) {
  const res = await fetch(path, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  });

  if (res.status === 204) {
    return null;
  }

  let payload;
  try {
    payload = await res.json();
  } catch {
    payload = null;
  }

  if (!res.ok) {
    const message = payload?.error?.message || `Request failed (${res.status})`;
    throw new Error(message);
  }

  return payload;
}

async function loadTasks() {
  setStatus('Loading…');
  const payload = await api('/v1/tasks');
  state.tasks = payload.items || [];
  if (state.selectedID && !selectedTask()) {
    clearSelection();
  }
  renderTasks();
}

async function loadRuntime() {
  try {
    await api('/healthz', { headers: {} });
    const adapters = await api('/v1/adapters', { headers: {} });
    const issues = (adapters.errors || []).length;
    if (issues > 0 || adapters.ready === false) {
      setRuntimeState('Server Online, Adapters Need Attention', 'is-warn');
      return;
    }
    setRuntimeState('Server Online, Adapters Ready', 'is-ok');
  } catch {
    setRuntimeState('Server Unavailable. Check Local Service.', 'is-error');
  }
}

createForm.addEventListener('submit', async (event) => {
  event.preventDefault();
  clearFieldErrors();
  const title = createTitle.value.trim();
  if (!title) {
    setFieldError('create', 'Provide a task title and try again.');
    setStatus('Provide a task title and try again.');
    return;
  }

  await setBusy(createSubmit, 'Adding…', async () => {
    try {
      const task = await api('/v1/tasks', {
        method: 'POST',
        body: JSON.stringify({ title }),
      });
      state.tasks = [task, ...state.tasks];
      createTitle.value = '';
      setStatus('Task added.');
      renderTasks();
    } catch (err) {
      const message = mapAPIError(err);
      setFieldError('create', message);
      setStatus(message);
    }
  });
});

editForm.addEventListener('submit', async (event) => {
  event.preventDefault();
  clearFieldErrors();
  const current = selectedTask();
  if (!current) {
    setStatus('Select a task from the list first.');
    return;
  }

  const title = editTitle.value.trim();
  if (!title) {
    setFieldError('edit', 'Provide a task title and try again.');
    setStatus('Provide a task title and try again.');
    return;
  }

  await setBusy(editSubmit, 'Saving…', async () => {
    try {
      const updated = await api(`/v1/tasks/${current.id}`, {
        method: 'PATCH',
        body: JSON.stringify({ title }),
      });
      state.tasks = state.tasks.map((task) => (task.id === updated.id ? updated : task));
      setStatus('Task saved.');
      isDirty = false;
      selectTask(updated);
    } catch (err) {
      const message = mapAPIError(err);
      setFieldError('edit', message);
      setStatus(message);
    }
  });
});

deleteBtn.addEventListener('click', async () => {
  const current = selectedTask();
  if (!current) {
    setStatus('Select a task from the list first.');
    return;
  }

  const ok = window.confirm(`Delete "${current.title}"? This action cannot be undone.`);
  if (!ok) {
    return;
  }

  await setBusy(deleteBtn, 'Deleting…', async () => {
    try {
      await api(`/v1/tasks/${current.id}`, { method: 'DELETE' });
      state.tasks = state.tasks.filter((task) => task.id !== current.id);
      clearSelection();
      isDirty = false;
      setStatus('Task deleted.');
    } catch (err) {
      setStatus(mapAPIError(err));
    }
  });
});

filterForm.addEventListener('input', () => {
  renderTasks();
});

Promise.all([loadTasks(), loadRuntime()])
  .then(() => {
    window.setInterval(loadRuntime, 30000);
  })
  .catch((err) => {
    setStatus(mapAPIError(err));
  });

// Live updates via SSE — reload task list on any server-side mutation.
// This keeps the web UI in sync with TUI, CLI, and agent writes without polling.
(function connectEvents() {
  const es = new EventSource('/v1/tasks/events');
  es.addEventListener('task.created', () => loadTasks());
  es.addEventListener('task.updated', () => loadTasks());
  es.addEventListener('task.deleted', () => loadTasks());
  es.addEventListener('task.status_changed', () => loadTasks());
  es.addEventListener('error', () => {
    es.close();
    // Reconnect after 3s — handles server restart or network blip.
    window.setTimeout(connectEvents, 3000);
  });
}());
