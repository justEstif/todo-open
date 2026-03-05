package core

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrNotFound     = errors.New("task not found")
)

type TaskService interface {
	CreateTask(ctx context.Context, title string) (Task, error)
	GetTask(ctx context.Context, id string) (Task, error)
	ListTasks(ctx context.Context) ([]Task, error)
	UpdateTask(ctx context.Context, id string, title string) (Task, error)
	DeleteTask(ctx context.Context, id string) error
}

type IDGenerator func() string

type TaskRepository interface {
	Create(ctx context.Context, task Task) (Task, error)
	GetByID(ctx context.Context, id string) (Task, error)
	List(ctx context.Context) ([]Task, error)
	Update(ctx context.Context, task Task) (Task, error)
	Delete(ctx context.Context, id string, deletedAt time.Time) error
}

type Service struct {
	repo  TaskRepository
	nowFn func() time.Time
	idFn  IDGenerator
}

func NewService(repo TaskRepository, nowFn func() time.Time, idFn IDGenerator) *Service {
	if nowFn == nil {
		nowFn = time.Now
	}
	if idFn == nil {
		idFn = func() string { return fmt.Sprintf("task_%d", nowFn().UnixNano()) }
	}
	return &Service{repo: repo, nowFn: nowFn, idFn: idFn}
}

func (s *Service) CreateTask(ctx context.Context, title string) (Task, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return Task{}, fmt.Errorf("title is required: %w", ErrInvalidInput)
	}
	now := s.nowFn().UTC()
	task := Task{ID: s.idFn(), Title: title, Status: TaskStatusOpen, CreatedAt: now, UpdatedAt: now}
	return s.repo.Create(ctx, task)
}

func (s *Service) GetTask(ctx context.Context, id string) (Task, error) {
	if strings.TrimSpace(id) == "" {
		return Task{}, fmt.Errorf("id is required: %w", ErrInvalidInput)
	}
	return s.repo.GetByID(ctx, id)
}

func (s *Service) ListTasks(ctx context.Context) ([]Task, error) {
	return s.repo.List(ctx)
}

func (s *Service) UpdateTask(ctx context.Context, id string, title string) (Task, error) {
	if strings.TrimSpace(id) == "" {
		return Task{}, fmt.Errorf("id is required: %w", ErrInvalidInput)
	}
	title = strings.TrimSpace(title)
	if title == "" {
		return Task{}, fmt.Errorf("title is required: %w", ErrInvalidInput)
	}

	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return Task{}, err
	}
	task.Title = title
	task.UpdatedAt = s.nowFn().UTC()
	return s.repo.Update(ctx, task)
}

func (s *Service) DeleteTask(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("id is required: %w", ErrInvalidInput)
	}
	return s.repo.Delete(ctx, id, s.nowFn().UTC())
}
