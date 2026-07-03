package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/ingvar-goryainov/agents-lint/internal/lint/rules"
	"github.com/ingvar-goryainov/agents-lint/internal/lint/skeleton"
	"github.com/ingvar-goryainov/agents-lint/internal/reporter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunInit(t *testing.T) {
	t.Run("default path, no existing file", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "AGENTS.md")

		var buf bytes.Buffer
		exitCode, err := runInit(&buf, path, false)
		require.NoError(t, err)
		assert.Equal(t, 0, exitCode)
		assert.Contains(t, buf.String(), path)

		got, readErr := os.ReadFile(path)
		require.NoError(t, readErr)
		assert.Equal(t, skeleton.Content(), got)
	})

	t.Run("explicit path argument", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "docs", "AGENTS.md")
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))

		var buf bytes.Buffer
		exitCode, err := runInit(&buf, path, false)
		require.NoError(t, err)
		assert.Equal(t, 0, exitCode)

		_, statErr := os.Stat(path)
		assert.NoError(t, statErr)
	})

	t.Run("existing file without --force is refused", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "AGENTS.md")
		require.NoError(t, os.WriteFile(path, []byte("hand-authored"), 0o644))

		var buf bytes.Buffer
		exitCode, err := runInit(&buf, path, false)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "--force")
		assert.Equal(t, 1, exitCode)
		assert.Empty(t, buf.String())

		got, readErr := os.ReadFile(path)
		require.NoError(t, readErr)
		assert.Equal(t, "hand-authored", string(got))
	})

	t.Run("existing file with --force is overwritten", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "AGENTS.md")
		require.NoError(t, os.WriteFile(path, []byte("hand-authored"), 0o644))

		var buf bytes.Buffer
		exitCode, err := runInit(&buf, path, true)
		require.NoError(t, err)
		assert.Equal(t, 0, exitCode)
		assert.Contains(t, buf.String(), path)

		got, readErr := os.ReadFile(path)
		require.NoError(t, readErr)
		assert.Equal(t, skeleton.Content(), got)
	})
}

// TestRunInit_ThenScan asserts the FR-CLI-03 guarantee end-to-end: a file
// written by runInit passes runScan with zero errors and the FR-OUT-04
// success line.
func TestRunInit_ThenScan(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "AGENTS.md")

	var initBuf bytes.Buffer
	exitCode, err := runInit(&initBuf, path, false)
	require.NoError(t, err)
	require.Equal(t, 0, exitCode)

	findings, err := rules.Run(path)
	require.NoError(t, err)
	require.Empty(t, findings)

	var scanBuf bytes.Buffer
	require.NoError(t, reporter.WriteText(&scanBuf, findings, len(rules.DefaultRules())+1, reporter.Options{}))
	assert.Contains(t, scanBuf.String(), "AGENTS.md is valid")

	scanExitCode, scanErr := runScan(&bytes.Buffer{}, path, "", formatText)
	require.NoError(t, scanErr)
	assert.Equal(t, 0, scanExitCode)
}
