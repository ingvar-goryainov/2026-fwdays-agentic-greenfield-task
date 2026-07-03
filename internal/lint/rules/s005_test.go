package rules_test

import (
	"path/filepath"
	"testing"

	"github.com/ingvar-goryainov/agents-lint/internal/lint"
	"github.com/ingvar-goryainov/agents-lint/internal/lint/rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestS005_Check(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		wantFinding bool
	}{
		{"valid frontmatter", "../../../testdata/S005/valid.md", false},
		{"malformed frontmatter", "../../../testdata/S005/invalid.md", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := lint.Parse(tt.path)
			require.NoError(t, err)

			findings := rules.S005{}.Check(doc)
			if !tt.wantFinding {
				assert.Empty(t, findings)
				return
			}
			require.Len(t, findings, 1)
			assert.Equal(t, rules.RuleS005, findings[0].RuleID)
			assert.Equal(t, lint.SeverityError, findings[0].Severity)
			assert.NotEmpty(t, findings[0].Message)
		})
	}
}

func TestS005_NoFrontmatter(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "AGENTS.md")
	require.NoError(t, writeTestFile(path, "## Agents\n\n### architect\n\nDescription.\n"))

	doc, err := lint.Parse(path)
	require.NoError(t, err)

	assert.Empty(t, rules.S005{}.Check(doc))
}
