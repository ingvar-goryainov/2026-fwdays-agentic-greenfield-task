package rules_test

import (
	"testing"

	"github.com/ingvar-goryainov/agents-lint/internal/lint"
	"github.com/ingvar-goryainov/agents-lint/internal/lint/rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestS002_Check(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		wantFinding bool
	}{
		{"has Agents section", "../../../testdata/S002/valid.md", false},
		{"missing Agent(s) section", "../../../testdata/S002/invalid.md", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := lint.Parse(tt.path)
			require.NoError(t, err)

			findings := rules.S002{}.Check(doc)
			if !tt.wantFinding {
				assert.Empty(t, findings)
				return
			}
			require.Len(t, findings, 1)
			assert.Equal(t, rules.RuleS002, findings[0].RuleID)
			assert.Equal(t, lint.SeverityError, findings[0].Severity)
		})
	}
}

func TestS002_CaseInsensitiveSpelling(t *testing.T) {
	for _, title := range []string{"## Agent", "## Agents", "## AGENT", "## agents"} {
		t.Run(title, func(t *testing.T) {
			dir := t.TempDir()
			path := dir + "/AGENTS.md"
			require.NoError(t, writeTestFile(path, title+"\n\n### architect\n\nDescription.\n"))

			doc, err := lint.Parse(path)
			require.NoError(t, err)
			assert.Empty(t, rules.S002{}.Check(doc))
		})
	}
}
