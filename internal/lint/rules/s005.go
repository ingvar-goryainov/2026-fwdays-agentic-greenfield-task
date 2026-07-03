package rules

import (
	"fmt"

	"github.com/ingvar-goryainov/agents-lint/internal/lint"
	"gopkg.in/yaml.v3"
)

// RuleS005 is the ID for the valid-frontmatter-YAML rule (FR-S005).
const RuleS005 = "S005"

// S005 checks that, if present, the YAML frontmatter block parses as
// valid YAML. Files with no frontmatter are not flagged.
type S005 struct{}

func (S005) ID() string { return RuleS005 }

func (S005) Check(doc *lint.Document) []lint.Finding {
	if doc.Frontmatter == nil {
		return nil
	}

	var out any
	if err := yaml.Unmarshal([]byte(doc.Frontmatter.Raw), &out); err != nil {
		return []lint.Finding{{
			RuleID:   RuleS005,
			Severity: lint.SeverityError,
			File:     doc.Path,
			Line:     doc.Frontmatter.Line,
			Message:  fmt.Sprintf("frontmatter is not valid YAML: %v", err),
		}}
	}
	return nil
}

// Doc documents FR-S005 (FR-CLI-07).
func (S005) Doc() lint.RuleDoc {
	return lint.RuleDoc{
		ID:          RuleS005,
		Description: "If the document has YAML frontmatter (delimited by `---`), it must parse as valid YAML. Files without frontmatter are not flagged.",
		Severity:    lint.SeverityError,
		Valid:       "---\nname: my-project\n---\n\n## Agents\n...\n",
		Invalid:     "---\nname: \"unterminated\n---\n\n## Agents\n...\n",
		FixGuidance: "Fix the YAML syntax error in the frontmatter block (check quoting, indentation, and colons), or remove the frontmatter entirely if it isn't needed.",
	}
}
