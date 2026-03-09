package api

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/justEstif/todo-open/internal/core"
	"github.com/justEstif/todo-open/internal/events"
)

// TaskEvent is an alias for events.Event, the canonical task mutation event type.
type TaskEvent = events.Event

type Client struct {
	baseURL    string
	http       *http.Client
	httpStream *http.Client // no timeout — for long-lived SSE connections
}

func New(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		http: &http.Client{
			Timeout: 3 * time.Second,
		},
		httpStream: &http.Client{
			// No timeout: SSE connections must stay open indefinitely.
			// The caller's context controls cancellation.
		},
	}
}

func (c *Client) Health() error {
	resp, err := c.http.Get(c.baseURL + "/healthz")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}
	return nil
}

func (c *Client) CreateTask(title string) (core.Task, error) {
	var out core.Task
	err := c.do(http.MethodPost, "/v1/tasks", map[string]string{"title": title}, http.StatusCreated, &out)
	return out, err
}

func (c *Client) ListTasks() ([]core.Task, error) {
	var out struct {
		Items []core.Task `json:"items"`
	}
	if err := c.do(http.MethodGet, "/v1/tasks", nil, http.StatusOK, &out); err != nil {
		return nil, err
	}
	return out.Items, nil
}

func (c *Client) GetTask(id string) (core.Task, error) {
	var out core.Task
	err := c.do(http.MethodGet, "/v1/tasks/"+id, nil, http.StatusOK, &out)
	return out, err
}

func (c *Client) UpdateTask(id string, title string) (core.Task, error) {
	var out core.Task
	err := c.do(http.MethodPatch, "/v1/tasks/"+id, map[string]string{"title": title}, http.StatusOK, &out)
	return out, err
}

func (c *Client) CompleteTask(id string) (core.Task, error) {
	var out core.Task
	err := c.do(http.MethodPost, "/v1/tasks/"+id+"/complete", struct{}{}, http.StatusOK, &out)
	return out, err
}

func (c *Client) PatchTaskStatus(id, status string) (core.Task, error) {
	var out core.Task
	err := c.do(http.MethodPatch, "/v1/tasks/"+id, map[string]string{"status": status}, http.StatusOK, &out)
	return out, err
}

func (c *Client) DeleteTask(id string) error {
	return c.do(http.MethodDelete, "/v1/tasks/"+id, nil, http.StatusNoContent, nil)
}

// SubscribeEvents connects to the SSE event stream and returns a channel of
// TaskEvents. The caller must call the returned cancel func to stop the stream
// and release the underlying connection. Events are dropped if the channel is
// full (buffer 64); the caller should read promptly.
func (c *Client) SubscribeEvents(ctx context.Context) (<-chan TaskEvent, func(), error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/v1/tasks/events", nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")

	resp, err := c.httpStream.Do(req)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	ch := make(chan TaskEvent, 64)
	cancel := func() {
		resp.Body.Close()
	}

	go func() {
		defer close(ch)
		scanner := bufio.NewScanner(resp.Body)
		var dataLine string
		for scanner.Scan() {
			line := scanner.Text()
			if after, ok := strings.CutPrefix(line, "data: "); ok {
				dataLine = after
			} else if line == "" && dataLine != "" {
				var e TaskEvent
				if err := json.Unmarshal([]byte(dataLine), &e); err == nil {
					select {
					case ch <- e:
					default:
					}
				}
				dataLine = ""
			}
		}
	}()

	return ch, cancel, nil
}

// do executes an API request, checks the expected status, and decodes the response body into dst.
// If dst is nil, no body decoding is performed.
func (c *Client) do(method, path string, payload any, expectStatus int, dst any) error {
	resp, err := c.doJSON(method, path, payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != expectStatus {
		return decodeAPIError(resp)
	}
	if dst != nil {
		return json.NewDecoder(resp.Body).Decode(dst)
	}
	return nil
}

func (c *Client) doJSON(method string, path string, payload any) (*http.Response, error) {
	var body io.Reader
	if payload != nil {
		var b bytes.Buffer
		if err := json.NewEncoder(&b).Encode(payload); err != nil {
			return nil, err
		}
		body = &b
	}

	req, err := http.NewRequest(method, c.baseURL+path, body)
	if err != nil {
		return nil, err
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return c.http.Do(req)
}

func decodeAPIError(resp *http.Response) error {
	var out struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err == nil && out.Error.Message != "" {
		return fmt.Errorf("api error (%s): %s", out.Error.Code, out.Error.Message)
	}
	return fmt.Errorf("unexpected status: %s", resp.Status)
}
