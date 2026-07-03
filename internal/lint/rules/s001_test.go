package rules_test

import (
	"path/filepath"
	"testing"

	"github.com/ingvar-goryainov/agents-lint/internal/lint"
	"github.com/ingvar-goryainov/agents-lint/internal/lint/rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckFileExists(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		wantFinding bool
	}{
		{"existing file", "../../../testdata/S001/valid.md", false},
		{"missing file", filepath.Join(t.TempDir(), "does-not-exist.md"), true},
		{"path is a directory", t.TempDir(), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			finding := rules.CheckFileExists(tt.path)
			if !tt.wantFinding {
				assert.Nil(t, finding)
				return
			}
			require.NotNil(t, finding)
			assert.Equal(t, rules.RuleS001, finding.RuleID)
			assert.Equal(t, lint.SeverityError, finding.Severity)
			assert.NotEmpty(t, finding.Message)
		})
	}
}
