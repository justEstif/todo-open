# Human UX Invariants

This document codifies how the system must behave from a human user's perspective. It serves as the north star for CLI and TUI implementation, ensuring humans always feel in control of their task list while collaborating with software agents.

---

## 1. Core Principle

**Humans set intent, agents execute.** The system must never make a human feel like they've lost control of their own task list. Every interaction should reinforce that humans are the ultimate authority over what work matters, its priority, and its completion status.

---

## 2. Workflow Inventory

### Create a Task (without dependencies)
- **Trigger**: Human needs to add new work to the system
- **Steps**: 
  1. Execute `todoopen create "Task title"`
  2. Optionally add flags for priority, project, tags, description
  3. System generates task ID and validates input
- **Expected feedback**:
  - Immediate confirmation showing task ID: "Created task abc123: Build login page"
  - Task appears in list views with correct metadata
- **Failure modes**:
  - Invalid title: clear error "Title cannot be empty"
  - Invalid metadata: specific field validation error
  - Server unavailable: local cache/write-through with clear "offline mode" indication

### Create a Task (with dependencies)
- **Trigger**: Human needs to create work that should wait for other tasks
- **Steps**:
  1. Execute `todoopen create "Dependent task" --depends-on task1,task2`
  2. System creates task with `trigger_ids` pointing to prerequisites
  3. Task initially shows `pending` status
- **Expected feedback**:
  - Confirmation: "Created task abc123: Dependent task (waiting on 2 prerequisites)"
  - Task appears in list with `pending` status and dependency info
- **Failure modes**:
  - Invalid dependency IDs: "Dependency task xyz not found"
  - Circular dependencies: "Circular dependency detected"
  - Dependency already completed: task created directly as `open`

### View Task List
- **Trigger**: Human wants to see current work status
- **Steps**:
  1. Execute `todoopen list` (or `todoopen ls`)
  2. Optionally apply filters: `todoopen list --status open --project my-project`
  3. System returns filtered, sorted task list
- **Expected feedback**:
  - Consistent tabular format with key columns
  - Stable sort order (same filter = same order)
  - Clear visual distinction between statuses
  - Agent claim state visible in appropriate column
- **Failure modes**:
  - No tasks match filter: "No tasks found matching criteria" (not empty output)
  - Invalid filter: "Invalid status value: invalid_status"
  - Server error: clear error message with retry suggestion

### Inspect a Single Task
- **Trigger**: Human wants detailed information about a specific task
- **Steps**:
  1. Execute `todoopen show task-abc123`
  2. System fetches full task object including metadata
  3. Format displays all fields, agent claim state, and dependencies
- **Expected feedback**:
  - Detailed view with all populated fields
  - Agent claim section if task is claimed: "Held by agent-xyz (expires in 5m)"
  - Dependency relationships clearly shown
  - Timestamps in human-readable local time
- **Failure modes**:
  - Task not found: "Task abc123 not found"
  - Permission error: "You do not have permission to view this task"
  - Version conflict:提示重试或显示最新版本

### Edit a Task Title or Description
- **Trigger**: Human needs to clarify or correct task information
- **Steps**:
  1. Execute `todoopen edit task-abc123 --title "New title"`
  2. Or: `todoopen edit task-abc123 --description "New description"`
  3. System validates and applies changes with optimistic concurrency
- **Expected feedback**:
  - Success: "Updated task abc123: New title"
  - Shows exactly what changed
  - Version bump visible in detailed view
- **Failure modes**:
  - Task not found: "Task abc123 not found"
  - Version conflict: "Task was modified by another process. Use --force to overwrite or fetch latest version"
  - Agent claim conflict: "Cannot edit: task is held by agent-xyz. Use --force-override to take control"

### Manually Complete a Task
- **Trigger**: Human wants to mark work as done
- **Steps**:
  1. Execute `todoopen complete task-abc123`
  2. System validates task can be completed (dependencies satisfied)
  3. Transitions task to `done` with timestamp
