package core

import "time"

type TaskStatus string

const (
	TaskStatusOpen     TaskStatus = "open"
	TaskStatusArchived TaskStatus = "archived"
)

type Task struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Status    TaskStatus `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
