package lint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindingConstruction(t *testing.T) {
	f := Finding{
		RuleID:   "S001",
		Severity: SeverityError,
		File:     "AGENTS.md",
		Line:     1,
		Message:  "AGENTS.md file must exist",
	}

	assert.Equal(t, "S001", f.RuleID)
	assert.Equal(t, SeverityError, f.Severity)
	assert.Equal(t, "AGENTS.md", f.File)
	assert.Equal(t, 1, f.Line)
	assert.Equal(t, "AGENTS.md file must exist", f.Message)
}
