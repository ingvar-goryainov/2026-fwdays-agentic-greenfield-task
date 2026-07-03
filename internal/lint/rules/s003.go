package rules

import (
	"fmt"
	"strings"

	"github.com/ingvar-goryainov/agents-lint/internal/lint"
)

// RuleS003 is the ID for the agent-block-completeness rule (FR-S003).
const RuleS003 = "S003"

// S003 checks that every agent block has a name, a role/description, and
// at least one of instructions, tools, or context.
type S003 struct{}

func (S003) ID() string { return RuleS003 }

func (S003) Check(doc *lint.Document) []lint.Finding {
	var findings []lint.Finding
	for _, b := range doc.AgentBlocks {
		if strings.TrimSpace(b.Name) == "" {
			findings = append(findings, lint.Finding{
				RuleID:   RuleS003,
				Severity: lint.SeverityError,
				File:     doc.Path,
				Line:     b.Line,
				Message:  "agent block has no name — give the `###` heading a non-empty name",
			})
		}
		if !b.HasRole {
			findings = append(findings, lint.Finding{
				RuleID:   RuleS003,
				Severity: lint.SeverityError,
				File:     doc.Path,
				Line:     b.Line,
				Message: fmt.Sprintf(
					"agent %q has no role/description — add a description paragraph or a **Role:** field",
					b.Name,
				),
			})
		}
		if !b.HasInstructions && !b.HasTools && !b.HasContext {
			findings = append(findings, lint.Finding{
				RuleID:   RuleS003,
				Severity: lint.SeverityError,
				File:     doc.Path,
				Line:     b.Line,
				Message: fmt.Sprintf(
					"agent %q has none of instructions, tools, or context — add a **Instructions:**, **Tools:**, or **Context:** field",
					b.Name,
				),
			})
		}
	}
	return findings
}
