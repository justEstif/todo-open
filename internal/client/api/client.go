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
)

// TaskEvent is a task mutation event received over the SSE stream.
type TaskEvent struct {
	Type      string           `json:"type"`
	Task      *core.Task       `json:"task,omitempty"`
	OldStatus *core.TaskStatus `json:"old_status,omitempty"`
	NewStatus *core.TaskStatus `json:"new_status,omitempty"`
	At        time.Time        `json:"at"`
}

type Client struct {
	baseURL string
	http    *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		http: &http.Client{
			Timeout: 3 * time.Second,
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
	payload := map[string]string{"title": title}
	resp, err := c.doJSON(http.MethodPost, "/v1/tasks", payload)
	if err != nil {
		return core.Task{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return core.Task{}, decodeAPIError(resp)
	}
	var out core.Task
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return core.Task{}, err
	}
	return out, nil
}

func (c *Client) ListTasks() ([]core.Task, error) {
	resp, err := c.doJSON(http.MethodGet, "/v1/tasks", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, decodeAPIError(resp)
	}
	var out struct {
		Items []core.Task `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out.Items, nil
}

func (c *Client) GetTask(id string) (core.Task, error) {
	resp, err := c.doJSON(http.MethodGet, "/v1/tasks/"+id, nil)
	if err != nil {
		return core.Task{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return core.Task{}, decodeAPIError(resp)
	}
	var out core.Task
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return core.Task{}, err
	}
	return out, nil
}

func (c *Client) UpdateTask(id string, title string) (core.Task, error) {
	payload := map[string]string{"title": title}
	resp, err := c.doJSON(http.MethodPatch, "/v1/tasks/"+id, payload)
	if err != nil {
		return core.Task{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return core.Task{}, decodeAPIError(resp)
	}
	var out core.Task
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return core.Task{}, err
	}
	return out, nil
}

func (c *Client) CompleteTask(id string) (core.Task, error) {
	resp, err := c.doJSON(http.MethodPost, "/v1/tasks/"+id+"/complete", struct{}{})
	if err != nil {
		return core.Task{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return core.Task{}, decodeAPIError(resp)
	}
	var out core.Task
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return core.Task{}, err
	}
	return out, nil
}

func (c *Client) PatchTaskStatus(id, status string) (core.Task, error) {
	payload := map[string]string{"status": status}
	resp, err := c.doJSON(http.MethodPatch, "/v1/tasks/"+id, payload)
	if err != nil {
		return core.Task{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return core.Task{}, decodeAPIError(resp)
	}
	var out core.Task
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return core.Task{}, err
	}
	return out, nil
}

func (c *Client) DeleteTask(id string) error {
	resp, err := c.doJSON(http.MethodDelete, "/v1/tasks/"+id, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return decodeAPIError(resp)
	}
	return nil
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

	resp, err := c.http.Do(req)
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
			if strings.HasPrefix(line, "data: ") {
				dataLine = strings.TrimPrefix(line, "data: ")
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