- **Expected feedback**:
  - "Completed task abc123: Build login page"
  - Automatic unblocking of dependent tasks visible
  - Completion timestamp shown in local time
- **Failure modes**:
  - Dependencies not met: "Cannot complete: dependencies not satisfied"
  - Task already done: "Task abc123 is already complete"
  - Agent claim conflict: "Cannot complete: task is held by agent-xyz. Use --force to override"

### Cancel / Delete a Task
- **Trigger**: Human decides work is no longer needed
- **Steps**:
  1. Execute `todoopen delete task-abc123`
  2. System prompts for confirmation with task details
  3. Human confirms deletion
  4. System removes or soft-deletes task
- **Expected feedback**:
  - Confirmation prompt: "Delete task abc123 'Build login page'? [y/N]"
  - Success: "Deleted task abc123: Build login page"
  - Any dependent tasks now show missing dependency
- **Failure modes**:
  - Agent holds task: "Cannot delete: task is actively held by agent-xyz"
  - Task has dependents: "Cannot delete: task has 2 dependent tasks. Use --force to remove anyway"
  - Confirmation refused: "Deletion cancelled"

### Prioritize (Set Priority or Reorder)
- **Trigger**: Human needs to adjust work importance
- **Steps**:
  1. Execute `todoopen prioritize task-abc123 --priority high`
  2. Or: `todoopen reorder task-abc123 --before task-def456`
  3. System validates and applies priority changes
- **Expected feedback**:
  - "Set priority of task abc123 to high"
  - List views reflect new priority order immediately
  - Agent work queue respects new priority on next poll
- **Failure modes**:
  - Invalid priority: "Priority must be one of: low, normal, high, critical"
  - Agent holds task: "Cannot prioritize: task is held by agent-xyz"
  - Version conflict: "Task was modified. Fetch latest version and retry"

### Watch Live Progress (SSE Stream)
- **Trigger**: Human wants real-time visibility into agent activity
- **Steps**:
  1. Execute `todoopen watch`
  2. System connects to SSE stream and displays events
  3. Human sees real-time task status changes
- **Expected feedback**:
  - Live updates showing agent claims, heartbeats, completions
  - Clear formatting: "agent-xyz claimed task-abc123 (expires in 5m)"
  - Stream remains stable with reconnect capabilities
- **Failure modes**:
  - Connection lost: "Connection lost, attempting to reconnect..."
  - Server unavailable: "Cannot connect to server. Check server status."
  - Permission error: "You do not have permission to access the event stream"

### Resolve a Conflict (Human Edits Claimed Task)
- **Trigger**: Human tries to edit a task an agent has claimed
- **Steps**:
  1. Human attempts edit: `todoopen edit task-abc123 --title "New title"`
  2. System detects agent claim and blocks edit
  3. Human chooses resolution strategy
- **Expected feedback**:
  - Clear error: "Cannot edit task abc123: held by agent-xyz (lease expires in 3m)"
  - Options presented: "Wait for agent to complete, force override, or cancel"
  - If forced: "Overrode agent-xyz claim on task abc123"
- **Failure modes**:
  - Agent heartbeat during edit: "Agent heartbeated during edit. Please retry."
  - Lease expired: "Agent lease expired. Task is now available."
  - Network timeout during resolution: "Resolution timeout. Task status unknown, please check."

### Review Completed Work
- **Trigger**: Human wants to audit what agents accomplished
- **Steps**:
  1. Execute `todoopen list --status done --since 1d`
  2. System returns recently completed tasks
  3. Human inspects details of specific completed tasks
- **Expected feedback**:
  - Clear list of completed tasks with completion timestamps
  - Agent attribution: "Completed by agent-xyz at 2:34 PM"
  - Audit trail visible in detailed view showing who did what when
- **Failure modes**:
  - No recent completions: "No tasks completed in the last 1 day"
  - Permission error: "You do not have permission to view completed tasks"
  - Incomplete audit trail: "Some completion details may be missing due to system issues"

---

## 3. UX Invariants (The Rules)

These are hard rules the CLI and TUI must never violate:

