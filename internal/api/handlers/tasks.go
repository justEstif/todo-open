package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ebeyene/todo-open/internal/core"
)

type TaskService interface {
	CreateTask(ctx context.Context, title string) (core.Task, error)
	GetTask(ctx context.Context, id string) (core.Task, error)
	ListTasks(ctx context.Context) ([]core.Task, error)
	UpdateTask(ctx context.Context, id string, title string) (core.Task, error)
	DeleteTask(ctx context.Context, id string) error
}

type TaskHandler struct {
	svc TaskService
}

func NewTaskHandler(svc TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

type taskPayload struct {
	Title string `json:"title"`
}

type errorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var payload taskPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "invalid JSON payload")
		return
	}
	task, err := h.svc.CreateTask(r.Context(), payload.Title)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.svc.ListTasks(r.Context())
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": tasks})
}

func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	task, err := h.svc.GetTask(r.Context(), id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var payload taskPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "invalid JSON payload")
		return
	}
	task, err := h.svc.UpdateTask(r.Context(), id, payload.Title)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.DeleteTask(r.Context(), id); err != nil {
		writeServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, core.ErrInvalidInput):
		writeError(w, http.StatusBadRequest, "validation_error", err.Error())
	case errors.Is(err, core.ErrNotFound):
		writeError(w, http.StatusNotFound, "not_found", "task not found")
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

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
