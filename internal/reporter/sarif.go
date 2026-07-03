package reporter

import (
	"encoding/json"
	"io"

	"github.com/ingvar-goryainov/agents-lint/internal/lint"
)

const sarifSchemaURI = "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json"

// sarifLog is the top-level SARIF v2.1.0 document.
type sarifLog struct {
	Schema  string     `json:"$schema"`
	Version string     `json:"version"`
	Runs    []sarifRun `json:"runs"`
}

type sarifRun struct {
	Tool    sarifTool     `json:"tool"`
	Results []sarifResult `json:"results"`
}

type sarifTool struct {
	Driver sarifDriver `json:"driver"`
}

type sarifDriver struct {
	Name  string                     `json:"name"`
	Rules []sarifReportingDescriptor `json:"rules"`
}

type sarifReportingDescriptor struct {
	ID string `json:"id"`
}

type sarifResult struct {
	RuleID    string          `json:"ruleId"`
	Level     string          `json:"level"`
	Message   sarifMessage    `json:"message"`
	Locations []sarifLocation `json:"locations"`
}

type sarifMessage struct {
	Text string `json:"text"`
}

type sarifLocation struct {
	PhysicalLocation sarifPhysicalLocation `json:"physicalLocation"`
}

type sarifPhysicalLocation struct {
	ArtifactLocation sarifArtifactLocation `json:"artifactLocation"`
	Region           *sarifRegion          `json:"region,omitempty"`
}

type sarifArtifactLocation struct {
	URI string `json:"uri"`
}

type sarifRegion struct {
	StartLine int `json:"startLine"`
}

// WriteSARIF renders findings as a SARIF v2.1.0 log (FR-OUT-03) to w. It
// takes no ruleCount/Options, unlike WriteText: SARIF has no "N rules
// passed" success line and no color option — an empty findings slice
// simply produces a log with empty "results" and "rules" arrays.
func WriteSARIF(w io.Writer, findings []lint.Finding) error {
	log := sarifLog{
		Schema:  sarifSchemaURI,
		Version: "2.1.0",
		Runs: []sarifRun{
			{
				Tool: sarifTool{
					Driver: sarifDriver{
						Name:  "agents-lint",
						Rules: sarifRules(findings),
					},
				},
				Results: sarifResults(findings),
			},
		},
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(log)
}

// sarifRules returns one sarifReportingDescriptor per distinct RuleID
// present in findings, in first-seen order.
func sarifRules(findings []lint.Finding) []sarifReportingDescriptor {
	rules := make([]sarifReportingDescriptor, 0, len(findings))
	seen := make(map[string]bool, len(findings))
	for _, f := range findings {
		if seen[f.RuleID] {
			continue
		}
		seen[f.RuleID] = true
		rules = append(rules, sarifReportingDescriptor{ID: f.RuleID})
	}
	return rules
}

func sarifResults(findings []lint.Finding) []sarifResult {
	results := make([]sarifResult, 0, len(findings))
	for _, f := range findings {
		results = append(results, sarifResult{
			RuleID:  f.RuleID,
			Level:   sarifLevel(f.Severity),
			Message: sarifMessage{Text: f.Message},
			Locations: []sarifLocation{
				{PhysicalLocation: newPhysicalLocation(f)},
			},
		})
	}
	return results
}

// newPhysicalLocation builds a physicalLocation from a Finding's File
// and Line. Line == 0 (today, only FR-S001's file-not-found case) omits
// Region entirely rather than emitting an invalid region.startLine: 0.
func newPhysicalLocation(f lint.Finding) sarifPhysicalLocation {
	loc := sarifPhysicalLocation{
		ArtifactLocation: sarifArtifactLocation{URI: f.File},
	}
	if f.Line > 0 {
		loc.Region = &sarifRegion{StartLine: f.Line}
	}
	return loc
}

// sarifLevel maps lint.Severity to a SARIF result level.
func sarifLevel(s lint.Severity) string {
	switch s {
	case lint.SeverityError:
		return "error"
	case lint.SeverityWarning:
		return "warning"
	default:
		return string(s)
	}
}
