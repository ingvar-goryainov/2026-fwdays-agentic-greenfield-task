package reporter_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/ingvar-goryainov/agents-lint/internal/lint"
	"github.com/ingvar-goryainov/agents-lint/internal/reporter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// sarifDoc mirrors just the fields these tests assert on, decoded from
// WriteSARIF's JSON output.
type sarifDoc struct {
	Schema  string `json:"$schema"`
	Version string `json:"version"`
	Runs    []struct {
		Tool struct {
			Driver struct {
				Name  string `json:"name"`
				Rules []struct {
					ID string `json:"id"`
				} `json:"rules"`
			} `json:"driver"`
		} `json:"tool"`
		Results []struct {
			RuleID  string `json:"ruleId"`
			Level   string `json:"level"`
			Message struct {
				Text string `json:"text"`
			} `json:"message"`
			Locations []struct {
				PhysicalLocation struct {
					ArtifactLocation struct {
						URI string `json:"uri"`
					} `json:"artifactLocation"`
					Region *struct {
						StartLine int `json:"startLine"`
					} `json:"region"`
				} `json:"physicalLocation"`
			} `json:"locations"`
		} `json:"results"`
	} `json:"runs"`
}

func writeSARIF(t *testing.T, findings []lint.Finding) sarifDoc {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, reporter.WriteSARIF(&buf, findings))

	var doc sarifDoc
	require.NoError(t, json.Unmarshal(buf.Bytes(), &doc))
	require.Len(t, doc.Runs, 1, "SARIF log must contain exactly one run")
	return doc
}

func TestWriteSARIF_SingleFinding(t *testing.T) {
	doc := writeSARIF(t, []lint.Finding{
		{RuleID: "S002", Severity: lint.SeverityError, File: "AGENTS.md", Line: 3, Message: "missing ## Agents section"},
	})

	assert.Equal(t, "2.1.0", doc.Version)
	assert.Equal(t, "agents-lint", doc.Runs[0].Tool.Driver.Name)
	require.Len(t, doc.Runs[0].Results, 1)

	result := doc.Runs[0].Results[0]
	assert.Equal(t, "S002", result.RuleID)
	assert.Equal(t, "missing ## Agents section", result.Message.Text)
}

func TestWriteSARIF_FindingsOrderPreserved(t *testing.T) {
	doc := writeSARIF(t, []lint.Finding{
		{RuleID: "S002", Severity: lint.SeverityError, File: "AGENTS.md", Line: 1},
		{RuleID: "S004", Severity: lint.SeverityWarning, File: "AGENTS.md", Line: 2},
		{RuleID: "S005", Severity: lint.SeverityError, File: "AGENTS.md", Line: 3},
	})

	require.Len(t, doc.Runs[0].Results, 3)
	assert.Equal(t, []string{"S002", "S004", "S005"}, []string{
		doc.Runs[0].Results[0].RuleID,
		doc.Runs[0].Results[1].RuleID,
		doc.Runs[0].Results[2].RuleID,
	})
}

func TestWriteSARIF_EmptyFindings(t *testing.T) {
	doc := writeSARIF(t, nil)

	assert.Empty(t, doc.Runs[0].Results)
	assert.Empty(t, doc.Runs[0].Tool.Driver.Rules)
}

func TestWriteSARIF_SeverityLevelMapping(t *testing.T) {
	doc := writeSARIF(t, []lint.Finding{
		{RuleID: "S002", Severity: lint.SeverityError, File: "AGENTS.md", Line: 1},
		{RuleID: "C002", Severity: lint.SeverityWarning, File: "AGENTS.md", Line: 2},
	})

	require.Len(t, doc.Runs[0].Results, 2)
	assert.Equal(t, "error", doc.Runs[0].Results[0].Level)
	assert.Equal(t, "warning", doc.Runs[0].Results[1].Level)
}

func TestWriteSARIF_RegionOmittedWhenLineIsZero(t *testing.T) {
	doc := writeSARIF(t, []lint.Finding{
		{RuleID: "S001", Severity: lint.SeverityError, File: "AGENTS.md", Line: 0, Message: "AGENTS.md not found"},
	})

	require.Len(t, doc.Runs[0].Results, 1)
	loc := doc.Runs[0].Results[0].Locations[0].PhysicalLocation
	assert.Equal(t, "AGENTS.md", loc.ArtifactLocation.URI)
	assert.Nil(t, loc.Region, "region must be omitted when Line == 0")
}

func TestWriteSARIF_RegionIncludesStartLine(t *testing.T) {
	doc := writeSARIF(t, []lint.Finding{
		{RuleID: "S002", Severity: lint.SeverityError, File: "AGENTS.md", Line: 12},
	})

	require.Len(t, doc.Runs[0].Results, 1)
	loc := doc.Runs[0].Results[0].Locations[0].PhysicalLocation
	require.NotNil(t, loc.Region)
	assert.Equal(t, 12, loc.Region.StartLine)
}

func TestWriteSARIF_RuleCatalogDeduplicated(t *testing.T) {
	doc := writeSARIF(t, []lint.Finding{
		{RuleID: "S002", Severity: lint.SeverityError, File: "AGENTS.md", Line: 1},
		{RuleID: "S004", Severity: lint.SeverityWarning, File: "AGENTS.md", Line: 2},
		{RuleID: "S002", Severity: lint.SeverityError, File: "AGENTS.md", Line: 5},
	})

	rules := doc.Runs[0].Tool.Driver.Rules
	require.Len(t, rules, 2)
	assert.Equal(t, "S002", rules[0].ID)
	assert.Equal(t, "S004", rules[1].ID)
}
