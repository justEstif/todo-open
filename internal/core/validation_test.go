package core

import (
	"strings"
	"testing"
)

func TestValidateTaskJSONLValidRecord(t *testing.T) {
	t.Parallel()

	input := `{"id":"task_1","title":"Ship MVP","status":"open","created_at":"2026-03-05T18:00:00Z","updated_at":"2026-03-05T18:00:00Z","priority":"normal","tags":["mvp"],"version":1,"ext":{"kanban":{"column":"todo"}}}`
	issues, err := ValidateTaskJSONL(strings.NewReader(input+"\n"), ValidationModeStrict)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %+v", issues)
	}
}

func TestValidateTaskJSONLReportsLineAndContext(t *testing.T) {
	t.Parallel()

	input := strings.Join([]string{
		`{"id":"task_1","title":"ok","status":"open","created_at":"2026-03-05T18:00:00Z","updated_at":"2026-03-05T18:00:00Z"}`,
		`{"id":"task_2","title":"","status":"doing","created_at":"bad","updated_at":"2026-03-05T18:00:00Z","extra":"x"}`,
	}, "\n") + "\n"

	issues, err := ValidateTaskJSONL(strings.NewReader(input), ValidationModeStrict)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) == 0 {
		t.Fatal("expected issues")
	}
	for _, issue := range issues {
		if issue.Line != 2 {
			t.Fatalf("expected issue on line 2, got line %d (%+v)", issue.Line, issue)
		}
		if issue.Context == "" {
			t.Fatalf("expected context in issue: %+v", issue)
		}
	}
}

func TestValidateTaskJSONLCompatAllowsUnknownTopLevel(t *testing.T) {
	t.Parallel()

	input := `{"id":"task_1","title":"Ship","status":"open","created_at":"2026-03-05T18:00:00Z","updated_at":"2026-03-05T18:00:00Z","extra":"kept"}`
	issues, err := ValidateTaskJSONL(strings.NewReader(input+"\n"), ValidationModeCompat)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 0 {
		t.Fatalf("expected no issues in compat mode, got %+v", issues)
	}
}
