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
}
