package rules

import "github.com/ingvar-goryainov/agents-lint/internal/lint"

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

// Run validates the AGENTS.md file at path: it checks file existence
// first (S001) and, only if that succeeds, parses the file and runs the
// remaining schema rules against it. A non-nil error indicates a problem
// unrelated to the file simply not existing (e.g. a read/parse failure);
// schema violations are reported as Findings, not errors.
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
	return findings, nil
}
