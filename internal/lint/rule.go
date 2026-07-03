package lint

// Rule validates a parsed Document and reports any Findings. Concrete
// implementations live in internal/lint/rules, not here — this package
// stays free of rule logic so internal/lint/rules can import it without
// creating an import cycle.
type Rule interface {
	ID() string
	Check(doc *Document) []Finding
}

// CodebaseRule validates a parsed Document against the surrounding
// repository, not just against itself. repoRoot is the directory rules
// resolve relative paths and evidence-file lookups against. Kept separate
// from Rule (rather than adding repoRoot to Rule.Check) so the schema
// rules, which never need a repo root, are unaffected.
type CodebaseRule interface {
	ID() string
	Check(doc *Document, repoRoot string) []Finding
}
