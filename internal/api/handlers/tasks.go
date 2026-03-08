package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/justEstif/todo-open/internal/core"
)

type TaskHandler struct {
	svc core.TaskService
}

func NewTaskHandler(svc core.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

type taskPayload struct {
	Title      string   `json:"title"`
	Status     string   `json:"status,omitempty"`
	TriggerIDs []string `json:"trigger_ids,omitempty"`
}

type errorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var payload taskPayload
	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "invalid JSON payload")
		return
	}
	task, err := h.svc.CreateTask(r.Context(), payload.Title, payload.TriggerIDs...)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeTaskJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := core.ListFilter{}
	if s := q.Get("status"); s != "" {
		filter.Status = core.TaskStatus(s)
	}
	if q.Get("is_blocked") == "true" {
		filter.IsBlocked = true
	}
	tasks, err := h.svc.ListTasks(r.Context(), filter)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": tasks})
}

func (h *TaskHandler) Complete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	task, err := h.svc.CompleteTask(r.Context(), id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeTaskJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	task, err := h.svc.GetTask(r.Context(), id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeTaskJSON(w, http.StatusOK, task)
}

// Upsert handles PUT /v1/tasks/{id} — idempotent create-or-update.
func (h *TaskHandler) Upsert(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var payload taskPayload
	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "invalid JSON payload")
		return
	}

	var ifMatch *int
	if hdr := r.Header.Get("If-Match"); hdr != "" {
		v, err := parseETag(hdr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid_etag", "invalid If-Match header")
			return
		}
		ifMatch = &v
	}

	task, created, err := h.svc.UpsertTask(r.Context(), id, payload.Title, ifMatch)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	if created {
		writeTaskJSON(w, http.StatusCreated, task)
	} else {
		writeTaskJSON(w, http.StatusOK, task)
	}
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var payload taskPayload
	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "invalid JSON payload")
		return
	}

	// Status-only idempotent patch.
	if payload.Title == "" && payload.Status != "" {
		task, err := h.svc.PatchStatus(r.Context(), id, core.TaskStatus(payload.Status))
		if err != nil {
			writeServiceError(w, err)
			return
		}
		writeTaskJSON(w, http.StatusOK, task)
		return
	}

	// ETag / If-Match enforcement.
	if ifMatch := r.Header.Get("If-Match"); ifMatch != "" {
		expectedVersion, err := parseETag(ifMatch)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid_etag", "invalid If-Match header")
			return
		}
		current, err := h.svc.GetTask(r.Context(), id)
		if err != nil {
			writeServiceError(w, err)
			return
		}
		if current.Version != expectedVersion {
			writeError(w, http.StatusConflict, "version_conflict", "ETag mismatch; resource was modified")
			return
		}
	}

	task, err := h.svc.UpdateTask(r.Context(), id, payload.Title)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeTaskJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.DeleteTask(r.Context(), id); err != nil {
		writeServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// NextTask handles GET /v1/tasks/next — returns the highest-priority unclaimed open task.
func (h *TaskHandler) NextTask(w http.ResponseWriter, r *http.Request) {
	task, err := h.svc.NextTask(r.Context())
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeTaskJSON(w, http.StatusOK, task)
}

// Claim handles POST /v1/tasks/{id}/claim.
type claimPayload struct {
	AgentID         string `json:"agent_id"`
	LeaseTTLSeconds int    `json:"lease_ttl_seconds"`
}

func (h *TaskHandler) Claim(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var payload claimPayload
	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "invalid JSON payload")
		return
	}
	task, err := h.svc.ClaimTask(r.Context(), id, payload.AgentID, payload.LeaseTTLSeconds)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeTaskJSON(w, http.StatusOK, task)
}

// Heartbeat handles POST /v1/tasks/{id}/heartbeat.
type agentPayload struct {
	AgentID string `json:"agent_id"`
}

func (h *TaskHandler) Heartbeat(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var payload agentPayload
	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "invalid JSON payload")
		return
	}
	task, err := h.svc.HeartbeatTask(r.Context(), id, payload.AgentID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeTaskJSON(w, http.StatusOK, task)
}

// Release handles POST /v1/tasks/{id}/release.
func (h *TaskHandler) Release(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var payload agentPayload
	if err := decodeJSON(r, &payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "invalid JSON payload")
		return
	}
	task, err := h.svc.ReleaseTask(r.Context(), id, payload.AgentID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeTaskJSON(w, http.StatusOK, task)
}

// parseETag parses a quoted ETag string like `"42"` into an integer.
func parseETag(etag string) (int, error) {
	if len(etag) >= 2 && etag[0] == '"' && etag[len(etag)-1] == '"' {
		etag = etag[1 : len(etag)-1]
	}
	v, err := strconv.Atoi(etag)
	if err != nil {
		return 0, fmt.Errorf("invalid etag: %w", err)
	}
	return v, nil
}

func decodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return err
	}
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("trailing data")
	}
	return nil
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, core.ErrInvalidInput):
		writeError(w, http.StatusBadRequest, "validation_error", err.Error())
	case errors.Is(err, core.ErrNotFound):
		writeError(w, http.StatusNotFound, "not_found", "task not found")
	case errors.Is(err, core.ErrConflict):
		writeError(w, http.StatusConflict, "conflict", err.Error())
	case errors.Is(err, core.ErrForbidden):
		writeError(w, http.StatusForbidden, "forbidden", err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal_error", "internal server error")
	}
}

func writeError(w http.ResponseWriter, status int, code string, message string) {
	resp := errorResponse{}
	resp.Error.Code = code
	resp.Error.Message = message
	writeJSON(w, status, resp)
}

// writeTaskJSON writes a task response and sets the ETag header.
func writeTaskJSON(w http.ResponseWriter, status int, task core.Task) {
	w.Header().Set("ETag", fmt.Sprintf(`"%d"`, task.Version))
	writeJSON(w, status, task)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
