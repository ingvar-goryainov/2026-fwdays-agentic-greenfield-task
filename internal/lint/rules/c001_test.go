package rules_test

import (
	"testing"

	"github.com/ingvar-goryainov/agents-lint/internal/lint"
	"github.com/ingvar-goryainov/agents-lint/internal/lint/rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// repoRoot is the actual repository root relative to this package
// directory (internal/lint/rules), used so C001 fixtures can reference
// real files without needing a synthetic repo layout.
const repoRoot = "../../.."

func TestC001_ID(t *testing.T) {
	assert.Equal(t, "C001", rules.C001{}.ID())
}

func TestC001_Check(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		wantLines int
	}{
		{"real paths and non-path code spans", "../../../testdata/C001/valid.md", 0},
		{"broken paths", "../../../testdata/C001/invalid.md", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := lint.Parse(tt.path)
			require.NoError(t, err)

			findings := rules.C001{}.Check(doc, repoRoot)
			require.Len(t, findings, tt.wantLines)
			for _, f := range findings {
				assert.Equal(t, rules.RuleC001, f.RuleID)
				assert.Equal(t, lint.SeverityError, f.Severity)
				assert.Equal(t, doc.Path, f.File)
				assert.NotZero(t, f.Line)
			}
		})
	}
}

func TestC001_Check_NoCodeSpans(t *testing.T) {
	doc := &lint.Document{Path: "AGENTS.md"}
	assert.Empty(t, rules.C001{}.Check(doc, repoRoot))
}

// TestC001_LooksLikePathHeuristic exercises the "looks like a path"
// heuristic through the public Check API rather than calling the
// unexported looksLikePath directly (matching this codebase's convention
// of testing unexported helpers via their exported entry point, e.g.
// internal/config's toConfig is only exercised through Load). repoRoot is
// an empty temp directory, so every candidate text is guaranteed not to
// exist — a finding means the text was classified as a path candidate; no
// finding means it was correctly ignored.
func TestC001_LooksLikePathHeuristic(t *testing.T) {
	tests := []struct {
		name        string
		text        string
		wantFinding bool
	}{
		{"contains whitespace", "some path/file.go", false},
		{"cli flag", "--config", false},
		{"short cli flag", "-c", false},
		{"url", "https://example.com", false},
		{"glob", "**/*.go", false},
		{"root-relative rest path", "/api/users", false},
		{"multi-segment path", "internal/lint/rule.go", true},
		{"directory-shaped path", "internal/lint/", true},
		{"bare filename with extension", "go.mod", true},
		{"dotfile with extension", "README.md", true},
		{"leading-dot config file", ".agents-lint.yaml", true},
		{"bare word no dot", "terraform", false},
		{"decimal version", "1.22", false},
		{"semver version", "v1.2.3", false},
		{"template placeholder with slash", "internal/rules/<rule_id>.go", false},
		{"bare template placeholder", "<rule_id>.go", false},
		{"go qualified identifier, exported", "fmt.Print", false},
		{"go qualified identifier, exported multi-word", "os.ReadFile", false},
	}

	emptyRoot := t.TempDir()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := &lint.Document{
				Path:      "AGENTS.md",
				CodeSpans: []lint.CodeSpan{{Text: tt.text, Line: 7}},
			}
			findings := rules.C001{}.Check(doc, emptyRoot)
			if !tt.wantFinding {
				assert.Empty(t, findings)
				return
			}
			require.Len(t, findings, 1)
			assert.Equal(t, 7, findings[0].Line)
		})
	}
}
