package cmd

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ingvar-goryainov/agents-lint/internal/lint"
	"github.com/ingvar-goryainov/agents-lint/internal/lint/rules"
	"github.com/ingvar-goryainov/agents-lint/internal/reporter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// failingWriter returns an error on every Write, used to exercise
// runScan's reporter.WriteText error-propagation path.
type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) {
	return 0, errors.New("write failed")
}

func TestRunScan(t *testing.T) {
	t.Run("clean file exits 0", func(t *testing.T) {
		var buf bytes.Buffer
		exitCode, err := runScan(&buf, "../../../testdata/S003/valid.md")
		require.NoError(t, err)
		assert.Equal(t, 0, exitCode)
		assert.Contains(t, buf.String(), "AGENTS.md is valid")
	})

	t.Run("error-severity finding exits 1", func(t *testing.T) {
		var buf bytes.Buffer
		exitCode, err := runScan(&buf, "../../../testdata/S002/invalid.md")
		require.NoError(t, err)
		assert.Equal(t, 1, exitCode)
		assert.Contains(t, buf.String(), "S002")
	})

	t.Run("unreadable file reports an error and exits 1", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("chmod-based unreadable file simulation is unix-specific")
		}
		if os.Geteuid() == 0 {
			t.Skip("permissions are not enforced when running as root")
		}

		dir := t.TempDir()
		path := filepath.Join(dir, "AGENTS.md")
		require.NoError(t, os.WriteFile(path, []byte("## Agents\n"), 0o644))
		require.NoError(t, os.Chmod(path, 0o000))
		t.Cleanup(func() { _ = os.Chmod(path, 0o644) })

		var buf bytes.Buffer
		exitCode, err := runScan(&buf, path)
		assert.Error(t, err)
		assert.Equal(t, 1, exitCode)
		assert.Empty(t, buf.String(), "no reporter output should be written on an unexpected error")
	})

	t.Run("reporter write failure is propagated", func(t *testing.T) {
		exitCode, err := runScan(failingWriter{}, "../../../testdata/S003/valid.md")
		assert.Error(t, err)
		assert.Equal(t, 1, exitCode)
	})
}

func TestRunScan_OutputMatchesReporter(t *testing.T) {
	path := "../../../testdata/S002/invalid.md"

	findings, err := rules.Run(path)
	require.NoError(t, err)

	var want bytes.Buffer
	require.NoError(t, reporter.WriteText(&want, findings, len(rules.DefaultRules())+1, reporter.Options{}))

	var got bytes.Buffer
	_, err = runScan(&got, path)
	require.NoError(t, err)

	assert.Equal(t, want.String(), got.String())
}

func TestExitCodeForFindings(t *testing.T) {
	tests := []struct {
		name     string
		findings []lint.Finding
		want     int
	}{
		{"no findings", nil, 0},
		{
			"only warnings",
			[]lint.Finding{
				{RuleID: "C002", Severity: lint.SeverityWarning},
				{RuleID: "C002", Severity: lint.SeverityWarning},
			},
			0,
		},
		{
			"at least one error",
			[]lint.Finding{
				{RuleID: "C002", Severity: lint.SeverityWarning},
				{RuleID: "S002", Severity: lint.SeverityError},
			},
			1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, exitCodeForFindings(tt.findings))
		})
	}
}
