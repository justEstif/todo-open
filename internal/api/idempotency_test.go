package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/justEstif/todo-open/internal/adapters"
	"github.com/justEstif/todo-open/internal/api"
	"github.com/justEstif/todo-open/internal/api/middleware"
	"github.com/justEstif/todo-open/internal/core"
	"github.com/justEstif/todo-open/internal/events"
	"github.com/justEstif/todo-open/internal/store/memory"
)

func newTestServer(t *testing.T, ids ...string) (*httptest.Server, *middleware.IdempotencyStore) {
	t.Helper()
	repo := memory.NewTaskRepo()
	i := 0
	idFn := func() string {
		if i < len(ids) {
			id := ids[i]
			i++
			return id
		}
		return "task_extra"
	}
	svc := core.NewService(repo, func() time.Time { return time.Now().UTC() }, idFn)
	idem := middleware.NewIdempotencyStore()
	ts := httptest.NewServer(api.NewRouterWithIdempotency(svc, adapters.Runtime{}, events.NewBroker(), idem))
	t.Cleanup(ts.Close)
	return ts, idem
}

func doReq(t *testing.T, baseURL, method, path string, body any, extraHeaders map[string]string) *http.Response {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatal(err)
		}
	}
	req, err := http.NewRequest(method, baseURL+path, &buf)
	if err != nil {
		t.Fatal(err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range extraHeaders {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	return resp
}

func decodeBody(t *testing.T, resp *http.Response) map[string]any {
	t.Helper()
	defer resp.Body.Close()
	var out map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	return out
}

// TestPutUpsert_Create verifies PUT creates a new task (201).
func TestPutUpsert_Create(t *testing.T) {
	t.Parallel()
	ts, _ := newTestServer(t)

	resp := doReq(t, ts.URL, http.MethodPut, "/v1/tasks/agent-task-1",
		map[string]string{"title": "new task"}, nil)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
	body := decodeBody(t, resp)
	if body["id"] != "agent-task-1" {
		t.Fatalf("expected id agent-task-1, got %v", body["id"])
	}
}

// TestPutUpsert_IdenticalContent verifies PUT with same content is a 200 no-op.
func TestPutUpsert_IdenticalContent(t *testing.T) {
	t.Parallel()
	ts, _ := newTestServer(t)

	// Create via PUT.
	resp1 := doReq(t, ts.URL, http.MethodPut, "/v1/tasks/idempotent-task",
		map[string]string{"title": "hello"}, nil)
	if resp1.StatusCode != http.StatusCreated {
		resp1.Body.Close()
		t.Fatalf("first PUT expected 201, got %d", resp1.StatusCode)
	}
	body1 := decodeBody(t, resp1)
	v1 := body1["version"].(float64)

	// Second PUT with identical content — no-op 200, version unchanged.
	resp2 := doReq(t, ts.URL, http.MethodPut, "/v1/tasks/idempotent-task",
		map[string]string{"title": "hello"}, nil)
	if resp2.StatusCode != http.StatusOK {
		resp2.Body.Close()
		t.Fatalf("second PUT (same content) expected 200, got %d", resp2.StatusCode)
	}
	body2 := decodeBody(t, resp2)
	if body2["version"].(float64) != v1 {
		t.Fatalf("expected version unchanged (%v), got %v", v1, body2["version"])
	}
}

// TestPutUpsert_DifferentContentNoIfMatch verifies 409 when content differs without If-Match.
func TestPutUpsert_DifferentContentNoIfMatch(t *testing.T) {
	t.Parallel()
	ts, _ := newTestServer(t)

	resp1 := doReq(t, ts.URL, http.MethodPut, "/v1/tasks/conflict-task",
		map[string]string{"title": "original"}, nil)
	if resp1.StatusCode != http.StatusCreated {
		resp1.Body.Close()
		t.Fatalf("expected 201, got %d", resp1.StatusCode)
	}
	resp1.Body.Close()

	resp2 := doReq(t, ts.URL, http.MethodPut, "/v1/tasks/conflict-task",
		map[string]string{"title": "different"}, nil)
	if resp2.StatusCode != http.StatusConflict {
		resp2.Body.Close()
		t.Fatalf("expected 409, got %d", resp2.StatusCode)
	}
	resp2.Body.Close()
}

// TestPutUpsert_UpdateWithIfMatch verifies 200 update when If-Match matches.
func TestPutUpsert_UpdateWithIfMatch(t *testing.T) {
	t.Parallel()
	ts, _ := newTestServer(t)

	resp1 := doReq(t, ts.URL, http.MethodPut, "/v1/tasks/match-task",
		map[string]string{"title": "v1"}, nil)
	if resp1.StatusCode != http.StatusCreated {
		resp1.Body.Close()
		t.Fatalf("expected 201, got %d", resp1.StatusCode)
	}
	etag := resp1.Header.Get("ETag")
	resp1.Body.Close()

	resp2 := doReq(t, ts.URL, http.MethodPut, "/v1/tasks/match-task",
		map[string]string{"title": "v2"}, map[string]string{"If-Match": etag})
	if resp2.StatusCode != http.StatusOK {
		resp2.Body.Close()
		t.Fatalf("expected 200, got %d", resp2.StatusCode)
	}
	body2 := decodeBody(t, resp2)
	if body2["title"] != "v2" {
		t.Fatalf("expected title v2, got %v", body2["title"])
	}
}

// TestPatchStatus_Idempotent verifies PATCH with same status is a no-op (version unchanged).
func TestPatchStatus_Idempotent(t *testing.T) {
	t.Parallel()
	ts, _ := newTestServer(t, "patch-task")

	doJSON(t, ts.URL, http.MethodPost, "/v1/tasks", map[string]string{"title": "task"}, http.StatusCreated)

	// First PATCH to done.
	r1 := doReq(t, ts.URL, http.MethodPatch, "/v1/tasks/patch-task",
		map[string]string{"status": "done"}, nil)
	b1 := decodeBody(t, r1)

	// Second PATCH same status — should be no-op, version unchanged.
	r2 := doReq(t, ts.URL, http.MethodPatch, "/v1/tasks/patch-task",
		map[string]string{"status": "done"}, nil)
	if r2.StatusCode != http.StatusOK {
		r2.Body.Close()
		t.Fatalf("expected 200, got %d", r2.StatusCode)
	}
	b2 := decodeBody(t, r2)
	if b1["version"] != b2["version"] {
		t.Fatalf("expected same version %v, got %v", b1["version"], b2["version"])
	}
}

// TestIdempotencyKey_Deduplication verifies X-Idempotency-Key returns cached response.
func TestIdempotencyKey_Deduplication(t *testing.T) {
	t.Parallel()
	ts, _ := newTestServer(t, "idem-task-1", "idem-task-2")

	headers := map[string]string{"X-Idempotency-Key": "create-key-1"}

	// First call — creates task.
	resp1 := doReq(t, ts.URL, http.MethodPost, "/v1/tasks",
		map[string]string{"title": "my task"}, headers)
	if resp1.StatusCode != http.StatusCreated {
		resp1.Body.Close()
		t.Fatalf("expected 201, got %d", resp1.StatusCode)
	}
	body1 := decodeBody(t, resp1)

	// Second call with same key — should return cached 201 with same body.
	resp2 := doReq(t, ts.URL, http.MethodPost, "/v1/tasks",
		map[string]string{"title": "my task"}, headers)
	if resp2.StatusCode != http.StatusCreated {
		resp2.Body.Close()
		t.Fatalf("expected cached 201, got %d", resp2.StatusCode)
	}
	body2 := decodeBody(t, resp2)
	if body1["id"] != body2["id"] {
		t.Fatalf("expected same task id %v, got %v", body1["id"], body2["id"])
	}
	if resp2.Header.Get("X-Idempotency-Replayed") != "true" {
		t.Fatal("expected X-Idempotency-Replayed header on cached response")
	}
}
