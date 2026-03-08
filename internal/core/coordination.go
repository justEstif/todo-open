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
// Returns nil if not present.
func getAgentExt(t Task) *AgentInfo {
	if t.Ext == nil {
		return nil
	}
	extMap, ok := toMap(t.Ext)
	if !ok {
		return nil
	}
	agentRaw, ok := extMap["agent"]
	if !ok || agentRaw == nil {
		return nil
	}
	// Round-trip through JSON to handle map[string]any.
	b, err := json.Marshal(agentRaw)
	if err != nil {
		return nil
	}
	var info AgentInfo
	if err := json.Unmarshal(b, &info); err != nil {
		return nil
	}
	return &info
}

// setAgentExt sets ext["agent"] on a task's ext map.
func setAgentExt(t *Task, info AgentInfo) {
	extMap := getOrInitExtMap(t)
	extMap["agent"] = info
	t.Ext = extMap
}

// clearAgentExt removes ext["agent"] from a task's ext map.
func clearAgentExt(t *Task) {
	if t.Ext == nil {
		return
	}
	extMap, ok := toMap(t.Ext)
	if !ok {
		return
	}
	delete(extMap, "agent")
	if len(extMap) == 0 {
		t.Ext = nil
	} else {
		t.Ext = extMap
	}
}

func getOrInitExtMap(t *Task) map[string]any {
	if t.Ext == nil {
		return map[string]any{}
	}
	m, ok := toMap(t.Ext)
	if !ok {
		return map[string]any{}
	}
	return m
}

func toMap(v any) (map[string]any, bool) {
	if m, ok := v.(map[string]any); ok {
		return m, true
	}
	// Handle case where ext was set as a typed struct — re-encode.
	b, err := json.Marshal(v)
	if err != nil {
		return nil, false
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, false
	}
	return m, true
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
	return s.repo.Update(ctx, task)
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
	clearAgentExt(&task)
	task.Status = TaskStatusOpen
	task.UpdatedAt = now
	task.Version++
	return s.repo.Update(ctx, task)
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