1. **Destructive operations require explicit confirmation** — delete, bulk status change, and force-override operations must never fire on first keypress or single flag. Always require explicit human confirmation with task details shown.

2. **Agent claim state must always be visible** — whenever a task is displayed (list, detail, edit), the agent claim state must be clearly visible including who holds it and when the lease expires.

3. **Conflict errors must show what changed** — version conflict errors must never just say "version mismatch". They must show the specific fields that changed and by whom (human vs agent).

4. **Humans can always release agent claims** — there must be a clear mechanism (with appropriate warnings) for a human to override an agent's claim and take control of a task.

5. **List output must be stable** — the same filter criteria must produce the same order of tasks every time. No random sorting or unstated ordering.

6. **Empty states must be informative** — never show blank output when no tasks match criteria. Always show a human-friendly message like "No open tasks" or "No tasks found for project 'X'".

7. **All mutations must print what changed** — silent success is never acceptable. Every successful operation must output exactly what changed in human-readable format.

8. **Timestamps in local time for humans** — all timestamps shown to humans must be converted to local time. UTC should only be used in machine output/API responses.

9. **Agent lease expiry must be clearly shown** — if an agent is heartbeating a task, the CLI/TUI must always show "held by <agent> (expires in Xm)" with a countdown.

10. **Priority changes and deletions are human-only** — the CLI should reject priority updates and delete operations if called with an agent auth header. Only humans can make these business decisions.

---

## 4. Conflict Resolution UX

### Detailed Walkthrough: Human Edits Agent-Claimed Task

When a human tries to edit a task that an agent has claimed:

```bash
$ todoopen edit task-abc123 --title "Fix the login bug"
Error: Cannot edit task abc123: held by agent-frontend-worker (lease expires in 3m 42s)

This task is currently being worked on by an agent. You have several options:

1. Wait for the agent to complete the task (estimated 3m 42s remaining)
2. Force override the agent's claim and take control
3. Cancel this edit

What would you like to do? [wait/override/cancel] 
```

**If human chooses "override":**
```
Warning: Forcing agent-frontend-worker to release task abc123.
The agent may lose unsaved work and this could disrupt workflow.

Force override task abc123 'Fix the login bug'? [y/N] y
Overrode agent-frontend-worker claim on task abc123
Updated task abc123: Fix the login bug
```

**If human chooses "wait":**
```
Waiting for agent-frontend-worker to complete task abc123...
Agent completed task abc123 (took 2m 15s)
Task abc123 is now available for editing
```

**If agent heartbeats during edit attempt:**
```
Error: Agent heartbeated during edit operation.
The task state has changed. Please fetch the latest version and retry.

Current state:
  Title: Fix the authentication bug (changed by agent)
  Status: in_progress
  Agent: agent-frontend-worker (expires in 4m 20s)
```

### Error Message Wording

Standard conflict error format:
```
Error: Cannot [operation] task [id]: [reason]

[clear explanation of what happened]
[options for resolution with exact verbs]

What would you like to do? [option1/option2/cancel]
```

Examples:
```
Error: Cannot edit task abc123: held by agent-xyz

This task is currently claimed by agent-xyz until 2:34 PM.
The agent is actively working on this task.

Options: wait (3m remaining), force-override, cancel
What would you like to do? [wait/override/cancel] 
```

```
Error: Cannot delete task abc123: has 2 dependent tasks

Deleting this task would leave its dependents in an inconsistent state.
Dependent tasks:
  - task-def456: Implement the feature
  - task-ghi789: Write tests

Options: force-delete (breaks dependents), cancel
What would you like to do? [force/cancel]
```

---

## 5. CLI Output Contract

### Task List Row Format

```
ID          STATUS     PRIORITY   AGENT           TITLE                        PROJECT
abc123      open       high       -               Fix login bug                auth
def456      in_progress  normal    agent-backend   Add user profile            frontend
ghi789      done       critical   -               Deploy to production        infra
```

