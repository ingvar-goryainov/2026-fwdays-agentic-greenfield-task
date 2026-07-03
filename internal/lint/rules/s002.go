package rules

import (
	"strings"

	"github.com/ingvar-goryainov/agents-lint/internal/lint"
)

// RuleS002 is the ID for the required-sections rule (FR-S002).
const RuleS002 = "S002"

// S002 checks that the document has a top-level "## Agent" or
// "## Agents" section.
type S002 struct{}

func (S002) ID() string { return RuleS002 }

func (S002) Check(doc *lint.Document) []lint.Finding {
	for _, h := range doc.H2Sections {
		t := strings.ToLower(strings.TrimSpace(h.Text))
		if t == "agent" || t == "agents" {
			return nil
		}
	}
	return []lint.Finding{{
		RuleID:   RuleS002,
		Severity: lint.SeverityError,
		File:     doc.Path,
		Message:  "missing required section — add a top-level `## Agent` or `## Agents` heading",
	}}
}
