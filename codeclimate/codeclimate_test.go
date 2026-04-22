package codeclimate_test

import (
	"strings"
	"testing"

	"github.com/Omochice/codeclimate-to-codequality/codeclimate"
)

func TestParse(t *testing.T) {
	t.Run("parses a single issue without null terminator", func(t *testing.T) {
		input := `{"type":"issue","check_name":"test-check","description":"test description","severity":"blocker","fingerprint":"fp1","location":{"path":"a.rb","lines":{"begin":1}}}`
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

	t.Run("parses multiple issues separated by null byte", func(t *testing.T) {
		input := "{\"type\":\"issue\",\"check_name\":\"check-1\",\"description\":\"desc1\",\"severity\":\"major\",\"fingerprint\":\"fp1\",\"location\":{\"path\":\"a.rb\",\"lines\":{\"begin\":1}}}\x00{\"type\":\"issue\",\"check_name\":\"check-2\",\"description\":\"desc2\",\"severity\":\"minor\",\"fingerprint\":\"fp2\",\"location\":{\"path\":\"b.rb\",\"lines\":{\"begin\":2}}}"
		reader := strings.NewReader(input)

		issues, err := codeclimate.Parse(reader)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(issues) != 2 {
			t.Fatalf("expected length %d, got %d", 2, len(issues))
		}
		if issues[0].CheckName != "check-1" {
			t.Fatalf("got %v, want %v", issues[0].CheckName, "check-1")
		}
		if issues[1].CheckName != "check-2" {
			t.Fatalf("got %v, want %v", issues[1].CheckName, "check-2")
		}
	})

	t.Run("parses multiple issues terminated by null byte", func(t *testing.T) {
		input := "{\"type\":\"issue\",\"check_name\":\"check-1\",\"description\":\"desc1\",\"severity\":\"major\",\"fingerprint\":\"fp1\",\"location\":{\"path\":\"a.rb\",\"lines\":{\"begin\":1}}}\x00{\"type\":\"issue\",\"check_name\":\"check-2\",\"description\":\"desc2\",\"severity\":\"minor\",\"fingerprint\":\"fp2\",\"location\":{\"path\":\"b.rb\",\"lines\":{\"begin\":2}}}\x00"
		reader := strings.NewReader(input)

		issues, err := codeclimate.Parse(reader)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(issues) != 2 {
			t.Fatalf("expected length %d, got %d", 2, len(issues))
		}
	})

	t.Run("returns empty slice for empty input", func(t *testing.T) {
		reader := strings.NewReader("")

		issues, err := codeclimate.Parse(reader)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(issues) != 0 {
			t.Fatalf("expected length %d, got %d", 0, len(issues))
		}
	})

	t.Run("trims whitespace around each JSON object", func(t *testing.T) {
		input := "\n  {\"type\":\"issue\",\"check_name\":\"check-1\",\"description\":\"desc\",\"severity\":\"info\",\"fingerprint\":\"fp1\",\"location\":{\"path\":\"a.rb\",\"lines\":{\"begin\":1}}}  \n\x00\n  {\"type\":\"issue\",\"check_name\":\"check-2\",\"description\":\"desc\",\"severity\":\"info\",\"fingerprint\":\"fp2\",\"location\":{\"path\":\"b.rb\",\"lines\":{\"begin\":2}}}  \n"
		reader := strings.NewReader(input)

		issues, err := codeclimate.Parse(reader)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(issues) != 2 {
			t.Fatalf("expected length %d, got %d", 2, len(issues))
		}
	})

	t.Run("skips whitespace-only segments between null delimiters", func(t *testing.T) {
		input := "{\"type\":\"issue\",\"check_name\":\"check-1\",\"description\":\"desc\",\"severity\":\"info\",\"fingerprint\":\"fp1\",\"location\":{\"path\":\"a.rb\",\"lines\":{\"begin\":1}}}\x00\n\x00{\"type\":\"issue\",\"check_name\":\"check-2\",\"description\":\"desc\",\"severity\":\"info\",\"fingerprint\":\"fp2\",\"location\":{\"path\":\"b.rb\",\"lines\":{\"begin\":2}}}"
		reader := strings.NewReader(input)

		issues, err := codeclimate.Parse(reader)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(issues) != 2 {
			t.Fatalf("expected length %d, got %d", 2, len(issues))
		}
	})

	t.Run("returns error for invalid JSON in a segment", func(t *testing.T) {
		input := "{\"type\":\"issue\"}\x00{invalid json"
		reader := strings.NewReader(input)

		issues, err := codeclimate.Parse(reader)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if issues != nil {
			t.Fatalf("expected nil, got %v", issues)
		}
	})

	t.Run("parses issue with positions instead of lines", func(t *testing.T) {
		input := `{"type":"issue","check_name":"pos-check","description":"desc","severity":"minor","fingerprint":"fp2","location":{"path":"b.rb","positions":{"begin":{"line":5,"column":10}}}}`
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
		input := `{"type":"issue","check_name":"cat-check","description":"desc","categories":["Bug Risk","Style"],"severity":"info","fingerprint":"fp3","location":{"path":"c.rb","lines":{"begin":1}}}`
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
