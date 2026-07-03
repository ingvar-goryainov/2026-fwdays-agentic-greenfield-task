// Package config loads and validates .agents-lint.yaml, implementing
// FR-CFG-01..04: the default config location, path/rule overrides, the
// no-config-file default behavior, and schema validation with clear
// errors. It depends only on internal/lint (for the Severity type and
// Finding), not on internal/lint/rules, so it stays reusable by any
// caller that supplies its own list of valid rule IDs.
package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/ingvar-goryainov/agents-lint/internal/lint"
)

// DefaultFileName is the config file name looked up in the current
// directory when no explicit path is given (FR-CFG-01).
const DefaultFileName = ".agents-lint.yaml"

// RuleConfig holds per-rule overrides. Enabled and Severity are pointers
// so "not mentioned in the config" (leave default behavior alone) is
// distinguishable from "explicitly set".
type RuleConfig struct {
	Enabled  *bool
	Severity *lint.Severity
}

// Config is the parsed, validated contents of a config file.
type Config struct {
	Path  string
	Rules map[string]RuleConfig
}

// rawConfig/rawRuleConfig are the YAML decode targets, kept free of
// lint.Severity so invalid severity strings can be validated (and
// reported) explicitly rather than failing during YAML decoding.
type rawConfig struct {
	Path  string                   `yaml:"path"`
	Rules map[string]rawRuleConfig `yaml:"rules"`
}

type rawRuleConfig struct {
	Enabled  *bool   `yaml:"enabled"`
	Severity *string `yaml:"severity"`
}

// Resolve loads the effective config for a scan (FR-CFG-01, FR-CFG-03,
// FR-CLI-04). explicitPath, if non-empty, is used as-is and must exist;
// an empty explicitPath falls back to DefaultFileName in the current
// directory, and a missing default file yields a zero-value Config with
// no error (FR-CFG-03: no config file means default behavior).
func Resolve(explicitPath string, validRuleIDs []string) (*Config, error) {
	path := explicitPath
	if path == "" {
		path = DefaultFileName
		if _, err := os.Stat(path); err != nil {
			return &Config{}, nil
		}
	}
	return Load(path, validRuleIDs)
}

// Load parses and validates the config file at path against validRuleIDs
// (FR-CFG-04). An empty file is treated as a valid, empty config.
func Load(path string, validRuleIDs []string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config %s: %w", path, err)
	}
	defer func() { _ = f.Close() }()

	var raw rawConfig
	dec := yaml.NewDecoder(f)
	dec.KnownFields(true)
	if err := dec.Decode(&raw); err != nil {
		if errors.Is(err, io.EOF) {
			return &Config{}, nil
		}
		return nil, fmt.Errorf("config %s: %w", path, err)
	}

	cfg, err := raw.toConfig(validRuleIDs)
	if err != nil {
		return nil, fmt.Errorf("config %s: %w", path, err)
	}
	return cfg, nil
}

// toConfig converts raw, decoded YAML into a validated Config, checking
// the two things structural YAML decoding can't: that every rule ID
// referenced under "rules" is known, and that every "severity" value is
// "error" or "warning" (FR-CFG-04).
func (r rawConfig) toConfig(validRuleIDs []string) (*Config, error) {
	valid := make(map[string]bool, len(validRuleIDs))
	for _, id := range validRuleIDs {
		valid[id] = true
	}

	cfg := &Config{Path: r.Path}
	if len(r.Rules) > 0 {
		cfg.Rules = make(map[string]RuleConfig, len(r.Rules))
	}
	for id, rc := range r.Rules {
		if !valid[id] {
			return nil, fmt.Errorf(
				"rule %q: unknown rule ID — valid rule IDs: %s",
				id, strings.Join(validRuleIDs, ", "),
			)
		}

		out := RuleConfig{Enabled: rc.Enabled}
		if rc.Severity != nil {
			sev := lint.Severity(*rc.Severity)
			if sev != lint.SeverityError && sev != lint.SeverityWarning {
				return nil, fmt.Errorf(
					"rule %s: invalid severity %q — must be %q or %q",
					id, *rc.Severity, lint.SeverityError, lint.SeverityWarning,
				)
			}
			out.Severity = &sev
		}
		cfg.Rules[id] = out
	}
	return cfg, nil
}

// Apply filters out findings for rules disabled in c and overrides the
// severity of findings for rules with an explicit severity override
// (FR-CFG-02). Findings for rules not mentioned in c.Rules pass through
// unchanged.
func (c *Config) Apply(findings []lint.Finding) []lint.Finding {
	if c == nil || len(c.Rules) == 0 {
		return findings
	}

	out := make([]lint.Finding, 0, len(findings))
	for _, f := range findings {
		rc, ok := c.Rules[f.RuleID]
		if ok && rc.Enabled != nil && !*rc.Enabled {
			continue
		}
		if ok && rc.Severity != nil {
			f.Severity = *rc.Severity
		}
		out = append(out, f)
	}
	return out
}
