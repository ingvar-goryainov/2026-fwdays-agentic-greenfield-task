package rules_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ingvar-goryainov/agents-lint/internal/lint"
	"github.com/ingvar-goryainov/agents-lint/internal/lint/rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestC002_ID(t *testing.T) {
	assert.Equal(t, "C002", rules.C002{}.ID())
}

func TestC002_Check(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		wantFinding bool
	}{
		{"tool with evidence, and a non-catalog bare word ignored", "../../../testdata/C002/valid.md", false},
		{"tool referenced with no evidence", "../../../testdata/C002/invalid.md", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := lint.Parse(tt.path)
			require.NoError(t, err)

			findings := rules.C002{}.Check(doc, repoRoot)
			if !tt.wantFinding {
				assert.Empty(t, findings)
				return
			}
			require.Len(t, findings, 1)
			assert.Equal(t, rules.RuleC002, findings[0].RuleID)
			assert.Equal(t, lint.SeverityWarning, findings[0].Severity)
			assert.Equal(t, doc.Path, findings[0].File)
			assert.NotZero(t, findings[0].Line)
		})
	}
}

func TestC002_Check_NoCodeSpans(t *testing.T) {
	doc := &lint.Document{Path: "AGENTS.md"}
	assert.Empty(t, rules.C002{}.Check(doc, repoRoot))
}

func TestHasEvidence(t *testing.T) {
	t.Run("glob pattern matches", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "main.tf"), []byte(""), 0o644))
		doc := &lint.Document{Path: "AGENTS.md", CodeSpans: []lint.CodeSpan{{Text: "terraform", Line: 1}}}
		assert.Empty(t, rules.C002{}.Check(doc, dir))
	})

	t.Run("exact filename matches", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte(""), 0o644))
		doc := &lint.Document{Path: "AGENTS.md", CodeSpans: []lint.CodeSpan{{Text: "go", Line: 1}}}
		assert.Empty(t, rules.C002{}.Check(doc, dir))
	})

	t.Run("directory name matches", func(t *testing.T) {
		dir := t.TempDir()
		require.NoError(t, os.Mkdir(filepath.Join(dir, "k8s"), 0o755))
		doc := &lint.Document{Path: "AGENTS.md", CodeSpans: []lint.CodeSpan{{Text: "kubectl", Line: 1}}}
		assert.Empty(t, rules.C002{}.Check(doc, dir))
	})

	t.Run("no patterns match", func(t *testing.T) {
		dir := t.TempDir()
		doc := &lint.Document{Path: "AGENTS.md", CodeSpans: []lint.CodeSpan{{Text: "yarn", Line: 3}}}
		findings := rules.C002{}.Check(doc, dir)
		require.Len(t, findings, 1)
		assert.Equal(t, 3, findings[0].Line)
	})
}
