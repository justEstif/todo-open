package core

import "context"

type TaskService interface {
	CreateTask(ctx context.Context, title string) (Task, error)
}
