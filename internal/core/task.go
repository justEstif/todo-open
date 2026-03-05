package core

import "time"

type TaskStatus string

type TaskPriority string

const (
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
	Ext          any          `json:"ext,omitempty"`
}
