package lint_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ingvar-goryainov/agents-lint/internal/lint"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "AGENTS.md")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	return path
}

func TestParse_NoFrontmatter(t *testing.T) {
	path := writeFile(t, "## Agents\n\n### architect\n\nDescription.\n")
	doc, err := lint.Parse(path)
	require.NoError(t, err)
	assert.Nil(t, doc.Frontmatter)
}

func TestParse_ValidFrontmatter(t *testing.T) {
	content := "---\nkey: value\n---\n\n## Agents\n\n### architect\n\nDescription.\n"
	path := writeFile(t, content)
	doc, err := lint.Parse(path)
	require.NoError(t, err)
	require.NotNil(t, doc.Frontmatter)
	assert.Equal(t, "key: value\n", doc.Frontmatter.Raw)
	assert.Equal(t, 1, doc.Frontmatter.Line)
	// Line numbers after the frontmatter block must still align with the
	// original file, not the stripped body.
	require.Len(t, doc.H2Sections, 1)
	assert.Equal(t, 5, doc.H2Sections[0].Line)
}

func TestParse_UnterminatedFrontmatterDelimiter(t *testing.T) {
	content := "---\nkey: value\n\n## Agents\n"
	path := writeFile(t, content)
	doc, err := lint.Parse(path)
	require.NoError(t, err)
	require.NotNil(t, doc.Frontmatter)
	assert.Equal(t, "key: value\n\n## Agents\n", doc.Frontmatter.Raw)
	assert.Empty(t, doc.H2Sections)
}

func TestParse_NoAgentSection(t *testing.T) {
	path := writeFile(t, "## Overview\n\nSome text.\n")
	doc, err := lint.Parse(path)
	require.NoError(t, err)
	require.Len(t, doc.H2Sections, 1)
	assert.Equal(t, "Overview", doc.H2Sections[0].Text)
	assert.Empty(t, doc.AgentBlocks)
}

func TestParse_MultipleAgentBlocks(t *testing.T) {
	content := `## Agents

### architect

Senior architect.

**Role:** Designs things.

### implementer

Writes code.

**Tools:** go build
`
	path := writeFile(t, content)
	doc, err := lint.Parse(path)
	require.NoError(t, err)
	require.Len(t, doc.AgentBlocks, 2)
	assert.Equal(t, "architect", doc.AgentBlocks[0].Name)
	assert.True(t, doc.AgentBlocks[0].HasRole)
	assert.Equal(t, "implementer", doc.AgentBlocks[1].Name)
	assert.True(t, doc.AgentBlocks[1].HasTools)
}

func TestParse_BoldLabelVariants(t *testing.T) {
	tests := []struct {
		name            string
		label           string
		wantInstruction bool
	}{
		{"Instructions label", "**Instructions:** do things", true},
		{"Responsibilities label", "**Responsibilities:** do things", true},
		{"Unrelated label", "**Note:** do things", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := "## Agents\n\n### architect\n\nDescription.\n\n" + tt.label + "\n"
			path := writeFile(t, content)
			doc, err := lint.Parse(path)
			require.NoError(t, err)
			require.Len(t, doc.AgentBlocks, 1)
			assert.Equal(t, tt.wantInstruction, doc.AgentBlocks[0].HasInstructions)
		})
	}
}

func TestParse_AgentSectionTitleCaseInsensitive(t *testing.T) {
	path := writeFile(t, "## AGENT\n\n### architect\n\nDescription.\n")
	doc, err := lint.Parse(path)
	require.NoError(t, err)
	require.Len(t, doc.AgentBlocks, 1)
}

func TestParse_H3OutsideAgentSectionIgnored(t *testing.T) {
	content := "## Repository structure\n\n### not-an-agent\n\nSome text.\n"
	path := writeFile(t, content)
	doc, err := lint.Parse(path)
	require.NoError(t, err)
	assert.Empty(t, doc.AgentBlocks)
}

func TestParse_MissingFile(t *testing.T) {
	_, err := lint.Parse(filepath.Join(t.TempDir(), "does-not-exist.md"))
	assert.Error(t, err)
}

func TestParse_H1HeadingEndsAgentSection(t *testing.T) {
	content := "## Agents\n\n### architect\n\nDescription.\n\n# Unrelated title\n\n### not-an-agent\n\nText.\n"
	path := writeFile(t, content)
	doc, err := lint.Parse(path)
	require.NoError(t, err)
	require.Len(t, doc.AgentBlocks, 1)
	assert.Equal(t, "architect", doc.AgentBlocks[0].Name)
}

func TestParse_ContextLabel(t *testing.T) {
	content := "## Agents\n\n### architect\n\nDescription.\n\n**Context:** `docs/requirements.md`\n"
	path := writeFile(t, content)
	doc, err := lint.Parse(path)
	require.NoError(t, err)
	require.Len(t, doc.AgentBlocks, 1)
	assert.True(t, doc.AgentBlocks[0].HasContext)
}

func TestParse_EmptyHeadingHasLineZero(t *testing.T) {
	content := "## Agents\n\n### \n\nSome text.\n"
	path := writeFile(t, content)
	doc, err := lint.Parse(path)
	require.NoError(t, err)
	require.Len(t, doc.AgentBlocks, 1)
	assert.Empty(t, doc.AgentBlocks[0].Name)
	assert.Zero(t, doc.AgentBlocks[0].Line)
}

func TestParse_CodeSpansAcrossNesting(t *testing.T) {
	content := `## Agents

Paragraph mentions ` + "`internal/lint/rule.go`" + `.

- List item references ` + "`go.mod`" + `

> Blockquote references ` + "`Dockerfile`" + `
`
	path := writeFile(t, content)
	doc, err := lint.Parse(path)
	require.NoError(t, err)
	require.Len(t, doc.CodeSpans, 3)

	texts := []string{doc.CodeSpans[0].Text, doc.CodeSpans[1].Text, doc.CodeSpans[2].Text}
	assert.ElementsMatch(t, []string{"internal/lint/rule.go", "go.mod", "Dockerfile"}, texts)
	for _, cs := range doc.CodeSpans {
		assert.NotZero(t, cs.Line)
	}
}

func TestParse_CodeSpanInFrontmatterNotCollected(t *testing.T) {
	content := "---\nkey: `not-a-real-span`\n---\n\n## Agents\n\nText with `real-span` here.\n"
	path := writeFile(t, content)
	doc, err := lint.Parse(path)
	require.NoError(t, err)
	require.Len(t, doc.CodeSpans, 1)
	assert.Equal(t, "real-span", doc.CodeSpans[0].Text)
}

func TestParse_FrontmatterClosingDelimiterAtEOF(t *testing.T) {
	content := "---\nkey: value\n---"
	path := writeFile(t, content)
	doc, err := lint.Parse(path)
	require.NoError(t, err)
	require.NotNil(t, doc.Frontmatter)
	assert.Equal(t, "key: value\n", doc.Frontmatter.Raw)
}
