package cmd

import (
	"bytes"
	"testing"

	"github.com/ingvar-goryainov/agents-lint/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func boolPtr(b bool) *bool { return &b }

func TestResolvePath(t *testing.T) {
	tests := []struct {
		name    string
		argPath string
		cfgPath string
		want    string
	}{
		{"arg wins over config path", "arg.md", "config.md", "arg.md"},
		{"config path wins over default", "", "config.md", "config.md"},
		{"default used when both empty", "", "", defaultAgentsMDPath},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, resolvePath(tt.argPath, tt.cfgPath))
		})
	}
}

func TestEffectiveRuleCount(t *testing.T) {
	knownIDs := []string{"S001", "S002", "S003", "S004", "S005"}

	tests := []struct {
		name string
		cfg  *config.Config
		want int
	}{
		{"no config", &config.Config{}, 5},
		{
			"one disabled rule",
			&config.Config{Rules: map[string]config.RuleConfig{
				"S004": {Enabled: boolPtr(false)},
			}},
			4,
		},
		{
			"enabled explicitly true does not reduce count",
			&config.Config{Rules: map[string]config.RuleConfig{
				"S004": {Enabled: boolPtr(true)},
			}},
			5,
		},
		{
			"only one map entry per rule ID, so re-assignment still counts once",
			&config.Config{Rules: map[string]config.RuleConfig{
				"S004": {Enabled: boolPtr(false)},
				"S005": {Enabled: boolPtr(false)},
			}},
			3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, effectiveRuleCount(tt.cfg, knownIDs))
		})
	}
}

func TestRunScan_Config(t *testing.T) {
	t.Run("disabling a rule removes its finding and reduces the rule count", func(t *testing.T) {
		var buf bytes.Buffer
		exitCode, err := runScan(&buf, "../../../testdata/config/s004_only.md", "../../../testdata/config/disable_s004.yaml")
		require.NoError(t, err)
		assert.Equal(t, 0, exitCode)
		assert.Contains(t, buf.String(), "AGENTS.md is valid (4 rules passed)")
	})

	t.Run("severity override to warning changes exit code from 1 to 0", func(t *testing.T) {
		var buf bytes.Buffer
		exitCode, err := runScan(&buf, "../../../testdata/config/s004_only.md", "../../../testdata/config/severity_s004_warning.yaml")
		require.NoError(t, err)
		assert.Equal(t, 0, exitCode)
		assert.Contains(t, buf.String(), "warning  S004")
	})

	t.Run("config path is used when no positional argument is given", func(t *testing.T) {
		var buf bytes.Buffer
		exitCode, err := runScan(&buf, "", "../../../testdata/config/path_override.yaml")
		require.NoError(t, err)
		assert.Equal(t, 0, exitCode)
		assert.Contains(t, buf.String(), "AGENTS.md is valid")
	})

	t.Run("a positional argument still overrides config path", func(t *testing.T) {
		var buf bytes.Buffer
		exitCode, err := runScan(&buf, "../../../testdata/config/s004_only.md", "../../../testdata/config/path_override.yaml")
		require.NoError(t, err)
		assert.Equal(t, 1, exitCode)
		assert.Contains(t, buf.String(), "S004")
	})

	t.Run("invalid config file returns an error without writing reporter output", func(t *testing.T) {
		var buf bytes.Buffer
		exitCode, err := runScan(&buf, "../../../testdata/S003/valid.md", "../../../testdata/config/unknown_rule_id.yaml")
		assert.Error(t, err)
		assert.Equal(t, 1, exitCode)
		assert.Empty(t, buf.String())
	})

	t.Run("missing --config path returns an error", func(t *testing.T) {
		var buf bytes.Buffer
		exitCode, err := runScan(&buf, "../../../testdata/S003/valid.md", "../../../testdata/config/does-not-exist.yaml")
		assert.Error(t, err)
		assert.Equal(t, 1, exitCode)
		assert.Empty(t, buf.String())
	})
}
