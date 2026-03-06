package api_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/justEstif/todo-open/internal/adapters"
	"github.com/justEstif/todo-open/internal/api"
	"github.com/justEstif/todo-open/internal/core"
	"github.com/justEstif/todo-open/internal/store/memory"
)

func TestWebAppStaticSurface(t *testing.T) {
	repo := memory.NewTaskRepo()
	svc := core.NewService(repo, time.Now, func() string { return "task_1" })
	ts := httptest.NewServer(api.NewRouter(svc, adapters.Runtime{}))
	t.Cleanup(ts.Close)

	resp := mustGet(t, ts.URL+"/")
	if got := resp.StatusCode; got != http.StatusOK {
		t.Fatalf("GET / status=%d want=%d", got, http.StatusOK)
	}
	html := readBody(t, resp)
	if !strings.Contains(html, "todo.open") {
		t.Fatalf("index missing app title")
	}
	if !strings.Contains(html, "/static/app.css") || !strings.Contains(html, "/static/app.js") {
		t.Fatalf("index missing static asset references")
	}
	if !strings.Contains(html, "runtime-status") {
		t.Fatalf("index missing runtime status element")
	}

	resp = mustGet(t, ts.URL+"/static/app.css")
	if got := resp.StatusCode; got != http.StatusOK {
		t.Fatalf("GET /static/app.css status=%d want=%d", got, http.StatusOK)
	}
	css := readBody(t, resp)
	if !strings.Contains(css, "--brand") {
		t.Fatalf("app.css appears incomplete")
	}

	resp = mustGet(t, ts.URL+"/static/app.js")
	if got := resp.StatusCode; got != http.StatusOK {
		t.Fatalf("GET /static/app.js status=%d want=%d", got, http.StatusOK)
	}
	js := readBody(t, resp)
	if !strings.Contains(js, "loadTasks") {
		t.Fatalf("app.js missing expected function")
	}
}

func mustGet(t *testing.T, url string) *http.Response {
	t.Helper()
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	return resp
}

func readBody(t *testing.T, resp *http.Response) string {
	t.Helper()
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
