package core

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"strings"
	"time"
)

type ValidationMode string

const (
	ValidationModeStrict ValidationMode = "strict"
	ValidationModeCompat ValidationMode = "compat"
)

type ValidationIssue struct {
	Line    int
	Field   string
	Message string
	Context string
}

var coreFields = map[string]struct{}{
	"id": {}, "title": {}, "status": {}, "created_at": {}, "updated_at": {},
	"description": {}, "project": {}, "tags": {}, "priority": {}, "due_at": {},
	"started_at": {}, "completed_at": {}, "deleted_at": {}, "parent_id": {},
	"assignee": {}, "estimate_minutes": {}, "sort_order": {}, "version": {}, "ext": {},
	"trigger_ids": {}, "blocking": {}, "blocked_by": {},
}

func ValidateTaskJSONL(r io.Reader, mode ValidationMode) ([]ValidationIssue, error) {
	scanner := bufio.NewScanner(r)
	issues := []ValidationIssue{}
	line := 0
	for scanner.Scan() {
		line++
		raw := strings.TrimSpace(scanner.Text())
		if raw == "" {
			continue
		}

		var record map[string]any
		if err := json.Unmarshal([]byte(raw), &record); err != nil {
			issues = append(issues, ValidationIssue{Line: line, Field: "<json>", Message: fmt.Sprintf("invalid JSON: %v", err), Context: raw})
			continue
		}
		issues = append(issues, validateTaskRecord(record, line, raw, mode)...)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return issues, nil
}

func validateTaskRecord(record map[string]any, line int, raw string, mode ValidationMode) []ValidationIssue {
	issues := []ValidationIssue{}

	reqString(record, "id", line, raw, &issues)
	title := reqString(record, "title", line, raw, &issues)
	status := reqString(record, "status", line, raw, &issues)
	createdAt := reqRFC3339UTC(record, "created_at", line, raw, &issues)
	reqRFC3339UTC(record, "updated_at", line, raw, &issues)

	if strings.TrimSpace(title) == "" {
		issues = append(issues, ValidationIssue{Line: line, Field: "title", Message: "must be non-empty", Context: raw})
	}

	if status != "" {
		switch TaskStatus(status) {
		case TaskStatusPending, TaskStatusOpen, TaskStatusInProgress, TaskStatusDone, TaskStatusArchived:
		default:
			issues = append(issues, ValidationIssue{Line: line, Field: "status", Message: "must be one of pending|open|in_progress|done|archived", Context: raw})
		}
	}

	if p, ok := record["priority"]; ok {
		if s, ok := p.(string); !ok {
			issues = append(issues, ValidationIssue{Line: line, Field: "priority", Message: "must be a string", Context: raw})
		} else {
			switch TaskPriority(s) {
			case TaskPriorityLow, TaskPriorityNormal, TaskPriorityHigh, TaskPriorityCritical:
			default:
				issues = append(issues, ValidationIssue{Line: line, Field: "priority", Message: "must be one of low|normal|high|critical", Context: raw})
			}
		}
	}

	if tags, ok := record["tags"]; ok {
		arr, ok := tags.([]any)
		if !ok {
			issues = append(issues, ValidationIssue{Line: line, Field: "tags", Message: "must be an array of strings", Context: raw})
		} else {
			seen := map[string]struct{}{}
			for _, v := range arr {
				s, ok := v.(string)
				if !ok {
					issues = append(issues, ValidationIssue{Line: line, Field: "tags", Message: "must contain only strings", Context: raw})
					break
				}
				if _, exists := seen[s]; exists {
					issues = append(issues, ValidationIssue{Line: line, Field: "tags", Message: "must contain unique values", Context: raw})
					break
				}
				seen[s] = struct{}{}
			}
		}
	}

	if v, ok := record["version"]; ok {
		n, ok := asInt(v)
		if !ok || n < 1 {
			issues = append(issues, ValidationIssue{Line: line, Field: "version", Message: "must be an integer >= 1", Context: raw})
		}
	}

	if v, ok := record["estimate_minutes"]; ok {
		n, ok := asInt(v)
		if !ok || n < 0 {
			issues = append(issues, ValidationIssue{Line: line, Field: "estimate_minutes", Message: "must be an integer >= 0", Context: raw})
		}
	}

	if v, ok := record["ext"]; ok {
		if _, ok := v.(map[string]any); !ok {
			issues = append(issues, ValidationIssue{Line: line, Field: "ext", Message: "must be an object", Context: raw})
		}
	}

	if status == string(TaskStatusDone) {
		completedAt := reqRFC3339UTC(record, "completed_at", line, raw, &issues)
		if !createdAt.IsZero() && !completedAt.IsZero() && completedAt.Before(createdAt) {
			issues = append(issues, ValidationIssue{Line: line, Field: "completed_at", Message: "must be >= created_at", Context: raw})
		}
	}
	if startedAt, ok := optRFC3339UTC(record, "started_at", line, raw, &issues); ok {
		if !createdAt.IsZero() && startedAt.Before(createdAt) {
			issues = append(issues, ValidationIssue{Line: line, Field: "started_at", Message: "must be >= created_at", Context: raw})
		}
	}
	if completedAt, ok := optRFC3339UTC(record, "completed_at", line, raw, &issues); ok {
		if !createdAt.IsZero() && completedAt.Before(createdAt) {
			issues = append(issues, ValidationIssue{Line: line, Field: "completed_at", Message: "must be >= created_at", Context: raw})
		}
	}

	if mode == ValidationModeStrict {
		for k := range record {
			if _, ok := coreFields[k]; !ok {
				issues = append(issues, ValidationIssue{Line: line, Field: k, Message: "unknown top-level field; place custom fields under ext", Context: raw})
			}
		}
	}

	return issues
}

func reqString(record map[string]any, field string, line int, raw string, issues *[]ValidationIssue) string {
	v, ok := record[field]
	if !ok {
		*issues = append(*issues, ValidationIssue{Line: line, Field: field, Message: "is required", Context: raw})
		return ""
	}
	s, ok := v.(string)
	if !ok {
		*issues = append(*issues, ValidationIssue{Line: line, Field: field, Message: "must be a string", Context: raw})
		return ""
	}
	return s
}

func reqRFC3339UTC(record map[string]any, field string, line int, raw string, issues *[]ValidationIssue) time.Time {
	v, ok := record[field]
	if !ok {
		*issues = append(*issues, ValidationIssue{Line: line, Field: field, Message: "is required", Context: raw})
		return time.Time{}
	}
	t, ok := parseRFC3339UTC(v)
	if !ok {
		*issues = append(*issues, ValidationIssue{Line: line, Field: field, Message: "must be RFC3339 UTC timestamp (e.g. 2026-03-05T18:00:00Z)", Context: raw})
		return time.Time{}
	}
	return t
}

func optRFC3339UTC(record map[string]any, field string, line int, raw string, issues *[]ValidationIssue) (time.Time, bool) {
	v, ok := record[field]
	if !ok {
		return time.Time{}, false
	}
	t, valid := parseRFC3339UTC(v)
	if !valid {
		*issues = append(*issues, ValidationIssue{Line: line, Field: field, Message: "must be RFC3339 UTC timestamp (e.g. 2026-03-05T18:00:00Z)", Context: raw})
		return time.Time{}, false
	}
	return t, true
}

func parseRFC3339UTC(v any) (time.Time, bool) {
	s, ok := v.(string)
	if !ok {
		return time.Time{}, false
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, false
	}
	return t, t.Location() == time.UTC
}

func asInt(v any) (int, bool) {
	f, ok := v.(float64)
	if !ok {
		return 0, false
	}
	if math.Mod(f, 1) != 0 {
		return 0, false
	}
	return int(f), true
}
