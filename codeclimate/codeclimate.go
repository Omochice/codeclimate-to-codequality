package codeclimate

import (
	"bytes"
	"encoding/json"
	"io"
)

type Issue struct {
	Type        string   `json:"type"`
	CheckName   string   `json:"check_name"`
	Description string   `json:"description"`
	Categories  []string `json:"categories"`
	Severity    string   `json:"severity"`
	Fingerprint string   `json:"fingerprint"`
	Location    Location `json:"location"`
}

type Location struct {
	Path      string     `json:"path"`
	Lines     *Lines     `json:"lines,omitempty"`
	Positions *Positions `json:"positions,omitempty"`
}

type Lines struct {
	Begin int `json:"begin"`
	End   int `json:"end,omitempty"`
}

type Positions struct {
	Begin Position  `json:"begin"`
	End   *Position `json:"end,omitempty"`
}

type Position struct {
	Line   int `json:"line"`
	Column int `json:"column,omitempty"`
}

// Parse decodes null-byte delimited CodeClimate JSON from r.
func Parse(r io.Reader) ([]Issue, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	segments := bytes.Split(data, []byte{0})
	issues := make([]Issue, 0, len(segments))

	for _, seg := range segments {
		seg = bytes.TrimSpace(seg)
		if len(seg) == 0 {
			continue
		}

		var issue Issue
		if err := json.Unmarshal(seg, &issue); err != nil {
			return nil, err
		}
		issues = append(issues, issue)
	}

	return issues, nil
}
