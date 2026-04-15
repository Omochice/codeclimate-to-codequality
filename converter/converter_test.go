package converter_test

import (
	"testing"

	"github.com/Omochice/codeclimate-to-codequality/codeclimate"
	"github.com/Omochice/codeclimate-to-codequality/converter"
)

func TestSeverity(t *testing.T) {
	tests := []struct {
		name     string
		severity string
		want     string
	}{
		{
			name:     "blocker maps to critical",
			severity: "blocker",
			want:     "critical",
		},
		{
			name:     "Blocker (capitalized) maps to critical",
			severity: "Blocker",
			want:     "critical",
		},
		{
			name:     "critical passes through",
			severity: "critical",
			want:     "critical",
		},
		{
			name:     "Critical (capitalized) passes through",
			severity: "Critical",
			want:     "critical",
		},
		{
			name:     "major passes through",
			severity: "major",
			want:     "major",
		},
		{
			name:     "Major (capitalized) passes through",
			severity: "Major",
			want:     "major",
		},
		{
			name:     "minor passes through",
			severity: "minor",
			want:     "minor",
		},
		{
			name:     "Minor (capitalized) passes through",
			severity: "Minor",
			want:     "minor",
		},
		{
			name:     "info passes through",
			severity: "info",
			want:     "info",
		},
		{
			name:     "Info (capitalized) passes through",
			severity: "Info",
			want:     "info",
		},
		{
			name:     "unknown severity maps to info",
			severity: "Unknown",
			want:     "info",
		},
		{
			name:     "empty severity maps to info",
			severity: "",
			want:     "info",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := converter.Severity(tt.severity)
			if got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIssues(t *testing.T) {
	t.Run("converts valid issue with lines", func(t *testing.T) {
		issues := []codeclimate.Issue{
			{
				Type:        "issue",
				CheckName:   "test-check",
				Description: "test description",
				Severity:    "blocker",
				Fingerprint: "fp1",
				Location: codeclimate.Location{
					Path:  "app/models/user.rb",
					Lines: &codeclimate.Lines{Begin: 42},
				},
			},
		}

		violations := converter.Issues(issues)
		if len(violations) != 1 {
			t.Fatalf("expected length %d, got %d", 1, len(violations))
		}

		v := violations[0]
		if v.Description != "test description" {
			t.Fatalf("got %v, want %v", v.Description, "test description")
		}
		if v.CheckName != "test-check" {
			t.Fatalf("got %v, want %v", v.CheckName, "test-check")
		}
		if v.Severity != "critical" {
			t.Fatalf("got %v, want %v", v.Severity, "critical")
		}
		if v.Location.Path != "app/models/user.rb" {
			t.Fatalf("got %v, want %v", v.Location.Path, "app/models/user.rb")
		}
		if v.Location.Lines.Begin != 42 {
			t.Fatalf("got %v, want %v", v.Location.Lines.Begin, 42)
		}
		if v.Fingerprint != "fp1" {
			t.Fatalf("got %v, want %v", v.Fingerprint, "fp1")
		}
	})

	t.Run("converts valid issue with positions", func(t *testing.T) {
		issues := []codeclimate.Issue{
			{
				Type:        "issue",
				CheckName:   "pos-check",
				Description: "position issue",
				Severity:    "minor",
				Fingerprint: "fp2",
				Location: codeclimate.Location{
					Path: "b.rb",
					Positions: &codeclimate.Positions{
						Begin: codeclimate.Position{Line: 10, Column: 5},
					},
				},
			},
		}

		violations := converter.Issues(issues)
		if len(violations) != 1 {
			t.Fatalf("expected length %d, got %d", 1, len(violations))
		}
		if violations[0].Location.Lines.Begin != 10 {
			t.Fatalf("got %v, want %v", violations[0].Location.Lines.Begin, 10)
		}
	})

	t.Run("skips issue with missing path", func(t *testing.T) {
		issues := []codeclimate.Issue{
			{
				CheckName:   "check",
				Description: "desc",
				Severity:    "info",
				Fingerprint: "fp",
				Location: codeclimate.Location{
					Path:  "",
					Lines: &codeclimate.Lines{Begin: 1},
				},
			},
		}

		violations := converter.Issues(issues)
		if len(violations) != 0 {
			t.Fatalf("expected length %d, got %d", 0, len(violations))
		}
	})

	t.Run("skips issue with missing line", func(t *testing.T) {
		issues := []codeclimate.Issue{
			{
				CheckName:   "check",
				Description: "desc",
				Severity:    "info",
				Fingerprint: "fp",
				Location: codeclimate.Location{
					Path: "a.rb",
				},
			},
		}

		violations := converter.Issues(issues)
		if len(violations) != 0 {
			t.Fatalf("expected length %d, got %d", 0, len(violations))
		}
	})

	t.Run("skips issue with missing check_name", func(t *testing.T) {
		issues := []codeclimate.Issue{
			{
				CheckName:   "",
				Description: "desc",
				Severity:    "info",
				Fingerprint: "fp",
				Location: codeclimate.Location{
					Path:  "a.rb",
					Lines: &codeclimate.Lines{Begin: 1},
				},
			},
		}

		violations := converter.Issues(issues)
		if len(violations) != 0 {
			t.Fatalf("expected length %d, got %d", 0, len(violations))
		}
	})

	t.Run("skips issue with missing description", func(t *testing.T) {
		issues := []codeclimate.Issue{
			{
				CheckName:   "check",
				Description: "",
				Severity:    "info",
				Fingerprint: "fp",
				Location: codeclimate.Location{
					Path:  "a.rb",
					Lines: &codeclimate.Lines{Begin: 1},
				},
			},
		}

		violations := converter.Issues(issues)
		if len(violations) != 0 {
			t.Fatalf("expected length %d, got %d", 0, len(violations))
		}
	})

	t.Run("skips issue with missing fingerprint", func(t *testing.T) {
		issues := []codeclimate.Issue{
			{
				CheckName:   "check",
				Description: "desc",
				Severity:    "info",
				Fingerprint: "",
				Location: codeclimate.Location{
					Path:  "a.rb",
					Lines: &codeclimate.Lines{Begin: 1},
				},
			},
		}

		violations := converter.Issues(issues)
		if len(violations) != 0 {
			t.Fatalf("expected length %d, got %d", 0, len(violations))
		}
	})

	t.Run("removes ./ prefix from path", func(t *testing.T) {
		issues := []codeclimate.Issue{
			{
				CheckName:   "check",
				Description: "desc",
				Severity:    "info",
				Fingerprint: "fp",
				Location: codeclimate.Location{
					Path:  "./app/models/user.rb",
					Lines: &codeclimate.Lines{Begin: 1},
				},
			},
		}

		violations := converter.Issues(issues)
		if len(violations) != 1 {
			t.Fatalf("expected length %d, got %d", 1, len(violations))
		}
		if violations[0].Location.Path != "app/models/user.rb" {
			t.Fatalf("got %v, want %v", violations[0].Location.Path, "app/models/user.rb")
		}
	})

	t.Run("handles empty array", func(t *testing.T) {
		issues := []codeclimate.Issue{}
		violations := converter.Issues(issues)
		if len(violations) != 0 {
			t.Fatalf("expected length %d, got %d", 0, len(violations))
		}
	})

	t.Run("processes multiple issues", func(t *testing.T) {
		issues := []codeclimate.Issue{
			{
				CheckName:   "check-a",
				Description: "desc a",
				Severity:    "critical",
				Fingerprint: "fp1",
				Location: codeclimate.Location{
					Path:  "a.rb",
					Lines: &codeclimate.Lines{Begin: 1},
				},
			},
			{
				CheckName:   "check-b",
				Description: "desc b",
				Severity:    "major",
				Fingerprint: "fp2",
				Location: codeclimate.Location{
					Path:  "b.rb",
					Lines: &codeclimate.Lines{Begin: 10},
				},
			},
		}

		violations := converter.Issues(issues)
		if len(violations) != 2 {
			t.Fatalf("expected length %d, got %d", 2, len(violations))
		}
		if violations[0].CheckName != "check-a" {
			t.Fatalf("got %v, want %v", violations[0].CheckName, "check-a")
		}
		if violations[1].CheckName != "check-b" {
			t.Fatalf("got %v, want %v", violations[1].CheckName, "check-b")
		}
	})

	t.Run("skips invalid issues and processes valid ones", func(t *testing.T) {
		issues := []codeclimate.Issue{
			{
				CheckName:   "check-a",
				Description: "desc a",
				Severity:    "critical",
				Fingerprint: "fp1",
				Location: codeclimate.Location{
					Path:  "a.rb",
					Lines: &codeclimate.Lines{Begin: 1},
				},
			},
			{
				CheckName:   "check-b",
				Description: "",
				Severity:    "major",
				Fingerprint: "fp2",
				Location: codeclimate.Location{
					Path:  "b.rb",
					Lines: &codeclimate.Lines{Begin: 10},
				},
			},
		}

		violations := converter.Issues(issues)
		if len(violations) != 1 {
			t.Fatalf("expected length %d, got %d", 1, len(violations))
		}
		if violations[0].CheckName != "check-a" {
			t.Fatalf("got %v, want %v", violations[0].CheckName, "check-a")
		}
	})
}
