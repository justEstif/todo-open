const statusEl = document.getElementById('status');
const taskListEl = document.getElementById('task-list');
const createForm = document.getElementById('create-form');
const createTitle = document.getElementById('create-title');
const editForm = document.getElementById('edit-form');
const editID = document.getElementById('edit-id');
const editTitle = document.getElementById('edit-title');
const deleteBtn = document.getElementById('delete-btn');

let tasks = [];

function setStatus(message) {
  statusEl.textContent = message;
}

function selectedTask() {
  return tasks.find((t) => t.id === editID.value);
}

function renderTasks() {
  taskListEl.innerHTML = '';
  if (tasks.length === 0) {
    setStatus('No tasks yet');
    return;
  }

  setStatus(`${tasks.length} task(s)`);
  for (const task of tasks) {
    const li = document.createElement('li');
    const top = document.createElement('div');
    top.className = 'row';

    const title = document.createElement('strong');
    title.className = 'grow';
    title.textContent = task.title;

    const edit = document.createElement('button');
    edit.type = 'button';
    edit.textContent = 'Edit';
    edit.addEventListener('click', () => {
      editID.value = task.id;
      editTitle.value = task.title;
      editTitle.focus();
    });

    top.append(title, edit);

    const meta = document.createElement('p');
    meta.className = 'muted';
    meta.textContent = `${task.id} • ${task.status}`;

    li.append(top, meta);
    taskListEl.appendChild(li);
  }
}

async function api(path, options = {}) {
  const res = await fetch(path, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  });

  if (res.status === 204) return null;
  const payload = await res.json();
  if (!res.ok) {
    const message = payload?.error?.message || `request failed (${res.status})`;
    throw new Error(message);
  }
  return payload;
}

async function loadTasks() {
  setStatus('Loading...');
  const payload = await api('/v1/tasks');
  tasks = payload.items || [];
  renderTasks();
}

createForm.addEventListener('submit', async (e) => {
  e.preventDefault();
  try {
    const title = createTitle.value.trim();
    if (!title) return;
    await api('/v1/tasks', {
      method: 'POST',
      body: JSON.stringify({ title }),
    });
    createTitle.value = '';
    await loadTasks();
  } catch (err) {
    setStatus(err.message);
  }
});

editForm.addEventListener('submit', async (e) => {
  e.preventDefault();
  try {
    const current = selectedTask();
    if (!current) {
      setStatus('Select a task to edit');
      return;
    }
    const title = editTitle.value.trim();
    await api(`/v1/tasks/${current.id}`, {
      method: 'PATCH',
      body: JSON.stringify({ title }),
    });
    await loadTasks();
    setStatus('Saved');
  } catch (err) {
    setStatus(err.message);
  }
});

deleteBtn.addEventListener('click', async () => {
  try {
    const current = selectedTask();
    if (!current) {
      setStatus('Select a task to delete');
      return;
    }
    await api(`/v1/tasks/${current.id}`, { method: 'DELETE' });
    editID.value = '';
    editTitle.value = '';
    await loadTasks();
    setStatus('Deleted');
  } catch (err) {
    setStatus(err.message);
  }
});

loadTasks().catch((err) => setStatus(err.message));
