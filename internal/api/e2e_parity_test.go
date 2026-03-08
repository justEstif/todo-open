package api_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/justEstif/todo-open/internal/adapters"
	"github.com/justEstif/todo-open/internal/api"
	"github.com/justEstif/todo-open/internal/core"
	"github.com/justEstif/todo-open/internal/events"
	"github.com/justEstif/todo-open/internal/store/memory"
)

// newE2EServer sets up a real in-memory server using a sequential ID generator.
func newE2EServer(t *testing.T, ids ...string) (*httptest.Server, *core.Service, *memory.TaskRepo) {
	t.Helper()
	repo := memory.NewTaskRepo()
	i := 0
	idFn := func() string {
		if i < len(ids) {
			id := ids[i]
			i++
			return id
		}
		i++
		return "task_extra"
	}
	svc := core.NewService(repo, time.Now, idFn)
	ts := httptest.NewServer(api.NewRouter(svc, adapters.Runtime{}, events.NewBroker()))
	t.Cleanup(ts.Close)
	return ts, svc, repo
}

// TestE2E_HumanFlow: create → PATCH status → complete.
func TestE2E_HumanFlow(t *testing.T) {
	t.Parallel()
	ts, _, _ := newE2EServer(t, "human-1")

	// 1. Create task.
	created := doJSON(t, ts.URL, http.MethodPost, "/v1/tasks",
		map[string]string{"title": "Human task"}, http.StatusCreated)
	id := created["id"].(string)
	if created["status"] != "open" {
		t.Fatalf("expected open, got %v", created["status"])
	}

	// 2. Update title via PATCH.
	updated := doJSON(t, ts.URL, http.MethodPatch, "/v1/tasks/"+id,
		map[string]string{"title": "Human task updated"}, http.StatusOK)
	if updated["title"] != "Human task updated" {
		t.Fatalf("unexpected title: %v", updated["title"])
	}
	if updated["version"].(float64) != 2 {
		t.Fatalf("expected version 2, got %v", updated["version"])
	}

	// 3. Complete.
	done := doJSON(t, ts.URL, http.MethodPost, "/v1/tasks/"+id+"/complete", nil, http.StatusOK)
	if done["status"] != "done" {
		t.Fatalf("expected done, got %v", done["status"])
	}
	if done["completed_at"] == nil {
		t.Fatal("expected completed_at to be set")
	}
}

// TestE2E_AgentFlow: create → appear in /next → claim → heartbeat → complete.
func TestE2E_AgentFlow(t *testing.T) {
	t.Parallel()
	ts, _, _ := newE2EServer(t, "agent-1")
	agentID := "test-agent"

	// 1. Create task.
	doJSON(t, ts.URL, http.MethodPost, "/v1/tasks",
		map[string]string{"title": "Agent task"}, http.StatusCreated)

	// 2. Appears in /next.
	next := doJSON(t, ts.URL, http.MethodGet, "/v1/tasks/next", nil, http.StatusOK)
	if next["id"] != "agent-1" {
		t.Fatalf("expected agent-1 in next, got %v", next["id"])
	}

	// 3. Claim.
	claimed := doJSON(t, ts.URL, http.MethodPost, "/v1/tasks/agent-1/claim",
		map[string]any{"agent_id": agentID, "lease_ttl_seconds": 60}, http.StatusOK)
	if claimed["status"] != "in_progress" {
		t.Fatalf("expected in_progress, got %v", claimed["status"])
	}
	ext := claimed["ext"].(map[string]any)
	agentExt := ext["agent"].(map[string]any)
	if agentExt["id"] != agentID {
		t.Fatalf("unexpected agent id: %v", agentExt["id"])
	}

	// 4. Heartbeat.
	heartbeated := doJSON(t, ts.URL, http.MethodPost, "/v1/tasks/agent-1/heartbeat",
		map[string]any{"agent_id": agentID}, http.StatusOK)
	if heartbeated["status"] != "in_progress" {
		t.Fatalf("expected in_progress after heartbeat, got %v", heartbeated["status"])
	}

	// 5. Complete via /complete.
	done := doJSON(t, ts.URL, http.MethodPost, "/v1/tasks/agent-1/complete", nil, http.StatusOK)
	if done["status"] != "done" {
		t.Fatalf("expected done, got %v", done["status"])
	}

	// 6. /next should now return 404.
	r := doReq(t, ts.URL, http.MethodGet, "/v1/tasks/next", nil, nil)
	r.Body.Close()
	if r.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 from /next after completion, got %d", r.StatusCode)
	}
}

