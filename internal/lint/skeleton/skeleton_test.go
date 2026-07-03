package skeleton

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ingvar-goryainov/agents-lint/internal/lint"
	"github.com/ingvar-goryainov/agents-lint/internal/lint/rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestContent_PassesAllSchemaRules is the regression test guaranteeing
// FR-CLI-03: the generated skeleton must always pass the real S001-S005
// rule set, not a hand-maintained assumption of what those rules check.
func TestContent_PassesAllSchemaRules(t *testing.T) {
	path := filepath.Join(t.TempDir(), "AGENTS.md")
	require.NoError(t, os.WriteFile(path, Content(), 0o644))

	findings, err := rules.Run(path)
	require.NoError(t, err)
	assert.Empty(t, findings, "generated skeleton must produce zero findings against S001-S005")
}

func TestContent_NotEmpty(t *testing.T) {
	assert.NotEmpty(t, Content())
}

// TestContent_ValidDocument sanity-checks the skeleton parses into the
// structural shape S003 expects, independent of the S001-S005 run above.
func TestContent_ValidDocument(t *testing.T) {
	path := filepath.Join(t.TempDir(), "AGENTS.md")
	require.NoError(t, os.WriteFile(path, Content(), 0o644))

	doc, err := lint.Parse(path)
	require.NoError(t, err)

	require.NotNil(t, doc.Frontmatter)
	require.Len(t, doc.AgentBlocks, 1)
	block := doc.AgentBlocks[0]
	assert.NotEmpty(t, block.Name)
	assert.True(t, block.HasRole)
	assert.True(t, block.HasInstructions || block.HasTools || block.HasContext)
}
