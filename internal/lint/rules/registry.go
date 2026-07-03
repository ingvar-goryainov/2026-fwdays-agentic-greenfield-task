package rules

import (
	"fmt"
	"os"

	"github.com/ingvar-goryainov/agents-lint/internal/lint"
)

// DefaultRules returns the schema rules in a fixed, deterministic order
// (S002, S003, S004, S005). S001 is run separately by Run, since it checks
// file existence before there is anything to parse.
func DefaultRules() []lint.Rule {
	return []lint.Rule{
		S002{},
		S003{},
		S004{},
		S005{},
	}
}

// CodebaseAwareRules returns the codebase-awareness rules (C001, C002) in
// a fixed, deterministic order. They run after the schema rules, against
// the same parsed Document, since they need it too.
func CodebaseAwareRules() []lint.CodebaseRule {
	return []lint.CodebaseRule{
		C001{},
		C002{},
	}
}

// Run validates the AGENTS.md file at path: it checks file existence
// first (S001) and, only if that succeeds, parses the file and runs the
// remaining schema rules and the codebase-awareness rules against it. A
// non-nil error indicates a problem unrelated to the file simply not
// existing (e.g. a read/parse failure); rule violations are reported as
// Findings, not errors.
func Run(path string) ([]lint.Finding, error) {
	if finding := CheckFileExists(path); finding != nil {
		return []lint.Finding{*finding}, nil
	}

	doc, err := lint.Parse(path)
	if err != nil {
		return nil, err
	}

	var findings []lint.Finding
	for _, rule := range DefaultRules() {
		findings = append(findings, rule.Check(doc)...)
	}

	repoRoot, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("rules: resolving repo root: %w", err)
	}
	for _, rule := range CodebaseAwareRules() {
		findings = append(findings, rule.Check(doc, repoRoot)...)
	}
	return findings, nil
}

// KnownRuleIDs returns every rule ID Run evaluates, including S001 (which
// runs outside the Rule interface). Used by configuration-support to
// validate rule IDs referenced in .agents-lint.yaml.
func KnownRuleIDs() []string {
	ids := []string{RuleS001}
	for _, r := range DefaultRules() {
		ids = append(ids, r.ID())
	}
	for _, r := range CodebaseAwareRules() {
		ids = append(ids, r.ID())
	}
	return ids
}

// documented is satisfied by any Rule/CodebaseRule whose concrete type
// also implements Doc(). Kept local to this file rather than added to the
// Rule/CodebaseRule interfaces themselves, so rule-docs-command stays
// additive: existing call sites of Rule/CodebaseRule are unaffected.
type documented interface {
	Doc() lint.RuleDoc
}

// RuleDocs returns the human-facing documentation for every rule Run
// evaluates, keyed by rule ID (FR-CLI-07). Mirrors KnownRuleIDs's set of
// IDs; S001 is added directly since, like KnownRuleIDs, it runs outside
// the Rule interface.
func RuleDocs() map[string]lint.RuleDoc {
	docs := map[string]lint.RuleDoc{
		RuleS001: s001Doc(),
	}
	for _, r := range DefaultRules() {
		if d, ok := r.(documented); ok {
			docs[r.ID()] = d.Doc()
		}
	}
	for _, r := range CodebaseAwareRules() {
		if d, ok := r.(documented); ok {
			docs[r.ID()] = d.Doc()
		}
	}
	return docs
}
