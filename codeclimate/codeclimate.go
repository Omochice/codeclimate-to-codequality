package codeclimate

import (
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

// Parse decodes a CodeClimate JSON array from r.
func Parse(r io.Reader) ([]Issue, error) {
	var issues []Issue

	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&issues); err != nil {
		return nil, err
	}

	if issues == nil {
		issues = []Issue{}
	}

	return issues, nil
}
