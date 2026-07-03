package rules

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ingvar-goryainov/agents-lint/internal/lint"
)

// RuleC001 is the ID for the file-path-reference rule (FR-C001).
const RuleC001 = "C001"

// C001 checks that every inline-code span in AGENTS.md that looks like a
// filesystem path resolves to an existing file or directory relative to
// the repo root.
type C001 struct{}

func (C001) ID() string { return RuleC001 }

func (C001) Check(doc *lint.Document, repoRoot string) []lint.Finding {
	var findings []lint.Finding
	for _, span := range doc.CodeSpans {
		if !looksLikePath(span.Text) {
			continue
		}
		if _, err := os.Stat(filepath.Join(repoRoot, span.Text)); err != nil {
			findings = append(findings, lint.Finding{
				RuleID:   RuleC001,
				Severity: lint.SeverityError,
				File:     doc.Path,
				Line:     span.Line,
				Message: fmt.Sprintf(
					"referenced path %q does not exist relative to the repo root — fix the path or update it if the file was moved/renamed",
					span.Text,
				),
			})
		}
	}
	return findings
}

// Doc documents FR-C001 (FR-CLI-07).
func (C001) Doc() lint.RuleDoc {
	return lint.RuleDoc{
		ID:          RuleC001,
		Description: "Every inline-code span that looks like a filesystem path (e.g. `internal/lint/rule.go`) must resolve to an existing file or directory relative to the repo root.",
		Severity:    lint.SeverityError,
		Valid:       "**Context:** `go.mod` and `internal/lint/rule.go`.",
		Invalid:     "**Context:** `internal/does/not/exist.go`.",
		FixGuidance: "Fix the path, or update/remove the reference if the file was moved, renamed, or deleted.",
	}
}

// looksLikePath reports whether text is a plausible filesystem path
// reference, per FR-C001's own "detected by ... patterns" wording: a
// heuristic, not exact parsing. See design.md for the reasoning behind
// each exclusion/inclusion.
func looksLikePath(text string) bool {
	switch {
	case text == "",
		strings.ContainsAny(text, " \t\n\r"),
		strings.HasPrefix(text, "-"),
		strings.Contains(text, "://"),
		strings.ContainsAny(text, "*?<>"),
		strings.HasPrefix(text, "/"):
		return false
	}

	if strings.Contains(text, "/") {
		return true
	}
	return isBareFileName(text)
}

// isBareFileName reports whether text is a single-segment "name.ext"-shaped
// token whose final dot-segment starts with a lowercase letter — this
// catches bare filenames like go.mod/README.md/.agents-lint.yaml (real
// file extensions are conventionally lowercase) while excluding both
// version-looking tokens like 1.22/v1.2.3 (digit-first) and Go-idiom
// qualified identifiers like fmt.Print/os.ReadFile (uppercase-first,
// since exported Go identifiers are capitalized) that are otherwise
// indistinguishable from a bare filename by shape alone.
func isBareFileName(text string) bool {
	idx := strings.LastIndex(text, ".")
	if idx <= 0 || idx == len(text)-1 {
		return false
	}
	ext := text[idx+1:]
	return ext[0] >= 'a' && ext[0] <= 'z'
}
