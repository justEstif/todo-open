package core

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

const defaultLeaseTTLSeconds = 300

// AgentInfo is the ephemeral claim state stored under ext["agent"].
type AgentInfo struct {
	ID              string    `json:"id"`
	ClaimedAt       time.Time `json:"claimed_at"`
	LeaseExpiresAt  time.Time `json:"lease_expires_at"`
	HeartbeatAt     time.Time `json:"heartbeat_at"`
	LeaseTTLSeconds int       `json:"lease_ttl_seconds"`
}

// getAgentExt extracts the agent claim info from a task's ext map.
// Returns nil if not present or on decode error.
func getAgentExt(t Task) *AgentInfo {
	raw, ok := t.Ext["agent"]
	if !ok {
		return nil
	}
	var info AgentInfo
	if err := json.Unmarshal(raw, &info); err != nil {
		return nil
	}
	return &info
}

// setAgentExt encodes info and stores it under ext["agent"].
func setAgentExt(t *Task, info AgentInfo) {
	b, err := json.Marshal(info)
	if err != nil {
		return
	}
	if t.Ext == nil {
		t.Ext = make(map[string]json.RawMessage)
	}
	t.Ext["agent"] = b
}

// clearAgentExt removes ext["agent"] from the task, nil-ing the map if empty.
func clearAgentExt(t *Task) {
	if len(t.Ext) == 0 {
		return
	}
	delete(t.Ext, "agent")
	if len(t.Ext) == 0 {
		t.Ext = nil
	}
}

// priorityRank maps priority to a sort key (higher = more urgent).
var priorityRank = map[TaskPriority]int{
	TaskPriorityCritical: 4,
	TaskPriorityHigh:     3,
	TaskPriorityNormal:   2,
	TaskPriorityLow:      1,
	"deferred":           0,
	"":                   2, // treat unset as normal
}

func isLeaseExpired(info *AgentInfo, now time.Time) bool {
	return info == nil || now.After(info.LeaseExpiresAt)
}

// NextTask returns the highest-priority unclaimed open task.
func (s *Service) NextTask(ctx context.Context) (Task, error) {
	all, err := s.repo.List(ctx)
	if err != nil {
		return Task{}, err
	}
	now := s.nowFn().UTC()
	var best *Task
	bestRank := -1
	for i := range all {
		t := &all[i]
		if t.Status != TaskStatusOpen {
			continue
		}
		agent := getAgentExt(*t)
		if agent != nil && agent.ID != "" && !isLeaseExpired(agent, now) {
			// Claimed and lease not expired.
			continue
		}
		rank := priorityRank[t.Priority]
		if best == nil || rank > bestRank {
			best = t
			bestRank = rank
		}
	}
	if best == nil {
		return Task{}, fmt.Errorf("no unclaimed open tasks: %w", ErrNotFound)
	}
	return *best, nil
}

// ClaimTask atomically claims a task for an agent.
func (s *Service) ClaimTask(ctx context.Context, id, agentID string, leaseTTLSeconds int) (Task, error) {
	if id == "" {
		return Task{}, fmt.Errorf("id is required: %w", ErrInvalidInput)
	}
	if agentID == "" {
		return Task{}, fmt.Errorf("agent_id is required: %w", ErrInvalidInput)
	}
	if leaseTTLSeconds <= 0 {
		leaseTTLSeconds = defaultLeaseTTLSeconds
	}

	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return Task{}, err
	}
	now := s.nowFn().UTC()

	// Check preconditions.
	if task.Status != TaskStatusOpen {
		return Task{}, fmt.Errorf("task status must be open to claim (got %s): %w", task.Status, ErrConflict)
	}
	existing := getAgentExt(task)
	if existing != nil && existing.ID != "" && !isLeaseExpired(existing, now) {
		return Task{}, fmt.Errorf("task already claimed by agent %s: %w", existing.ID, ErrConflict)
	}

	// Apply claim.
	oldStatus := task.Status
	lease := now.Add(time.Duration(leaseTTLSeconds) * time.Second)
	setAgentExt(&task, AgentInfo{
		ID:              agentID,
		ClaimedAt:       now,
		LeaseExpiresAt:  lease,
		HeartbeatAt:     now,
		LeaseTTLSeconds: leaseTTLSeconds,
	})
	task.Status = TaskStatusInProgress
	if task.StartedAt == nil {
		task.StartedAt = &now
	}
	task.UpdatedAt = now
	task.Version++
	result, err := s.repo.Update(ctx, task)
	if err == nil && task.Status != oldStatus {
		s.emitMutationEvent("task.status_changed", &result, &oldStatus, &task.Status)
	}
	return result, err
}

