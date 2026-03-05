package jsonl

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/justEstif/todo-open/internal/core"
)

const (
	defaultSchemaVersion = "todo.open.task.v1"
	metaDirName          = ".todoopen"
	metaFileName         = "meta.json"
	tasksFileName        = "tasks.jsonl"
)

type workspaceMeta struct {
	WorkspaceVersion int      `json:"workspace_version"`
	SchemaVersion    string   `json:"schema_version"`
	DefaultSort      []string `json:"default_sort,omitempty"`
	EnabledViews     []string `json:"enabled_views,omitempty"`
	EnabledAdapters  []string `json:"enabled_sync_adapters,omitempty"`
}

type TaskRepo struct {
	rootPath  string
	tasksPath string
	metaPath  string
}

func NewTaskRepo(rootPath string) *TaskRepo {
	return &TaskRepo{
		rootPath:  rootPath,
		tasksPath: filepath.Join(rootPath, tasksFileName),
		metaPath:  filepath.Join(rootPath, metaDirName, metaFileName),
	}
}

func (r *TaskRepo) Create(_ context.Context, task core.Task) (core.Task, error) {
	if err := r.ensureWorkspace(); err != nil {
		return core.Task{}, err
	}
	tasks, err := r.readAllTasks()
	if err != nil {
		return core.Task{}, err
	}
	for _, existing := range tasks {
		if existing.ID == task.ID {
			return core.Task{}, fmt.Errorf("task id already exists: %w", core.ErrInvalidInput)
		}
	}
	tasks = append(tasks, task)
	if err := r.writeAllTasks(tasks); err != nil {
		return core.Task{}, err
	}
	return task, nil
}

func (r *TaskRepo) GetByID(_ context.Context, id string) (core.Task, error) {
	if err := r.ensureWorkspace(); err != nil {
		return core.Task{}, err
	}
	tasks, err := r.readAllTasks()
	if err != nil {
		return core.Task{}, err
	}
	for _, task := range tasks {
		if task.ID == id {
			if task.DeletedAt != nil {
				return core.Task{}, core.ErrNotFound
			}
			return task, nil
		}
	}
	return core.Task{}, core.ErrNotFound
}

func (r *TaskRepo) List(_ context.Context) ([]core.Task, error) {
	if err := r.ensureWorkspace(); err != nil {
		return nil, err
	}
	tasks, err := r.readAllTasks()
	if err != nil {
		return nil, err
	}
	out := make([]core.Task, 0, len(tasks))
	for _, task := range tasks {
		if task.DeletedAt == nil {
			out = append(out, task)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(out[j].CreatedAt) })
	return out, nil
}

func (r *TaskRepo) Update(_ context.Context, task core.Task) (core.Task, error) {
	if err := r.ensureWorkspace(); err != nil {
		return core.Task{}, err
	}
	tasks, err := r.readAllTasks()
	if err != nil {
		return core.Task{}, err
	}
	updated := false
	for i := range tasks {
		if tasks[i].ID == task.ID {
			if tasks[i].DeletedAt != nil {
				return core.Task{}, core.ErrNotFound
			}
			tasks[i] = task
			updated = true
			break
		}
	}
	if !updated {
		return core.Task{}, core.ErrNotFound
	}
	if err := r.writeAllTasks(tasks); err != nil {
		return core.Task{}, err
	}
	return task, nil
}

func (r *TaskRepo) Delete(_ context.Context, id string, deletedAt time.Time) error {
	if err := r.ensureWorkspace(); err != nil {
		return err
	}
	tasks, err := r.readAllTasks()
	if err != nil {
		return err
	}
	updated := false
	for i := range tasks {
		if tasks[i].ID == id {
			if tasks[i].DeletedAt != nil {
				return core.ErrNotFound
			}
			tasks[i].DeletedAt = &deletedAt
			tasks[i].Status = core.TaskStatusArchived
			tasks[i].UpdatedAt = deletedAt
			tasks[i].Version++
			updated = true
			break
		}
	}
	if !updated {
		return core.ErrNotFound
	}
	return r.writeAllTasks(tasks)
}

func (r *TaskRepo) ensureWorkspace() error {
	if err := os.MkdirAll(r.rootPath, 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(r.metaPath), 0o755); err != nil {
		return err
	}
	if _, err := os.Stat(r.metaPath); errors.Is(err, os.ErrNotExist) {
		meta := workspaceMeta{WorkspaceVersion: 1, SchemaVersion: defaultSchemaVersion}
		return writeJSONAtomic(r.metaPath, meta)
	} else if err != nil {
		return err
	}

	var meta workspaceMeta
	data, err := os.ReadFile(r.metaPath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &meta); err != nil {
		return fmt.Errorf("invalid metadata file: %w", err)
	}
	if meta.WorkspaceVersion < 1 {
		return fmt.Errorf("unsupported workspace_version: %d", meta.WorkspaceVersion)
	}
	if strings.TrimSpace(meta.SchemaVersion) != defaultSchemaVersion {
		return fmt.Errorf("unsupported schema_version: %s", meta.SchemaVersion)
	}
	return nil
}

func (r *TaskRepo) readAllTasks() ([]core.Task, error) {
	f, err := os.Open(r.tasksPath)
	if errors.Is(err, os.ErrNotExist) {
		return []core.Task{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var tasks []core.Task
	scanner := bufio.NewScanner(f)
	line := 0
	for scanner.Scan() {
		line++
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}
		var task core.Task
		if err := json.Unmarshal([]byte(text), &task); err != nil {
			return nil, fmt.Errorf("invalid JSONL at line %d: %w", line, err)
		}
		tasks = append(tasks, task)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *TaskRepo) writeAllTasks(tasks []core.Task) error {
	var b strings.Builder
	enc := json.NewEncoder(&b)
	for _, task := range tasks {
		if err := enc.Encode(task); err != nil {
			return err
		}
	}
	return writeFileAtomic(r.tasksPath, []byte(b.String()))
}

func writeJSONAtomic(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return writeFileAtomic(path, data)
}

func writeFileAtomic(path string, data []byte) error {
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
