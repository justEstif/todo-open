package core

import "time"

type TaskStatus string

type TaskPriority string

const (
	// TaskStatusPending means the task is waiting on dependencies (trigger_ids) to complete.
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusOpen       TaskStatus = "open"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusDone       TaskStatus = "done"
	TaskStatusArchived   TaskStatus = "archived"
)

const (
	TaskPriorityLow      TaskPriority = "low"
	TaskPriorityNormal   TaskPriority = "normal"
	TaskPriorityHigh     TaskPriority = "high"
	TaskPriorityCritical TaskPriority = "critical"
)

type Task struct {
	ID           string       `json:"id"`
	Title        string       `json:"title"`
	Status       TaskStatus   `json:"status"`
	Priority     TaskPriority `json:"priority,omitempty"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	StartedAt    *time.Time   `json:"started_at,omitempty"`
	CompletedAt  *time.Time   `json:"completed_at,omitempty"`
	DeletedAt    *time.Time   `json:"deleted_at,omitempty"`
	Version      int          `json:"version"`
	Description  string       `json:"description,omitempty"`
	Project      string       `json:"project,omitempty"`
	ParentID     string       `json:"parent_id,omitempty"`
	Assignee     string       `json:"assignee,omitempty"`
	EstimateMins int          `json:"estimate_minutes,omitempty"`
	SortOrder    float64      `json:"sort_order,omitempty"`
	Tags         []string     `json:"tags,omitempty"`
	// TriggerIDs lists task IDs that must all reach "done" before this task transitions from pending to open.
	TriggerIDs []string `json:"trigger_ids,omitempty"`
	// Blocking lists task IDs that this task is blocking (outgoing dependency edges).
	Blocking []string `json:"blocking,omitempty"`
	// BlockedBy lists task IDs that are blocking this task (incoming dependency edges).
	BlockedBy []string `json:"blocked_by,omitempty"`
	Ext       any      `json:"ext,omitempty"`
}
