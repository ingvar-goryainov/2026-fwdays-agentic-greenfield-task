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
