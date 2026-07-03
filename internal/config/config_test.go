package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ingvar-goryainov/agents-lint/internal/config"
	"github.com/ingvar-goryainov/agents-lint/internal/lint"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var validRuleIDs = []string{"S001", "S002", "S003", "S004", "S005"}

func severityPtr(s lint.Severity) *lint.Severity { return &s }
func boolPtr(b bool) *bool                       { return &b }

func TestLoad_Valid(t *testing.T) {
	t.Run("path only", func(t *testing.T) {
		cfg, err := config.Load("../../testdata/config/valid_path.yaml", validRuleIDs)
		require.NoError(t, err)
		assert.Equal(t, "docs/AGENTS.md", cfg.Path)
		assert.Empty(t, cfg.Rules)
	})

	t.Run("rule overrides", func(t *testing.T) {
		cfg, err := config.Load("../../testdata/config/valid_rules.yaml", validRuleIDs)
		require.NoError(t, err)
		assert.Empty(t, cfg.Path)
		require.Contains(t, cfg.Rules, "S004")
		require.NotNil(t, cfg.Rules["S004"].Enabled)
		assert.False(t, *cfg.Rules["S004"].Enabled)

		require.Contains(t, cfg.Rules, "S002")
		require.NotNil(t, cfg.Rules["S002"].Severity)
		assert.Equal(t, lint.SeverityWarning, *cfg.Rules["S002"].Severity)
	})

	t.Run("empty file", func(t *testing.T) {
		cfg, err := config.Load("../../testdata/config/empty.yaml", validRuleIDs)
		require.NoError(t, err)
		assert.Equal(t, &config.Config{}, cfg)
	})
}

func TestLoad_Invalid(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		wantErrHas []string
	}{
		{
			"malformed yaml",
			"../../testdata/config/malformed.yaml",
			[]string{"../../testdata/config/malformed.yaml"},
		},
		{
			"unknown top-level key",
			"../../testdata/config/unknown_top_level_key.yaml",
			[]string{"strict"},
		},
		{
			"unknown rule key",
			"../../testdata/config/unknown_rule_key.yaml",
			[]string{"enalbed"},
		},
		{
			"unknown rule id",
			"../../testdata/config/unknown_rule_id.yaml",
			[]string{"S0O3", "S001", "S002", "S003", "S004", "S005"},
		},
		{
			"invalid severity",
			"../../testdata/config/invalid_severity.yaml",
			[]string{"S002", "critical"},
		},
		{
			"missing file",
			"../../testdata/config/does-not-exist.yaml",
			[]string{"../../testdata/config/does-not-exist.yaml"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := config.Load(tt.path, validRuleIDs)
			require.Error(t, err)
			assert.Nil(t, cfg)
			for _, want := range tt.wantErrHas {
				assert.Contains(t, err.Error(), want)
			}
		})
	}
}

func TestResolve(t *testing.T) {
	t.Run("explicit path used when given", func(t *testing.T) {
		cfg, err := config.Resolve("../../testdata/config/valid_path.yaml", validRuleIDs)
		require.NoError(t, err)
		assert.Equal(t, "docs/AGENTS.md", cfg.Path)
	})

	t.Run("explicit path missing is an error", func(t *testing.T) {
		cfg, err := config.Resolve("../../testdata/config/does-not-exist.yaml", validRuleIDs)
		require.Error(t, err)
		assert.Nil(t, cfg)
		assert.Contains(t, err.Error(), "does-not-exist.yaml")
	})

	t.Run("default file loaded when present", func(t *testing.T) {
		dir := t.TempDir()
		defaultPath := filepath.Join(dir, config.DefaultFileName)
		require.NoError(t, os.WriteFile(defaultPath, []byte("path: docs/AGENTS.md\n"), 0o644))

		wd, err := os.Getwd()
		require.NoError(t, err)
		require.NoError(t, os.Chdir(dir))
		t.Cleanup(func() { _ = os.Chdir(wd) })

		cfg, err := config.Resolve("", validRuleIDs)
		require.NoError(t, err)
		assert.Equal(t, "docs/AGENTS.md", cfg.Path)
	})

	t.Run("no explicit path and no default file returns empty config", func(t *testing.T) {
		dir := t.TempDir()

		wd, err := os.Getwd()
		require.NoError(t, err)
		require.NoError(t, os.Chdir(dir))
		t.Cleanup(func() { _ = os.Chdir(wd) })

		cfg, err := config.Resolve("", validRuleIDs)
		require.NoError(t, err)
		assert.Equal(t, &config.Config{}, cfg)
	})
}

func TestConfig_Apply(t *testing.T) {
	findings := []lint.Finding{
		{RuleID: "S002", Severity: lint.SeverityError, Message: "s002 finding"},
		{RuleID: "S004", Severity: lint.SeverityError, Message: "s004 finding"},
		{RuleID: "S005", Severity: lint.SeverityError, Message: "s005 finding"},
	}

	t.Run("nil config returns findings unchanged", func(t *testing.T) {
		var cfg *config.Config
		assert.Equal(t, findings, cfg.Apply(findings))
	})

	t.Run("empty config returns findings unchanged", func(t *testing.T) {
		cfg := &config.Config{}
		assert.Equal(t, findings, cfg.Apply(findings))
	})

	t.Run("disabled rule findings are removed, others untouched", func(t *testing.T) {
		cfg := &config.Config{Rules: map[string]config.RuleConfig{
			"S004": {Enabled: boolPtr(false)},
		}}
		got := cfg.Apply(findings)
		require.Len(t, got, 2)
		for _, f := range got {
			assert.NotEqual(t, "S004", f.RuleID)
		}
	})

	t.Run("severity override changes matching findings only", func(t *testing.T) {
		cfg := &config.Config{Rules: map[string]config.RuleConfig{
			"S002": {Severity: severityPtr(lint.SeverityWarning)},
		}}
		got := cfg.Apply(findings)
		require.Len(t, got, 3)
		for _, f := range got {
			if f.RuleID == "S002" {
				assert.Equal(t, lint.SeverityWarning, f.Severity)
			} else {
				assert.Equal(t, lint.SeverityError, f.Severity)
			}
		}
	})

	t.Run("severity can be raised from warning to error", func(t *testing.T) {
		cfg := &config.Config{Rules: map[string]config.RuleConfig{
			"C002": {Severity: severityPtr(lint.SeverityError)},
		}}
		got := cfg.Apply([]lint.Finding{{RuleID: "C002", Severity: lint.SeverityWarning}})
		require.Len(t, got, 1)
		assert.Equal(t, lint.SeverityError, got[0].Severity)
	})

	t.Run("rule not mentioned in config passes through unchanged", func(t *testing.T) {
		cfg := &config.Config{Rules: map[string]config.RuleConfig{
			"S002": {Enabled: boolPtr(false)},
		}}
		got := cfg.Apply(findings)
		require.Len(t, got, 2)
		assert.Equal(t, findings[1], got[0])
		assert.Equal(t, findings[2], got[1])
	})
}
