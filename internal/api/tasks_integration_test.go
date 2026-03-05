package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ebeyene/todo-open/internal/api"
	"github.com/ebeyene/todo-open/internal/core"
	"github.com/ebeyene/todo-open/internal/store/memory"
)

func TestTaskCRUDHappyPath(t *testing.T) {
	repo := memory.NewTaskRepo()
	svc := core.NewService(repo, func() time.Time { return time.Date(2026, 3, 5, 20, 0, 0, 0, time.UTC) }, func() string { return "task_1" })
	ts := httptest.NewServer(api.NewRouter(svc))
	t.Cleanup(ts.Close)

	created := doJSON(t, ts.URL, http.MethodPost, "/v1/tasks", map[string]string{"title": "first task"}, http.StatusCreated)
	if created["id"] != "task_1" {
		t.Fatalf("unexpected id: %v", created["id"])
	}
	if created["version"].(float64) != 1 {
		t.Fatalf("expected created version 1, got %v", created["version"])
	}

	_ = doJSON(t, ts.URL, http.MethodGet, "/v1/tasks/task_1", nil, http.StatusOK)
	list := doJSON(t, ts.URL, http.MethodGet, "/v1/tasks", nil, http.StatusOK)
	items := list["items"].([]any)
	if len(items) != 1 {
		t.Fatalf("expected 1 task, got %d", len(items))
	}

	updated := doJSON(t, ts.URL, http.MethodPatch, "/v1/tasks/task_1", map[string]string{"title": "updated"}, http.StatusOK)
	if updated["title"] != "updated" {
		t.Fatalf("unexpected title after update: %v", updated["title"])
	}
	if updated["version"].(float64) != 2 {
		t.Fatalf("expected updated version 2, got %v", updated["version"])
	}

	req, err := http.NewRequest(http.MethodDelete, ts.URL+"/v1/tasks/task_1", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204 on delete, got %d", resp.StatusCode)
	}
}

func TestTaskValidationFailures(t *testing.T) {
	repo := memory.NewTaskRepo()
	svc := core.NewService(repo, time.Now, func() string { return "task_1" })
	ts := httptest.NewServer(api.NewRouter(svc))
	t.Cleanup(ts.Close)

	badCreate := doJSON(t, ts.URL, http.MethodPost, "/v1/tasks", map[string]string{"title": "  "}, http.StatusBadRequest)
	errObj := badCreate["error"].(map[string]any)
	if errObj["code"] != "validation_error" {
		t.Fatalf("unexpected error code: %v", errObj["code"])
	}

	notFound := doJSON(t, ts.URL, http.MethodGet, "/v1/tasks/missing", nil, http.StatusNotFound)
	notFoundErr := notFound["error"].(map[string]any)
	if notFoundErr["code"] != "not_found" {
		t.Fatalf("unexpected error code: %v", notFoundErr["code"])
	}

	unknownField := doRawJSON(t, ts.URL, http.MethodPost, "/v1/tasks", `{"title":"ok","extra":"nope"}`, http.StatusBadRequest)
	unknownFieldErr := unknownField["error"].(map[string]any)
	if unknownFieldErr["code"] != "invalid_json" {
		t.Fatalf("unexpected error code: %v", unknownFieldErr["code"])
	}

	trailing := doRawJSON(t, ts.URL, http.MethodPost, "/v1/tasks", `{"title":"ok"}{"title":"again"}`, http.StatusBadRequest)
	trailingErr := trailing["error"].(map[string]any)
	if trailingErr["code"] != "invalid_json" {
		t.Fatalf("unexpected error code: %v", trailingErr["code"])
	}
}

func doJSON(t *testing.T, baseURL string, method string, path string, body any, wantStatus int) map[string]any {
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

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != wantStatus {
		t.Fatalf("%s %s: expected %d, got %d", method, path, wantStatus, resp.StatusCode)
	}

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	var out map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	return out
}

func doRawJSON(t *testing.T, baseURL string, method string, path string, raw string, wantStatus int) map[string]any {
	t.Helper()

	req, err := http.NewRequest(method, baseURL+path, bytes.NewBufferString(raw))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != wantStatus {
		t.Fatalf("%s %s: expected %d, got %d", method, path, wantStatus, resp.StatusCode)
	}

	var out map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	return out
}
