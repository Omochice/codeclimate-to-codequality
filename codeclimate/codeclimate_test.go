package codeclimate_test

import (
	"strings"
	"testing"

	"github.com/Omochice/codeclimate-to-codequality/codeclimate"
)

func TestParse(t *testing.T) {
	t.Run("parses valid JSON array", func(t *testing.T) {
		input := `[{"type":"issue","check_name":"test-check","description":"test description","severity":"blocker","fingerprint":"fp1","location":{"path":"a.rb","lines":{"begin":1}}}]`
		reader := strings.NewReader(input)

		issues, err := codeclimate.Parse(reader)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(issues) != 1 {
			t.Fatalf("expected length %d, got %d", 1, len(issues))
		}
		if issues[0].CheckName != "test-check" {
			t.Fatalf("got %v, want %v", issues[0].CheckName, "test-check")
		}
		if issues[0].Fingerprint != "fp1" {
			t.Fatalf("got %v, want %v", issues[0].Fingerprint, "fp1")
		}
	})

	t.Run("returns error for invalid JSON", func(t *testing.T) {
		input := `{invalid json`
		reader := strings.NewReader(input)

		issues, err := codeclimate.Parse(reader)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if issues != nil {
			t.Fatalf("expected nil, got %v", issues)
		}
	})

	t.Run("handles empty array", func(t *testing.T) {
		input := `[]`
		reader := strings.NewReader(input)

		issues, err := codeclimate.Parse(reader)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(issues) != 0 {
			t.Fatalf("expected length %d, got %d", 0, len(issues))
		}
	})

	t.Run("parses issue with positions instead of lines", func(t *testing.T) {
		input := `[{"type":"issue","check_name":"pos-check","description":"desc","severity":"minor","fingerprint":"fp2","location":{"path":"b.rb","positions":{"begin":{"line":5,"column":10}}}}]`
		reader := strings.NewReader(input)

		issues, err := codeclimate.Parse(reader)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(issues) != 1 {
			t.Fatalf("expected length %d, got %d", 1, len(issues))
		}
		if issues[0].Location.Positions == nil {
			t.Fatal("expected non-nil Positions")
		}
		if issues[0].Location.Positions.Begin.Line != 5 {
			t.Fatalf("got %v, want %v", issues[0].Location.Positions.Begin.Line, 5)
		}
		if issues[0].Location.Positions.Begin.Column != 10 {
			t.Fatalf("got %v, want %v", issues[0].Location.Positions.Begin.Column, 10)
		}
	})

	t.Run("parses issue with categories", func(t *testing.T) {
		input := `[{"type":"issue","check_name":"cat-check","description":"desc","categories":["Bug Risk","Style"],"severity":"info","fingerprint":"fp3","location":{"path":"c.rb","lines":{"begin":1}}}]`
		reader := strings.NewReader(input)

		issues, err := codeclimate.Parse(reader)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(issues) != 1 {
			t.Fatalf("expected length %d, got %d", 1, len(issues))
		}
		if len(issues[0].Categories) != 2 {
			t.Fatalf("expected %d categories, got %d", 2, len(issues[0].Categories))
		}
	})
}
