package app

import "testing"

func TestNewViewRegistry(t *testing.T) {
	r, err := NewViewRegistry()
	if err != nil {
		t.Fatalf("new view registry: %v", err)
	}

	names := r.Names()
	if len(names) != 1 || names[0] != "json" {
		t.Fatalf("names = %v, want [json]", names)
	}
}

func TestNewSyncRegistry(t *testing.T) {
	r, err := NewSyncRegistry()
	if err != nil {
		t.Fatalf("new sync registry: %v", err)
	}

	names := r.Names()
	if len(names) != 1 || names[0] != "noop" {
		t.Fatalf("names = %v, want [noop]", names)
	}
}
