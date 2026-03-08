package events_test

import (
	"context"
	"testing"
	"time"

	"github.com/justEstif/todo-open/internal/core"
	"github.com/justEstif/todo-open/internal/events"
	"github.com/justEstif/todo-open/internal/store/memory"
)

func newTestService(b *events.Broker) *events.EventEmittingService {
	repo := memory.NewTaskRepo()
	svc := core.NewService(repo, time.Now, nil)
	return events.NewEventEmittingService(svc, b)
}

func TestEventEmittingServiceCreate(t *testing.T) {
	t.Parallel()
	b := events.NewBroker()
	ch, unsub := b.Subscribe(4)
	defer unsub()

	svc := newTestService(b)
	task, err := svc.CreateTask(context.Background(), "hello")
	if err != nil {
		t.Fatal(err)
	}

	e := <-ch
	if e.Type != events.TypeCreated {
		t.Errorf("got event type %q, want %q", e.Type, events.TypeCreated)
	}
	if e.Task == nil || e.Task.ID != task.ID {
		t.Errorf("event task mismatch")
	}
}

func TestEventEmittingServiceUpdate(t *testing.T) {
	t.Parallel()
	b := events.NewBroker()
	svc := newTestService(b)

	task, _ := svc.CreateTask(context.Background(), "original")

	ch, unsub := b.Subscribe(4)
	defer unsub()

	_, err := svc.UpdateTask(context.Background(), task.ID, "updated")
	if err != nil {
		t.Fatal(err)
	}

	e := <-ch
	if e.Type != events.TypeUpdated {
		t.Errorf("got event type %q, want %q", e.Type, events.TypeUpdated)
	}
}

func TestEventEmittingServiceDelete(t *testing.T) {
	t.Parallel()
	b := events.NewBroker()
	svc := newTestService(b)

	task, _ := svc.CreateTask(context.Background(), "to delete")

	ch, unsub := b.Subscribe(4)
	defer unsub()

	if err := svc.DeleteTask(context.Background(), task.ID); err != nil {
		t.Fatal(err)
	}

	e := <-ch
	if e.Type != events.TypeDeleted {
		t.Errorf("got event type %q, want %q", e.Type, events.TypeDeleted)
	}
}

func TestEventEmittingServiceComplete(t *testing.T) {
	t.Parallel()
	b := events.NewBroker()
	svc := newTestService(b)

	task, _ := svc.CreateTask(context.Background(), "to complete")

	ch, unsub := b.Subscribe(4)
	defer unsub()

	_, err := svc.CompleteTask(context.Background(), task.ID)
	if err != nil {
		t.Fatal(err)
	}

	e := <-ch
	if e.Type != events.TypeStatusChanged {
		t.Errorf("got event type %q, want %q", e.Type, events.TypeStatusChanged)
	}
	if e.OldStatus == nil || *e.OldStatus != core.TaskStatusOpen {
		t.Errorf("expected old_status=open, got %v", e.OldStatus)
	}
	if e.NewStatus == nil || *e.NewStatus != core.TaskStatusDone {
		t.Errorf("expected new_status=done, got %v", e.NewStatus)
	}
}
