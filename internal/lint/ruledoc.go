package lint

// RuleDoc is a rule's human-facing documentation: what it checks, its
// default severity, a short example of input that passes and one that
// fails, and guidance on how to fix a violation. Populated by each rule
// alongside its Check logic and surfaced by `agents-lint docs <rule-id>`
// (FR-CLI-07).
type RuleDoc struct {
	ID          string
	Description string
	Severity    Severity
	Valid       string
	Invalid     string
	FixGuidance string
}
