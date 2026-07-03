package lint

// Rule validates a parsed Document and reports any Findings. Concrete
// implementations live in internal/lint/rules, not here — this package
// stays free of rule logic so internal/lint/rules can import it without
// creating an import cycle.
type Rule interface {
	ID() string
	Check(doc *Document) []Finding
}
