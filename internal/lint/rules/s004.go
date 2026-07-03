package rules

import (
	"fmt"
	"strings"

	"github.com/ingvar-goryainov/agents-lint/internal/lint"
)

// RuleS004 is the ID for the no-duplicate-agent-names rule (FR-S004).
const RuleS004 = "S004"

// S004 checks that no two agent blocks share the same name
// (case-insensitive, trimmed).
type S004 struct{}

func (S004) ID() string { return RuleS004 }

func (S004) Check(doc *lint.Document) []lint.Finding {
	seen := make(map[string]bool)
	var findings []lint.Finding
	for _, b := range doc.AgentBlocks {
		key := strings.ToLower(strings.TrimSpace(b.Name))
		if key == "" {
			continue
		}
		if seen[key] {
			findings = append(findings, lint.Finding{
				RuleID:   RuleS004,
				Severity: lint.SeverityError,
				File:     doc.Path,
				Line:     b.Line,
				Message: fmt.Sprintf(
					"duplicate agent name %q — agent names must be unique (case-insensitive)",
					b.Name,
				),
			})
			continue
		}
		seen[key] = true
	}
	return findings
}

// Doc documents FR-S004 (FR-CLI-07).
func (S004) Doc() lint.RuleDoc {
	return lint.RuleDoc{
		ID:          RuleS004,
		Description: "No two agent blocks may share the same name, compared case-insensitively after trimming whitespace.",
		Severity:    lint.SeverityError,
		Valid:       "### reviewer\n...\n\n### deployer\n...\n",
		Invalid:     "### reviewer\n...\n\n### Reviewer\n...\n",
		FixGuidance: "Rename one of the duplicate agent blocks so every `###` heading under Agent(s) is unique, ignoring case.",
	}
}
