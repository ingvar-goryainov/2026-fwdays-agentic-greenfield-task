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
