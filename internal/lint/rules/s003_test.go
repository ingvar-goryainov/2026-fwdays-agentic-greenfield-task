package rules_test

import (
	"path/filepath"
	"testing"

	"github.com/ingvar-goryainov/agents-lint/internal/lint"
	"github.com/ingvar-goryainov/agents-lint/internal/lint/rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestS003_Check(t *testing.T) {
	t.Run("complete block passes", func(t *testing.T) {
		doc, err := lint.Parse("../../../testdata/S003/valid.md")
		require.NoError(t, err)
		assert.Empty(t, rules.S003{}.Check(doc))
	})

	t.Run("block missing role and all optional fields", func(t *testing.T) {
		doc, err := lint.Parse("../../../testdata/S003/invalid.md")
		require.NoError(t, err)
		findings := rules.S003{}.Check(doc)
		require.Len(t, findings, 2)
		for _, f := range findings {
			assert.Equal(t, rules.RuleS003, f.RuleID)
			assert.Equal(t, lint.SeverityError, f.Severity)
		}
	})
}

func TestS003_MissingRoleOnly(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "AGENTS.md")
	content := "## Agents\n\n### architect\n\n**Tools:** go build\n"
	require.NoError(t, writeTestFile(path, content))

	doc, err := lint.Parse(path)
	require.NoError(t, err)

	findings := rules.S003{}.Check(doc)
	require.Len(t, findings, 1)
	assert.Contains(t, findings[0].Message, "role")
}

func TestS003_MissingName(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "AGENTS.md")
	content := "## Agents\n\n### \n\n**Role:** Does things.\n\n**Tools:** go build\n"
	require.NoError(t, writeTestFile(path, content))

	doc, err := lint.Parse(path)
	require.NoError(t, err)

	findings := rules.S003{}.Check(doc)
	require.Len(t, findings, 1)
	assert.Contains(t, findings[0].Message, "no name")
}

func TestS003_ResponsibilitiesSynonym(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "AGENTS.md")
	content := "## Agents\n\n### architect\n\nDescription.\n\n**Responsibilities:**\n- Do things\n"
	require.NoError(t, writeTestFile(path, content))

	doc, err := lint.Parse(path)
	require.NoError(t, err)

	assert.Empty(t, rules.S003{}.Check(doc))
}
