package reporter_test

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/ingvar-goryainov/agents-lint/internal/lint"
	"github.com/ingvar-goryainov/agents-lint/internal/reporter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// failingWriter returns an error on every Write, used to exercise
// WriteText's error-propagation paths.
type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) {
	return 0, errors.New("write failed")
}

func readGolden(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile("../../testdata/reporter/golden/" + name)
	require.NoError(t, err)
	return string(data)
}

func TestWriteText_Golden(t *testing.T) {
	tests := []struct {
		name      string
		findings  []lint.Finding
		ruleCount int
		golden    string
	}{
		{
			name: "single error",
			findings: []lint.Finding{
				{RuleID: "S002", Severity: lint.SeverityError, File: "AGENTS.md", Line: 3, Message: "missing ## Agents section"},
			},
			golden: "single_error.txt",
		},
		{
			name: "multiple findings preserve input order",
			findings: []lint.Finding{
				{RuleID: "S005", Severity: lint.SeverityError, File: "AGENTS.md", Line: 1, Message: "frontmatter is not valid YAML"},
				{RuleID: "S002", Severity: lint.SeverityError, File: "AGENTS.md", Line: 5, Message: "missing ## Agents section"},
				{RuleID: "S004", Severity: lint.SeverityWarning, File: "AGENTS.md", Line: 20, Message: `duplicate agent name "architect"`},
			},
			golden: "multi_mixed_order.txt",
		},
		{
			name: "mixed severities",
			findings: []lint.Finding{
				{RuleID: "S002", Severity: lint.SeverityError, File: "AGENTS.md", Line: 3, Message: "missing ## Agents section"},
				{RuleID: "S003", Severity: lint.SeverityError, File: "AGENTS.md", Line: 9, Message: `agent block "tester" has no role`},
				{RuleID: "C002", Severity: lint.SeverityWarning, File: "AGENTS.md", Line: 15, Message: `tool "kubectl" has no evidence in repo`},
			},
			golden: "mixed_severities.txt",
		},
		{
			name: "only warnings",
			findings: []lint.Finding{
				{RuleID: "C002", Severity: lint.SeverityWarning, File: "AGENTS.md", Line: 15, Message: `tool "kubectl" has no evidence in repo`},
				{RuleID: "C002", Severity: lint.SeverityWarning, File: "AGENTS.md", Line: 22, Message: `tool "terraform" has no evidence in repo`},
				{RuleID: "C002", Severity: lint.SeverityWarning, File: "AGENTS.md", Line: 30, Message: `tool "npm" has no evidence in repo`},
			},
			golden: "only_warnings.txt",
		},
		{
			name:      "no findings",
			findings:  nil,
			ruleCount: 5,
			golden:    "no_findings.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := reporter.WriteText(&buf, tt.findings, tt.ruleCount, reporter.Options{})
			require.NoError(t, err)
			assert.Equal(t, readGolden(t, tt.golden), buf.String())
		})
	}
}

func TestWriteText_SuccessLineIsEntireOutput(t *testing.T) {
	var buf bytes.Buffer
	err := reporter.WriteText(&buf, nil, 5, reporter.Options{})
	require.NoError(t, err)
	assert.Equal(t, "✓ AGENTS.md is valid (5 rules passed)\n", buf.String())
}

func TestWriteText_Color(t *testing.T) {
	finding := []lint.Finding{
		{RuleID: "S002", Severity: lint.SeverityError, File: "AGENTS.md", Line: 3, Message: "missing ## Agents section"},
	}

	t.Run("color disabled by default", func(t *testing.T) {
		var buf bytes.Buffer
		require.NoError(t, reporter.WriteText(&buf, finding, 0, reporter.Options{Color: false}))
		assert.NotContains(t, buf.String(), "\x1b[")
	})

	t.Run("color requested and NO_COLOR unset", func(t *testing.T) {
		t.Setenv("NO_COLOR", "")
		var buf bytes.Buffer
		require.NoError(t, reporter.WriteText(&buf, finding, 0, reporter.Options{Color: true}))
		assert.Contains(t, buf.String(), "\x1b[31merror\x1b[0m")
	})

	t.Run("NO_COLOR overrides Options.Color", func(t *testing.T) {
		t.Setenv("NO_COLOR", "1")
		var buf bytes.Buffer
		require.NoError(t, reporter.WriteText(&buf, finding, 0, reporter.Options{Color: true}))
		assert.NotContains(t, buf.String(), "\x1b[")
	})
}

func TestWriteText_WriteErrors(t *testing.T) {
	t.Run("success line write failure is propagated", func(t *testing.T) {
		err := reporter.WriteText(failingWriter{}, nil, 5, reporter.Options{})
		assert.Error(t, err)
	})

	t.Run("per-finding line write failure is propagated", func(t *testing.T) {
		err := reporter.WriteText(failingWriter{}, []lint.Finding{
			{RuleID: "S002", Severity: lint.SeverityError, File: "AGENTS.md", Line: 3, Message: "missing ## Agents section"},
		}, 0, reporter.Options{})
		assert.Error(t, err)
	})
}

func TestWriteText_UnrecognizedSeverityIsUncolored(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	var buf bytes.Buffer
	err := reporter.WriteText(&buf, []lint.Finding{
		{RuleID: "X001", Severity: lint.Severity("info"), File: "AGENTS.md", Line: 1, Message: "informational"},
	}, 0, reporter.Options{Color: true})
	require.NoError(t, err)
	assert.Equal(t, "info  X001  AGENTS.md:1 — informational\n0 error(s), 0 warning(s)\n", buf.String())
}

func TestWriteText_MessageFidelity(t *testing.T) {
	t.Run("long message is not truncated", func(t *testing.T) {
		longMessage := strings.Repeat("x", 250)
		var buf bytes.Buffer
		err := reporter.WriteText(&buf, []lint.Finding{
			{RuleID: "S003", Severity: lint.SeverityError, File: "AGENTS.md", Line: 9, Message: longMessage},
		}, 0, reporter.Options{})
		require.NoError(t, err)
		assert.Contains(t, buf.String(), longMessage)
	})

	t.Run("message with formatter-like substrings is not mangled", func(t *testing.T) {
		message := `has "quotes", an em dash — and  double  spaces`
		var buf bytes.Buffer
		err := reporter.WriteText(&buf, []lint.Finding{
			{RuleID: "S003", Severity: lint.SeverityWarning, File: "AGENTS.md", Line: 9, Message: message},
		}, 0, reporter.Options{})
		require.NoError(t, err)
		assert.Equal(t, "warning  S003  AGENTS.md:9 — "+message+"\n0 error(s), 1 warning(s)\n", buf.String())
	})
}
