package core_test

import (
	"context"
	"testing"
	"time"

	"github.com/justEstif/todo-open/internal/core"
	"github.com/justEstif/todo-open/internal/store/memory"
)

func fixedNow() func() time.Time {
	t := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	return func() time.Time { return t }
}

func TestCompleteTaskTransitionsPending(t *testing.T) {
	t.Parallel()
	repo := memory.NewTaskRepo()
	now := fixedNow()
	ids := []string{"t1", "t2", "t3"}
	i := 0
	svc := core.NewService(repo, now, func() string {
		id := ids[i]
		i++
		return id
	})
	ctx := context.Background()

	// Create two prerequisite tasks.
	t1, err := svc.CreateTask(ctx, "prereq 1")
	if err != nil {
		t.Fatal(err)
	}
	t2, err := svc.CreateTask(ctx, "prereq 2")
	if err != nil {
		t.Fatal(err)
	}

	// Manually insert a pending task with trigger_ids.
	pending := core.Task{
		ID:         "t3",
		Title:      "dependent task",
		Status:     core.TaskStatusPending,
		TriggerIDs: []string{t1.ID, t2.ID},
		CreatedAt:  now(),
		UpdatedAt:  now(),
		Version:    1,
	}
	if _, err := repo.Create(ctx, pending); err != nil {
		t.Fatal(err)
	}

	// Complete t1 — t3 still has t2 as pending trigger.
	if _, err := svc.CompleteTask(ctx, t1.ID); err != nil {
		t.Fatal(err)
	}
	got, err := repo.GetByID(ctx, "t3")
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != core.TaskStatusPending {
		t.Fatalf("expected t3 still pending after t1 done; got %s", got.Status)
	}

	// Complete t2 — now all triggers done; t3 should become open.
	if _, err := svc.CompleteTask(ctx, t2.ID); err != nil {
		t.Fatal(err)
	}
	got, err = repo.GetByID(ctx, "t3")
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != core.TaskStatusOpen {
		t.Fatalf("expected t3 open after all triggers done; got %s", got.Status)
	}
}

func TestListTasksFilter(t *testing.T) {
	t.Parallel()
	repo := memory.NewTaskRepo()
	now := fixedNow()
	ids := []string{"t1", "t2", "t3"}
	i := 0
	svc := core.NewService(repo, now, func() string {
		id := ids[i]
		i++
		return id
	})
	ctx := context.Background()

	_, _ = svc.CreateTask(ctx, "open task")
	// Insert a pending task with blocked_by.
	pending := core.Task{
		ID: "t2", Title: "pending task", Status: core.TaskStatusPending,
		BlockedBy: []string{"t1"}, CreatedAt: now(), UpdatedAt: now(), Version: 1,
	}
	_, _ = repo.Create(ctx, pending)
	// Insert a task with no blockers.
	free := core.Task{
		ID: "t3", Title: "free open", Status: core.TaskStatusOpen,
		CreatedAt: now(), UpdatedAt: now(), Version: 1,
	}
	_, _ = repo.Create(ctx, free)

	// Filter by status=pending.
	pendingList, err := svc.ListTasks(ctx, core.ListFilter{Status: core.TaskStatusPending})
	if err != nil {
		t.Fatal(err)
	}
	if len(pendingList) != 1 || pendingList[0].ID != "t2" {
		t.Fatalf("expected 1 pending task; got %v", pendingList)
	}

	// Filter by is_blocked=true.
	blocked, err := svc.ListTasks(ctx, core.ListFilter{IsBlocked: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(blocked) != 1 || blocked[0].ID != "t2" {
		t.Fatalf("expected 1 blocked task; got %v", blocked)
	}

	// No filter returns all.
	all, err := svc.ListTasks(ctx, core.ListFilter{})
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 3 {
		t.Fatalf("expected 3 tasks; got %d", len(all))
	}
}