// HeartbeatTask extends the lease for an agent-owned task.
func (s *Service) HeartbeatTask(ctx context.Context, id, agentID string) (Task, error) {
	if id == "" {
		return Task{}, fmt.Errorf("id is required: %w", ErrInvalidInput)
	}
	if agentID == "" {
		return Task{}, fmt.Errorf("agent_id is required: %w", ErrInvalidInput)
	}

	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return Task{}, err
	}
	agent := getAgentExt(task)
	if agent == nil || agent.ID == "" {
		return Task{}, fmt.Errorf("task has no active claim: %w", ErrForbidden)
	}
	if agent.ID != agentID {
		return Task{}, fmt.Errorf("agent_id mismatch: %w", ErrForbidden)
	}

	now := s.nowFn().UTC()
	ttl := agent.LeaseTTLSeconds
	if ttl <= 0 {
		ttl = defaultLeaseTTLSeconds
	}
	agent.LeaseExpiresAt = now.Add(time.Duration(ttl) * time.Second)
	agent.HeartbeatAt = now
	setAgentExt(&task, *agent)
	task.UpdatedAt = now
	task.Version++
	return s.repo.Update(ctx, task)
}

// ReleaseTask releases an agent's claim on a task, transitioning it back to open.
func (s *Service) ReleaseTask(ctx context.Context, id, agentID string) (Task, error) {
	if id == "" {
		return Task{}, fmt.Errorf("id is required: %w", ErrInvalidInput)
	}
	if agentID == "" {
		return Task{}, fmt.Errorf("agent_id is required: %w", ErrInvalidInput)
	}

	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return Task{}, err
	}
	agent := getAgentExt(task)
	if agent == nil || agent.ID == "" {
		return Task{}, fmt.Errorf("task has no active claim: %w", ErrForbidden)
	}
	if agent.ID != agentID {
		return Task{}, fmt.Errorf("agent_id mismatch: %w", ErrForbidden)
	}

	now := s.nowFn().UTC()
	oldStatus := task.Status
	clearAgentExt(&task)
	task.Status = TaskStatusOpen
	task.UpdatedAt = now
	task.Version++
	result, err := s.repo.Update(ctx, task)
	if err == nil && task.Status != oldStatus {
		s.emitMutationEvent("task.status_changed", &result, &oldStatus, &task.Status)
	}
	return result, err
}

// SweepExpiredLeases finds in_progress tasks with expired leases and transitions them back to open.
func (s *Service) SweepExpiredLeases(ctx context.Context) (int, error) {
	all, err := s.repo.List(ctx)
	if err != nil {
		return 0, err
	}
	now := s.nowFn().UTC()
	count := 0
	for _, task := range all {
		if task.Status != TaskStatusInProgress {
			continue
		}
		agent := getAgentExt(task)
		if agent == nil || agent.ID == "" {
			continue
		}
		if !now.After(agent.LeaseExpiresAt) {
			continue
		}
		clearAgentExt(&task)
		task.Status = TaskStatusOpen
		task.UpdatedAt = now
		task.Version++
		if _, err := s.repo.Update(ctx, task); err != nil {
			return count, fmt.Errorf("sweep: update task %s: %w", task.ID, err)
		}
		count++
	}
	return count, nil
}
