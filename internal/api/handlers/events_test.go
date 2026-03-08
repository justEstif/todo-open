package handlers_test

import (
	"bufio"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/justEstif/todo-open/internal/api/handlers"
	"github.com/justEstif/todo-open/internal/events"
)

func TestSSEHandlerReceivesEvent(t *testing.T) {
	t.Parallel()
	b := events.NewBroker()
	h := handlers.NewEventHandler(b)

	srv := httptest.NewServer(http.HandlerFunc(h.Stream))
	defer srv.Close()

	resp, err := http.Get(srv.URL) //nolint:noctx
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if ct := resp.Header.Get("Content-Type"); ct != "text/event-stream" {
		t.Errorf("Content-Type = %q, want text/event-stream", ct)
	}

	// Publish an event after connection is established.
	go func() {
		time.Sleep(20 * time.Millisecond)
		b.Publish(events.Event{Type: events.TypeCreated, At: time.Now()})
	}()

	scanner := bufio.NewScanner(resp.Body)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
		if line == "" && len(lines) >= 3 {
			break
		}
	}

	found := false
	for _, l := range lines {
		if strings.HasPrefix(l, "event: task.created") {
			found = true
		}
	}
	if !found {
		t.Errorf("SSE frame missing event line; got lines: %v", lines)
	}
}
