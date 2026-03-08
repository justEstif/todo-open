package events

import (
	"context"
	"time"

	"github.com/justEstif/todo-open/internal/core"
)

// EventEmittingService wraps a core.TaskService and publishes events after
// each successful mutation. Read-only methods are delegated without emission.
type EventEmittingService struct {
	inner  core.TaskService
	broker *Broker
	nowFn  func() time.Time
}

// NewEventEmittingService returns a new EventEmittingService.
func NewEventEmittingService(inner core.TaskService, broker *Broker) *EventEmittingService {
	return &EventEmittingService{inner: inner, broker: broker, nowFn: time.Now}
}

func (s *EventEmittingService) CreateTask(ctx context.Context, title string) (core.Task, error) {
	task, err := s.inner.CreateTask(ctx, title)
	if err != nil {
		return task, err
	}
	s.broker.Publish(Event{Type: TypeCreated, Task: &task, At: s.nowFn().UTC()})
	return task, nil
}

func (s *EventEmittingService) GetTask(ctx context.Context, id string) (core.Task, error) {
	return s.inner.GetTask(ctx, id)
}

func (s *EventEmittingService) ListTasks(ctx context.Context, filter core.ListFilter) ([]core.Task, error) {
	return s.inner.ListTasks(ctx, filter)
}

func (s *EventEmittingService) UpdateTask(ctx context.Context, id string, title string) (core.Task, error) {
	task, err := s.inner.UpdateTask(ctx, id, title)
	if err != nil {
		return task, err
	}
	s.broker.Publish(Event{Type: TypeUpdated, Task: &task, At: s.nowFn().UTC()})
	return task, nil
}

func (s *EventEmittingService) DeleteTask(ctx context.Context, id string) error {
	if err := s.inner.DeleteTask(ctx, id); err != nil {
		return err
	}
	s.broker.Publish(Event{Type: TypeDeleted, At: s.nowFn().UTC()})
	return nil
}

func (s *EventEmittingService) NextTask(ctx context.Context) (core.Task, error) {
	return s.inner.NextTask(ctx)
}

func (s *EventEmittingService) ClaimTask(ctx context.Context, id, agentID string, leaseTTLSeconds int) (core.Task, error) {
	task, err := s.inner.ClaimTask(ctx, id, agentID, leaseTTLSeconds)
	if err != nil {
		return task, err
	}
	s.broker.Publish(Event{Type: TypeUpdated, Task: &task, At: s.nowFn().UTC()})
	return task, nil
}

func (s *EventEmittingService) HeartbeatTask(ctx context.Context, id, agentID string) (core.Task, error) {
	return s.inner.HeartbeatTask(ctx, id, agentID)
}

func (s *EventEmittingService) ReleaseTask(ctx context.Context, id, agentID string) (core.Task, error) {
	task, err := s.inner.ReleaseTask(ctx, id, agentID)
	if err != nil {
		return task, err
	}
	s.broker.Publish(Event{Type: TypeUpdated, Task: &task, At: s.nowFn().UTC()})
	return task, nil
}

func (s *EventEmittingService) SweepExpiredLeases(ctx context.Context) (int, error) {
	return s.inner.SweepExpiredLeases(ctx)
}

func (s *EventEmittingService) CompleteTask(ctx context.Context, id string) (core.Task, error) {
	// capture old status before completing
	old, err := s.inner.GetTask(ctx, id)
	if err != nil {
		return core.Task{}, err
	}
	oldStatus := old.Status

	task, err := s.inner.CompleteTask(ctx, id)
	if err != nil {
		return task, err
	}
	newStatus := task.Status
	s.broker.Publish(Event{
		Type:      TypeStatusChanged,
		Task:      &task,
		OldStatus: &oldStatus,
		NewStatus: &newStatus,
		At:        s.nowFn().UTC(),
	})
	return task, nil
}
