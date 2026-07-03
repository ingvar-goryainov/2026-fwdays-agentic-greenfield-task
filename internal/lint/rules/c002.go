package rules

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ingvar-goryainov/agents-lint/internal/lint"
)

// RuleC002 is the ID for the tool/command-evidence rule (FR-C002).
const RuleC002 = "C002"

// toolEvidence maps a recognized tool/command name to the repo-root
// filenames/globs/directories that would corroborate the repo actually
// using it. Only inline-code spans whose text is an exact key here are
// considered a tool reference — see design.md for why the catalog is
// fixed rather than heuristic.
var toolEvidence = map[string][]string{
	"terraform": {"*.tf", ".terraform.lock.hcl"},
	"npm":       {"package.json", "package-lock.json"},
	"yarn":      {"yarn.lock"},
	"pnpm":      {"pnpm-lock.yaml"},
	"docker":    {"Dockerfile"},
	"go":        {"go.mod"},
	"make":      {"Makefile"},
	"cargo":     {"Cargo.toml"},
	"kubectl":   {"k8s/", "kubernetes/", "kustomization.yaml"},
	"git":       {".git"},
}

// C002 checks that every inline-code span in AGENTS.md matching a
// recognized tool/command name has corresponding evidence at the repo
// root: a config file, lockfile, or other marker file associated with
// that tool.
type C002 struct{}

func (C002) ID() string { return RuleC002 }

func (C002) Check(doc *lint.Document, repoRoot string) []lint.Finding {
	var findings []lint.Finding
	for _, span := range doc.CodeSpans {
		patterns, ok := toolEvidence[span.Text]
		if !ok {
			continue
		}
		if hasEvidence(repoRoot, patterns) {
			continue
		}
		findings = append(findings, lint.Finding{
			RuleID:   RuleC002,
			Severity: lint.SeverityWarning,
			File:     doc.Path,
			Line:     span.Line,
			Message: fmt.Sprintf(
				"AGENTS.md references tool %q but the repo has no corresponding evidence at its root (expected one of: %s) — confirm the tool is actually used, or remove the reference",
				span.Text, strings.Join(patterns, ", "),
			),
		})
	}
	return findings
}

// hasEvidence reports whether any of patterns is present directly at
// repoRoot (non-recursive — see design.md's Non-Goals). Glob-shaped
// patterns are matched with filepath.Glob; everything else is checked
// with os.Stat, which also covers directory-shaped patterns like "k8s/".
func hasEvidence(repoRoot string, patterns []string) bool {
	for _, p := range patterns {
		full := filepath.Join(repoRoot, p)
		if strings.ContainsAny(p, "*?[") {
			matches, err := filepath.Glob(full)
			if err == nil && len(matches) > 0 {
				return true
			}
			continue
		}
		if _, err := os.Stat(full); err == nil {
			return true
		}
	}
	return false
}
