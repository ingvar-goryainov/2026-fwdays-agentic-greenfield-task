package rules

import (
	"fmt"
	"os"

	"github.com/ingvar-goryainov/agents-lint/internal/lint"
)

// RuleS001 is the ID for the file-existence rule (FR-S001).
const RuleS001 = "S001"

// CheckFileExists implements FR-S001: the target AGENTS.md file must exist
// at the given path. It is not part of the Rule interface and runs before
// parsing, since every other rule needs a readable file to operate on.
func CheckFileExists(path string) *lint.Finding {
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return &lint.Finding{
			RuleID:   RuleS001,
			Severity: lint.SeverityError,
			File:     path,
			Message: fmt.Sprintf(
				"AGENTS.md not found at %q — create it (e.g. run `agents-lint init`) or pass the correct path",
				path,
			),
		}
	}
	return nil
}

// s001Doc documents FR-S001 (FR-CLI-07). CheckFileExists has no receiver
// type to hang a Doc() method off of, so this is a plain function the
// registry calls directly, mirroring how KnownRuleIDs handles S001.
func s001Doc() lint.RuleDoc {
	return lint.RuleDoc{
		ID:          RuleS001,
		Description: "The target AGENTS.md file must exist at the resolved path before any other rule can run.",
		Severity:    lint.SeverityError,
		Valid:       "A file exists at the resolved path, e.g. ./AGENTS.md.",
		Invalid:     "No file exists at the resolved path.",
		FixGuidance: "Create AGENTS.md (e.g. run `agents-lint init`), or pass the correct path as the scan argument or via .agents-lint.yaml's `path:` key.",
	}
}
