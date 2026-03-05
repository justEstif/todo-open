package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/justEstif/todo-open/internal/core"
)

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
