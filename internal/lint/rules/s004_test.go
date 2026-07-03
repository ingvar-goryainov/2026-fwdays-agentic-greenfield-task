package rules_test

import (
	"testing"

	"github.com/ingvar-goryainov/agents-lint/internal/lint"
	"github.com/ingvar-goryainov/agents-lint/internal/lint/rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestS004_Check(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		wantFinding bool
	}{
		{"distinct names", "../../../testdata/S004/valid.md", false},
		{"duplicate names, different casing", "../../../testdata/S004/invalid.md", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := lint.Parse(tt.path)
			require.NoError(t, err)

			findings := rules.S004{}.Check(doc)
			if !tt.wantFinding {
				assert.Empty(t, findings)
				return
			}
			require.Len(t, findings, 1)
			assert.Equal(t, rules.RuleS004, findings[0].RuleID)
			assert.Equal(t, lint.SeverityError, findings[0].Severity)
		})
	}
}
