package core

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrNotFound     = errors.New("task not found")
	ErrCycleDetected = errors.New("dependency cycle detected")
)

// ListFilter defines optional filters for listing tasks.
type ListFilter struct {
	Status    TaskStatus // if non-empty, filter by status
	IsBlocked bool       // if true, return only tasks that have non-empty blocked_by
}

type TaskService interface {
	CreateTask(ctx context.Context, title string) (Task, error)
	GetTask(ctx context.Context, id string) (Task, error)
	ListTasks(ctx context.Context, filter ListFilter) ([]Task, error)
	UpdateTask(ctx context.Context, id string, title string) (Task, error)
	DeleteTask(ctx context.Context, id string) error
	CompleteTask(ctx context.Context, id string) (Task, error)
}

type IDGenerator func() string

type TaskRepository interface {
	Create(ctx context.Context, task Task) (Task, error)
	GetByID(ctx context.Context, id string) (Task, error)
	List(ctx context.Context) ([]Task, error)
	Update(ctx context.Context, task Task) (Task, error)
}

type Service struct {
	repo  TaskRepository
	nowFn func() time.Time
	idFn  IDGenerator
}

func NewService(repo TaskRepository, nowFn func() time.Time, idFn IDGenerator) *Service {
	if nowFn == nil {
		nowFn = time.Now
	}
	if idFn == nil {
		idFn = func() string { return fmt.Sprintf("task_%d", nowFn().UnixNano()) }
	}
	return &Service{repo: repo, nowFn: nowFn, idFn: idFn}
}

func (s *Service) CreateTask(ctx context.Context, title string) (Task, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return Task{}, fmt.Errorf("title is required: %w", ErrInvalidInput)
	}
	now := s.nowFn().UTC()
	task := Task{ID: s.idFn(), Title: title, Status: TaskStatusOpen, CreatedAt: now, UpdatedAt: now, Version: 1}
	return s.repo.Create(ctx, task)
}

func (s *Service) GetTask(ctx context.Context, id string) (Task, error) {
	if strings.TrimSpace(id) == "" {
		return Task{}, fmt.Errorf("id is required: %w", ErrInvalidInput)
	}
	return s.repo.GetByID(ctx, id)
}

func (s *Service) ListTasks(ctx context.Context, filter ListFilter) ([]Task, error) {
	all, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	if filter.Status == "" && !filter.IsBlocked {
		return all, nil
	}
	out := make([]Task, 0, len(all))
	for _, t := range all {
		if filter.Status != "" && t.Status != filter.Status {
			continue
		}
		if filter.IsBlocked && len(t.BlockedBy) == 0 {
			continue
		}
		out = append(out, t)
	}
	return out, nil
}

func (s *Service) UpdateTask(ctx context.Context, id string, title string) (Task, error) {
	if strings.TrimSpace(id) == "" {
		return Task{}, fmt.Errorf("id is required: %w", ErrInvalidInput)
	}
	title = strings.TrimSpace(title)
	if title == "" {
		return Task{}, fmt.Errorf("title is required: %w", ErrInvalidInput)
	}

	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return Task{}, err
	}
	task.Title = title
	task.UpdatedAt = s.nowFn().UTC()
	task.Version++
	return s.repo.Update(ctx, task)
}

func (s *Service) DeleteTask(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("id is required: %w", ErrInvalidInput)
	}
	now := s.nowFn().UTC()
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	task.DeletedAt = &now
	task.Status = TaskStatusArchived
	task.UpdatedAt = now
	task.Version++
	_, err = s.repo.Update(ctx, task)
	return err
}

// CompleteTask sets the task status to done and evaluates pending tasks whose trigger_ids are now all done.
func (s *Service) CompleteTask(ctx context.Context, id string) (Task, error) {
	if strings.TrimSpace(id) == "" {
		return Task{}, fmt.Errorf("id is required: %w", ErrInvalidInput)
	}
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return Task{}, err
	}
	now := s.nowFn().UTC()
	task.Status = TaskStatusDone
	task.CompletedAt = &now
	task.UpdatedAt = now
	task.Version++
	task, err = s.repo.Update(ctx, task)
	if err != nil {
		return Task{}, err
	}
	if err := s.evaluatePendingTasks(ctx, id); err != nil {
		return task, err
	}
	return task, nil
}

// evaluatePendingTasks checks all pending tasks whose trigger_ids include completedID.
// If all of a pending task's trigger_ids are now done, it transitions that task to open.
func (s *Service) evaluatePendingTasks(ctx context.Context, completedID string) error {
	all, err := s.repo.List(ctx)
	if err != nil {
		return err
	}
	// Build a set of done task IDs for fast lookup.
	doneIDs := map[string]bool{}
	for _, t := range all {
		if t.Status == TaskStatusDone {
			doneIDs[t.ID] = true
		}
	}
	now := s.nowFn().UTC()
	for _, t := range all {
		if t.Status != TaskStatusPending || len(t.TriggerIDs) == 0 {
			continue
		}
		// Check if this task depends on the completed task.
		dependsOnCompleted := false
		for _, trigID := range t.TriggerIDs {
			if trigID == completedID {
				dependsOnCompleted = true
				break
			}
		}
		if !dependsOnCompleted {
			continue
		}
		// Check if all trigger_ids are now done.
		allDone := true
		for _, trigID := range t.TriggerIDs {
			if !doneIDs[trigID] {
				allDone = false
				break
			}
		}
		if allDone {
			t.Status = TaskStatusOpen
			t.UpdatedAt = now
			t.Version++
			if _, err := s.repo.Update(ctx, t); err != nil {
				return err
			}
		}
	}
	return nil
}

// detectCycle performs DFS cycle detection on dependency edges.
// adj maps task ID -> its trigger_ids (edges task depends on).
func detectCycle(adj map[string][]string) bool {
	// 0=unvisited, 1=in-stack, 2=done
	state := map[string]int{}
	var dfs func(id string) bool
	dfs = func(id string) bool {
		if state[id] == 1 {
			return true
		}
		if state[id] == 2 {
			return false
		}
		state[id] = 1
		for _, dep := range adj[id] {
			if dfs(dep) {
				return true
			}
		}
		state[id] = 2
		return false
	}
	for id := range adj {
		if state[id] == 0 && dfs(id) {
			return true
		}
	}
	return false
}
