package converter

import (
	"strings"

	"github.com/Omochice/codeclimate-to-codequality/codeclimate"
	"github.com/Omochice/codeclimate-to-codequality/codequality"
)

// Severity maps a CodeClimate severity to a GitLab Code Quality severity.
// GitLab does not support "blocker", so it is mapped to "critical".
func Severity(severity string) string {
	switch strings.ToLower(severity) {
	case "blocker":
		return "critical"
	case "critical", "major", "minor", "info":
		return strings.ToLower(severity)
	default:
		return "info"
	}
}

func line(issue codeclimate.Issue) int {
	if issue.Location.Lines != nil {
		return issue.Location.Lines.Begin
	}
	if issue.Location.Positions != nil {
		return issue.Location.Positions.Begin.Line
	}
	return 0
}

// Issues converts CodeClimate issues into CodeQuality violations.
// Issues that lack a path, line, check_name, description, or fingerprint are skipped.
func Issues(issues []codeclimate.Issue) []codequality.Violation {
	violations := make([]codequality.Violation, 0, len(issues))

	for _, issue := range issues {
		l := line(issue)
		if issue.Location.Path == "" || l == 0 || issue.CheckName == "" || issue.Description == "" || issue.Fingerprint == "" {
			continue
		}

		path := strings.TrimPrefix(issue.Location.Path, "./")

		violation := codequality.Violation{
			Description: issue.Description,
			CheckName:   issue.CheckName,
			Fingerprint: issue.Fingerprint,
			Severity:    Severity(issue.Severity),
			Location: codequality.Location{
				Path: path,
				Lines: codequality.Lines{
					Begin: l,
				},
			},
		}

		violations = append(violations, violation)
	}

	return violations
}
