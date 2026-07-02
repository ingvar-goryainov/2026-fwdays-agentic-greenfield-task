// Package lint is the core validation library for agents-lint. It is kept
// independent of cmd/ so future consumers (HTTP API, editor extensions) can
// depend on it without depending on the CLI.
package lint

// Severity indicates whether a Finding should fail the scan (SeverityError)
// or only be reported (SeverityWarning).
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
)

// Finding is a single validation result produced by a rule.
type Finding struct {
	RuleID   string
	Severity Severity
	File     string
	Line     int
	Message  string
}