// TestE2E_DependencyFlow: create A and B (trigger_ids=[A]) → complete A → B becomes open.
func TestE2E_DependencyFlow(t *testing.T) {
	t.Parallel()
	ts, _, _ := newE2EServer(t, "dep-a", "dep-b")

	// 1. Create task A.
	a := doJSON(t, ts.URL, http.MethodPost, "/v1/tasks",
		map[string]string{"title": "Task A"}, http.StatusCreated)
	aID := a["id"].(string)

	// 2. Create task B with trigger_ids=[A].
	bResp := doReq(t, ts.URL, http.MethodPost, "/v1/tasks", map[string]any{
		"title":       "Task B",
		"trigger_ids": []string{aID},
	}, nil)
	b := decodeBody(t, bResp)
	if bResp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 for B, got %d (body: %v)", bResp.StatusCode, b)
	}
	bID := b["id"].(string)

	// B should be pending.
	bGet := doJSON(t, ts.URL, http.MethodGet, "/v1/tasks/"+bID, nil, http.StatusOK)
	if bGet["status"] != "pending" {
		t.Fatalf("expected B to be pending, got %v", bGet["status"])
	}

	// 3. Complete A.
	doJSON(t, ts.URL, http.MethodPost, "/v1/tasks/"+aID+"/complete", nil, http.StatusOK)

	// 4. B should now be open.
	bAfter := doJSON(t, ts.URL, http.MethodGet, "/v1/tasks/"+bID, nil, http.StatusOK)
	if bAfter["status"] != "open" {
		t.Fatalf("expected B to be open after A completed, got %v", bAfter["status"])
	}
}

// TestE2E_ConcurrentClaim: two goroutines race to claim same task → exactly one gets 200, other gets 409.
func TestE2E_ConcurrentClaim(t *testing.T) {
	t.Parallel()
	ts, _, _ := newE2EServer(t, "race-task")

	doJSON(t, ts.URL, http.MethodPost, "/v1/tasks",
		map[string]string{"title": "race task"}, http.StatusCreated)

	var (
		successes int32
		conflicts int32
		wg        sync.WaitGroup
	)

	for i := range 2 {
		wg.Add(1)
		go func(agentID string) {
			defer wg.Done()
			resp := doReq(t, ts.URL, http.MethodPost, "/v1/tasks/race-task/claim",
				map[string]any{"agent_id": agentID, "lease_ttl_seconds": 60}, nil)
			resp.Body.Close()
			switch resp.StatusCode {
			case http.StatusOK:
				atomic.AddInt32(&successes, 1)
			case http.StatusConflict:
				atomic.AddInt32(&conflicts, 1)
			default:
				t.Errorf("unexpected status %d from claim", resp.StatusCode)
			}
		}("agent-" + string(rune('A'+i)))
	}
	wg.Wait()

	if successes != 1 {
		t.Fatalf("expected exactly 1 successful claim, got %d", successes)
	}
	if conflicts != 1 {
		t.Fatalf("expected exactly 1 conflict, got %d", conflicts)
	}
}

// TestE2E_LeaseExpiry: claim task, manually expire lease via service, sweeper returns task to open.
func TestE2E_LeaseExpiry(t *testing.T) {
	t.Parallel()
	repo := memory.NewTaskRepo()
	i := 0
	ids := []string{"sweep-1"}
	// Use a controllable clock.
	now := time.Now().UTC()
	nowFn := func() time.Time { return now }
	svc := core.NewService(repo, nowFn, func() string {
		if i < len(ids) {
			id := ids[i]
			i++
			return id
		}
		return "extra"
	})
	ts := httptest.NewServer(api.NewRouter(svc, adapters.Runtime{}, events.NewBroker()))
	t.Cleanup(ts.Close)

	// Create and claim task with 1s TTL.
	doJSON(t, ts.URL, http.MethodPost, "/v1/tasks",
		map[string]string{"title": "sweepable"}, http.StatusCreated)
	doJSON(t, ts.URL, http.MethodPost, "/v1/tasks/sweep-1/claim",
		map[string]any{"agent_id": "expiring-agent", "lease_ttl_seconds": 1}, http.StatusOK)

	// Verify in_progress.
	before := doJSON(t, ts.URL, http.MethodGet, "/v1/tasks/sweep-1", nil, http.StatusOK)
	if before["status"] != "in_progress" {
		t.Fatalf("expected in_progress, got %v", before["status"])
	}

	// Advance clock past TTL and run sweeper directly on service.
	now = now.Add(10 * time.Second)
	swept, err := svc.SweepExpiredLeases(context.Background())
	if err != nil {
		t.Fatalf("sweep error: %v", err)
	}
	if swept != 1 {
		t.Fatalf("expected 1 swept, got %d", swept)
	}

	// Task should be back to open.
	after := doJSON(t, ts.URL, http.MethodGet, "/v1/tasks/sweep-1", nil, http.StatusOK)
	if after["status"] != "open" {
		t.Fatalf("expected open after sweep, got %v", after["status"])
	}
	// ext.agent should be cleared.
	if after["ext"] != nil {
		ext, ok := after["ext"].(map[string]any)
		if ok && ext["agent"] != nil {
			t.Fatalf("expected ext.agent cleared after sweep, got %v", ext["agent"])
		}
	}
}