**Columns and rules:**
- `ID`: First 8 characters of task ID, truncated if longer
- `STATUS`: Right-aligned, color-coded
- `PRIORITY`: Color-coded (critical=red, high=yellow, normal=green, low=gray)
- `AGENT`: "-" if no agent, otherwise agent ID truncated to 12 chars
- `TITLE`: Truncated to 30 characters with "..." if longer
- `PROJECT`: Truncated to 12 characters with "..." if longer

### Single Task Detail View

Field order and display rules:
```
Task: abc123
Title: Build the login page
Status: in_progress
Priority: high
Project: authentication
Tags: [frontend, security]
Description: Implement OAuth2 login with GitHub
Created: 2026-03-08 14:23:45 (2 hours ago)
Updated: 2026-03-08 16:15:30 (5 minutes ago)

Dependencies:
  Waiting for: setup-database (done)
  Blocking: implement-profile (pending)

Agent Claim:
  Held by: agent-frontend-worker
  Claimed: 2026-03-08 16:10:15
  Expires: 2026-03-08 16:40:15 (in 25 minutes)
  Last heartbeat: 2026-03-08 16:35:15
```

**Display rules:**
- Always show: ID, Title, Status, Priority, timestamps
- Show only if non-empty: Project, Tags, Description, Dependencies, Agent Claim
- Omit entirely if empty: started_at, completed_at (unless task is done)
- Timestamps: Local time with relative time in parentheses
- Agent claim: Special formatting with expiration countdown

### Success Messages

Pattern: `[verb] [what]` with task details

Examples:
```
Created task abc123: Build login page
Updated task abc123: title changed to "Fix login bug"
Completed task abc123: Build login page
Deleted task abc123: Build login page
Set priority of task abc123 to critical
Overrode agent-xyz claim on task abc123
```

### Error Messages

Pattern: `Error: [what failed]` followed by helpful hint

Examples:
```
Error: Task abc123 not found
Hint: Use 'todoopen list' to see available tasks
```
```
Error: Cannot delete task abc123: held by agent-xyz
Hint: Use --force-override to take control, or wait for agent to complete
```

Always send to stderr, never stdout.

### Confirmation Prompts

Exact wording pattern:
```
[Action] task [id] '[title]'? [y/N]
```

Examples:
```
Delete task abc123 'Build login page'? [y/N]
Force override agent-xyz on task abc123 'Build login page'? [y/N]
```

Default to "N" (no) - human must explicitly type "y" to proceed.

---

## 6. What the TUI Must Implement

The TUI must have these capabilities to satisfy the UX invariants:

- [ ] **Task list view with stable sorting** - Same filter always produces same order
- [ ] **Agent claim visibility in all views** - Always show who holds what and when it expires
- [ ] **Real-time status updates via SSE** - Live updates of agent activity
- [ ] **Confirmation dialogs for destructive actions** - Never delete without explicit confirmation
- [ ] **Conflict resolution interface** - Clear options when editing agent-claimed tasks
- [ ] **Detailed task view with all metadata** - Complete task information with local timestamps
- [ ] **Empty state messages** - Informative messages when no tasks match criteria
- [ ] **Change feedback display** - Clear indication of what changed after every operation
- [ ] **Agent lease expiry countdowns** - Visual countdown showing time remaining on claims
- [ ] **Priority and deletion protection** - Block agent attempts to change priority or delete
- [ ] **Error display with hints** - Clear error messages and next-step suggestions
- [ ] **Filtering and search** - By status, project, tags, and free text
- [ ] **Dependency visualization** - Show blocking/blocked-by relationships
- [ ] **Audit trail for completed work** - Show who completed what and when
- [ ] **Local timestamp display** - All human-visible times in local timezone
- [ ] **Keyboard shortcuts for common actions** - Efficient navigation and operation

---

## 7. Cross-References

- [Agent Primitives](./agent-primitives.md) - Technical contract for agent behavior and endpoints
- [Human vs. Agent Model](./human-agent-model.md) - Mental model and responsibility boundaries
- [API Documentation](./api.md) - Full HTTP API reference for all client implementations