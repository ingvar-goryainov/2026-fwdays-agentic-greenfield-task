package rules_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ingvar-goryainov/agents-lint/internal/lint"
	"github.com/ingvar-goryainov/agents-lint/internal/lint/rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultRules_FixedOrder(t *testing.T) {
	got := rules.DefaultRules()
	require.Len(t, got, 4)
	want := []string{rules.RuleS002, rules.RuleS003, rules.RuleS004, rules.RuleS005}
	for i, id := range want {
		assert.Equal(t, id, got[i].ID())
	}
}

func TestKnownRuleIDs(t *testing.T) {
	want := []string{
		rules.RuleS001, rules.RuleS002, rules.RuleS003, rules.RuleS004, rules.RuleS005,
	}
	assert.Equal(t, want, rules.KnownRuleIDs())
}

func TestRun_MissingFileShortCircuits(t *testing.T) {
	path := filepath.Join(t.TempDir(), "does-not-exist.md")

	findings, err := rules.Run(path)
	require.NoError(t, err)
	require.Len(t, findings, 1)
	assert.Equal(t, rules.RuleS001, findings[0].RuleID)
	assert.Equal(t, lint.SeverityError, findings[0].Severity)
}

func TestRun_ExistingFileRunsAllRules(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "AGENTS.md")
	// Missing Agent(s) section (S002) and has no agent blocks (so S003/S004
	// find nothing to flag); no frontmatter (S005 passes).
	require.NoError(t, writeTestFile(path, "## Overview\n\nNo agents here.\n"))

	findings, err := rules.Run(path)
	require.NoError(t, err)
	require.Len(t, findings, 1)
	assert.Equal(t, rules.RuleS002, findings[0].RuleID)
}

func TestRun_CleanFileHasNoFindings(t *testing.T) {
	findings, err := rules.Run("../../../testdata/S003/valid.md")
	require.NoError(t, err)
	assert.Empty(t, findings)
}

func TestRun_UnreadableFilePropagatesError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("chmod-based unreadable file simulation is unix-specific")
	}
	if os.Geteuid() == 0 {
		t.Skip("permissions are not enforced when running as root")
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "AGENTS.md")
	require.NoError(t, writeTestFile(path, "## Agents\n"))
	require.NoError(t, os.Chmod(path, 0o000))
	t.Cleanup(func() { _ = os.Chmod(path, 0o644) })

	findings, err := rules.Run(path)
	assert.Error(t, err)
	assert.Nil(t, findings)
}

func TestRun_AggregatesInFixedOrder(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "AGENTS.md")
	content := "---\nname: \"unterminated\nversion: 1\n---\n\n## Overview\n\nNo agents here.\n"
	require.NoError(t, writeTestFile(path, content))

	findings, err := rules.Run(path)
	require.NoError(t, err)
	require.Len(t, findings, 2)
	assert.Equal(t, rules.RuleS002, findings[0].RuleID)
	assert.Equal(t, rules.RuleS005, findings[1].RuleID)
}
